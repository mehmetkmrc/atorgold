package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool
var DBErr error
var url = "postgres://postgres:ITM-2020@localhost:5432/ator_gold?sslmode=disable"

func InitiliazeDatabaseConnection() {

	DBPool, DBErr = pgxpool.New(context.Background(), url)
	if DBErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", DBErr)
		os.Exit(1)
	}
}
