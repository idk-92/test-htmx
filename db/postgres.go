package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Connect initializes a TCP connection pool for a mysql connection
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	godotenv.Load()

	var (
		err       error
		dbUser    = os.Getenv("DB_USER") // e.g. 'my-db-user'
		dbPwd     = os.Getenv("DB_PASS") // e.g. 'my-db-password'
		dbTCPHost = os.Getenv("DB_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort    = os.Getenv("DB_PORT") // e.g. '3306'
		dbName    = os.Getenv("DB_NAME") // e.g. 'my-database'
		// tls       = os.Getenv("DB_TLS")

		// connection pool
		maxOpenConnsStr = os.Getenv("DB_MAXOPEN_CONNS") //  e.g. 20
		// maxIdleConnsStr       = os.Getenv("DB_MAXIDLE_CONNS")       //  e.g. 2
		maxIdleSecondsStr     = os.Getenv("DB_MAXIDLE_SECONDS")     //  e.g. 10
		maxLifeTimeSecondsStr = os.Getenv("DB_MAXLIFETIME_SECONDS") //  e.g. 1800
	)
	if dbName == "" {
		return nil, errors.New("DB_NAME empty")
	}
	if dbTCPHost == "" {
		dbTCPHost = "0.0.0.0"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		return nil, errors.New("DB_USER empty")
	}
	if dbPwd == "" {
		return nil, errors.New("DB_PASS empty")
	}
	maxOpenConns := 20
	// maxIdleConns := 2
	maxIdleSeconds := 10
	maxLifeTimeSeconds := 1800
	if maxOpenConnsStr != "" {
		maxOpenConns, err = strconv.Atoi(maxOpenConnsStr)
		if err != nil {
			return nil, fmt.Errorf("DB_MAXOPEN_CONNS: %w", err)
		}
	}
	// if maxIdleConnsStr != "" {
	// 	maxIdleConns, err = strconv.Atoi(maxIdleConnsStr)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("DB_MAXIDLE_CONNS: %w", err)
	// 	}
	// }
	if maxIdleSecondsStr != "" {
		maxIdleSeconds, err = strconv.Atoi(maxIdleSecondsStr)
		if err != nil {
			return nil, fmt.Errorf("DB_MAXIDLE_SECONDS: %w", err)
		}
	}

	if maxLifeTimeSecondsStr != "" {
		maxLifeTimeSeconds, err = strconv.Atoi(maxLifeTimeSecondsStr)
		if err != nil {
			return nil, fmt.Errorf("DB_MAXLIFETIME_SECONDS: %w", err)
		}
	}

	fmt.Println(ctx, "connecting to postgresql db db=%s user=%s host=%s port=%v", dbName, dbUser, dbTCPHost, dbPort)
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	dbPool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	// defer dbPool.Close()

	// if tls != "" {
	// 	dbURI = fmt.Sprintf("%s&tls=%s", dbURI, tls)
	// }
	dbPool.Config().MaxConns = int32(maxOpenConns)
	dbPool.Config().MaxConnIdleTime = time.Duration(maxIdleSeconds) * time.Second
	dbPool.Config().MaxConnLifetime = time.Duration(maxLifeTimeSeconds) * time.Second
	// db.SetMaxOpenConns(maxOpenConns)
	// db.SetMaxIdleConns(maxIdleConns)
	// db.SetConnMaxIdleTime(time.Duration(maxIdleSeconds) * time.Second)
	// db.SetConnMaxLifetime(time.Duration(maxLifeTimeSeconds) * time.Second)

	// err = db.PingContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}
	return dbPool, nil
}
