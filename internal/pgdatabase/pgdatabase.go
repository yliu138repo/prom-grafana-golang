package pgdatabase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgresdb struct {
	Pool *pgxpool.Pool
}

var (
	pgInstance *Postgresdb
	pgOnce     sync.Once
)

/*
func NewPostgresDB(ctx context.Context, connString string) (*Postgresdb, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)

		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		pgInstance = &Postgresdb{DB: db}
	})
	return pgInstance, nil
}
*/

func NewPostgresDB(ctx context.Context, connString string) (*Postgresdb, error) {

	log.Printf("Attempting to connect with: %s\n", connString)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %v", err)
	}

	// Set some pool configuration
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &Postgresdb{Pool: pool}, nil
}

func (pg *Postgresdb) Ping(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := pg.Pool.Ping(ctxWithTimeout)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

func (pg *Postgresdb) GetStats() *pgxpool.Stat {
	return pg.Pool.Stat()
}

func (pg *Postgresdb) PingWithRetry(ctx context.Context) error {
	for i := 0; i < 3; i++ {
		err := pg.Ping(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("failed to ping database after 3 attempts")
}

func (pg *Postgresdb) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}

func getAllUsersFromDB(ctx context.Context, pool *pgxpool.Pool) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := "SELECT id, name, email, created_at, updated_at FROM public.users"
	rows, err := pool.Query(ctx, query)
	if err != nil {
		fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	users := []UserModel{}

	for rows.Next() {
		user := UserModel{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Errorf("unable to query users: %w", err)
			return
		}

		users = append(users, user)
	}

}

type UserModel struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt *time.Time // nullable
}
