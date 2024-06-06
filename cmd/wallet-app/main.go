// This is the application entry point.
// All the wiring needed for the application
// to start is typically done here.
// This should be lean.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Oloruntobi1/grey/internal/config"
	"github.com/Oloruntobi1/grey/internal/db/migrations"
	db "github.com/Oloruntobi1/grey/internal/db/sqlc"
	"github.com/Oloruntobi1/grey/internal/repositories"
	"github.com/Oloruntobi1/grey/internal/transport/http/domains/users"
	"github.com/Oloruntobi1/grey/internal/transport/http/domains/wallets"
	"github.com/Oloruntobi1/grey/internal/transport/http/handlers"
	"github.com/Oloruntobi1/grey/pkg/logger"
	"github.com/Oloruntobi1/grey/pkg/metrics"
	"github.com/Oloruntobi1/grey/pkg/tracer"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func getAppEnv() (string, string, error) {
	env := os.Getenv("MY_ENV")
	var filename string

	switch env {
	case "":
		filename = ".env"
	case "development":
		filename = "dev.env"
	case "test":
		filename = "test.env"
	case "production":
		filename = "prod.env"
	default:
		return "", "", fmt.Errorf("invalid environment: %v", env)
	}

	return env, filename, nil
}

func main() {
	// We start with making sure the environment is correct.
	// Prevent app from starting if an error is encountered
	// in this process.
	env, envFile, err := getAppEnv()
	if err != nil {
		log.Fatal(err)
	}

	if env == "" {
		env = "local"
	}

	// Then we load the appropriate env file based on the environment
	// where the application is running.
	// Prevent app from starting if an error is also
	// encountered in this process.
	err = godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file for %s environment", envFile, env)
	}

	// We start a context that will be used throughout the application
	ctx := context.Background()

	// After that we initialize our traces and metrics
	// if the project will be utilizing it
	serviceName := "grey-wallet-application"
	tp, err := tracer.StartTracer(ctx, serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	mp, err := metrics.SetupMetrics(ctx, serviceName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Next thing is to connect to a database.
	// Could be any but in this example we will
	// be using postgres.
	dbCfg := config.GetDatabaseConfig()

	if envFile == ".env" {
		dbCfg = fmt.Sprintf("%s?sslmode=disable", dbCfg)
	}

	if envFile == "dev.env" {
		dbCfg = fmt.Sprintf("%s?sslmode=disable", dbCfg)
	}

	pgxConfig, err := pgxpool.ParseConfig(dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	pgxConfig.ConnConfig.Tracer = &MyQueryTracer{}
	connPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Now that we have obtained our connection pool
	// we can run our migrations if applicable

	// But firstly we need to initialize our logger.
	logger := logger.NewSlog(ctx)

	// Then we will run our migrations
	err = migrations.RunDBMigration(ctx, connPool, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Obtain all queries
	dbQueries := db.New(connPool)

	// Use queries to initiliaze repositories
	userRepository := repositories.NewUserRepository(dbQueries)
	walletRepository := repositories.NewWalletRepository(dbQueries)

	userService := users.NewUserService(userRepository)
	walletService := wallets.NewWalletService(walletRepository)

	userHandler := handlers.NewUserHandler(*userService, logger)
	walletHandler := handlers.NewWalletHandler(*walletService, logger)

	// TODO: attach the tracing to middleware

	router := handlers.SetupRouter(ctx, *userHandler, *walletHandler)

	log.Fatal(http.ListenAndServe(":9191", router))
}

type MyQueryTracer struct{}

func (t *MyQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, d pgx.TraceQueryStartData) context.Context {
	ctx, span := otel.Tracer("").Start(ctx, "execute_query")
	defer span.End()

	span.SetAttributes(attribute.String("sql", d.SQL))
	argStrings := make([]string, len(d.Args))
	for i, arg := range d.Args {
		argStrings[i] = fmt.Sprintf("%v", arg) // Use fmt.Sprintf to handle any type
	}
	span.SetAttributes(attribute.StringSlice("db.statement.parameters", argStrings))

	return trace.ContextWithSpan(ctx, span)
}

func (t *MyQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, d pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx) // Assuming you're using OpenTelemetry
	if d.Err != nil {
		span.SetStatus(codes.Error, "database error")
		span.RecordError(d.Err)
	} else {
		span.AddEvent("successfully executed SQL statement")
	}
	span.End()
}
