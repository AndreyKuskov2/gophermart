package storage

import (
	"context"
	"fmt"

	// "github.com/golang-migrate/migrate/v4"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgres(dbURI string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db storage: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot ping database: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot create migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("cannot create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("cannot to apply migrations: %v", err)
	}

	return &Postgres{
		DB: pool,
	}, nil
}

func (db *Postgres) CreateUser(ctx context.Context, user models.UserCreditials) (int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, fmt.Errorf("cannot hashing password: %v", err)
	}

	var userID int
	if err := db.DB.QueryRow(ctx, checkUserIsExists, user.Login).Scan(&userID); err == nil {
		return 0, ErrUserIsExist
	}

	if err := db.DB.QueryRow(ctx, createNewUser, user.Login, string(passwordHash)).Scan(&userID); err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("cannot create user: %v", err)
	}
	return userID, nil
}

func (db *Postgres) GetUserByLogin(ctx context.Context, user models.UserCreditials) (int, error) {
	var userID int
	var passwordHash string

	if err := db.DB.QueryRow(ctx, getUserPasswordByLogin, user.Login).Scan(&userID, &passwordHash); err != nil {
		return 0, fmt.Errorf("user not found: %v", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	if err != nil {
		return 0, ErrInvalidData
	}
	return userID, nil
}

func (db *Postgres) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Orders, error) {
	rows, err := db.DB.Query(ctx, getOrderByNumber, orderNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Orders])
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (db *Postgres) CreateNewOrder(ctx context.Context, order *models.Orders) error {
	if _, err := db.DB.Exec(ctx, createOrder, order.Number, order.Status, order.Accrual, order.UserID); err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Orders, error) {
	rows, err := db.DB.Query(ctx, getOrdersByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Orders])
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (db *Postgres) GetUserBalance(ctx context.Context, userID string) (*models.Balance, error) {
	var balance models.Balance
	if err := db.DB.QueryRow(ctx, getUserBalance, userID, "PROCESSED").Scan(&balance.Current, &balance.Withdrawn); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (db *Postgres) CreateWithdrawal(ctx context.Context, withdrawal *models.WithdrawBalance) error {
	if _, err := db.DB.Exec(ctx, createWithdraw, withdrawal.UserID, withdrawal.OrderNumber, withdrawal.Amount); err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetWithdrawalByUserID(ctx context.Context, userID string) ([]models.WithdrawBalance, error) {
	rows, err := db.DB.Query(ctx, getWithdrawalByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	withdrawBalance, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.WithdrawBalance])
	if err != nil {
		return nil, err
	}
	return withdrawBalance, nil
}

func (db *Postgres) GetPendingOrders(ctx context.Context) ([]models.Orders, error) {
	rows, err := db.DB.Query(ctx, getPendingOrders, "NEW", "PROCESSING")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Orders])
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (db *Postgres) UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual *float32) error {
	if _, err := db.DB.Exec(ctx, updateOrderStatus, status, accrual, orderNumber); err != nil {
		return err
	}
	return nil
}
