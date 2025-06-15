package app

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/AndreyKuskov2/gophermart/internal/client"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"go.uber.org/zap"
)

type IOrdersStorage interface {
	GetPendingOrders() ([]models.Orders, error)
	UpdateOrderStatus(orderNumber, status string, accrual *int) error
}

type AccrualProcessor struct {
	storage       IOrdersStorage
	accrualClient *client.Client
	Log           *logger.Logger
	workerCount   int
}

func NewAccrualProcessor(orderRepository IOrdersStorage, accrualClient *client.Client, log *logger.Logger) *AccrualProcessor {
	return &AccrualProcessor{
		storage:       orderRepository,
		accrualClient: accrualClient,
		Log:           log,
	}
}

func (p *AccrualProcessor) Run(ctx context.Context, interval int, workerCount int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processPendingOrders(ctx, workerCount)
		}
	}
}

func (p *AccrualProcessor) processPendingOrders(ctx context.Context, workerCount int) {
	orders, err := p.storage.GetPendingOrders()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.Log.Log.Info("no pending orders found")
			return
		}
		p.Log.Log.Error("failed to get pending orders", zap.Error(err))
		return
	}

	jobs := make(chan models.Orders)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for order := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					p.processOrder(ctx, order)
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, order := range orders {
			select {
			case jobs <- order:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
}

func (p *AccrualProcessor) processOrder(ctx context.Context, order models.Orders) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, retryAfter, err := p.accrualClient.GetOrderInfo(reqCtx, order.Number)
	if err != nil {
		p.Log.Log.Info("failed to get order info", zap.String("order_number", order.Number), zap.Error(err))
		return
	}

	if retryAfter > 0 {
		p.Log.Log.Info("accrual service is busy, retrying later",
			zap.String("order_number", order.Number), zap.Int("retry_after", retryAfter))
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(retryAfter) * time.Second):
			return
		}
	}

	if response == nil {
		p.Log.Log.Info("empty response from accrual service", zap.String("order_number", order.Number))
		return
	}

	var newAccrual *int
	if response.Status == "PROCESSED" {
		newAccrual = &response.Accrual
	}

	if err = p.storage.UpdateOrderStatus(order.Number, response.Status, newAccrual); err != nil {
		p.Log.Log.Info("failed to update order accrual", zap.String("order_number", order.Number), zap.Error(err))
		return
	}
}
