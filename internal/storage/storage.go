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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type Postgres struct {
	DB  *pgxpool.Pool
	Ctx context.Context
}

func NewPostgres(ctx context.Context, dbURI string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db storage: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
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
		DB:  pool,
		Ctx: ctx,
	}, nil
}

func (db *Postgres) CreateUser(user models.UserCreditials) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return fmt.Errorf("cannot hashing password: %v", err)
	}

	var userID int
	if err := db.DB.QueryRow(db.Ctx, checkUserIsExists, user.Login).Scan(&userID); err == nil {
		return ErrUserIsExist
	}

	if _, err := db.DB.Exec(db.Ctx, createNewUser, user.Login, string(passwordHash)); err != nil {
		fmt.Println(err)
		return fmt.Errorf("cannot create user: %v", err)
	}
	return nil
}

func (db *Postgres) GetUserByLogin(user models.UserCreditials) error {
	var passwordHash string

	if err := db.DB.QueryRow(db.Ctx, getUserPasswordByLogin, user.Login).Scan(&passwordHash); err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	if err != nil {
		return ErrInvalidData
	}
	return nil
}
