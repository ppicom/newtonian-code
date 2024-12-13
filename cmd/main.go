package main

import (
	"database/sql"
	"log"

	usecases "github.com/ppicom/newtonian/internal/application/use_cases"
	"github.com/ppicom/newtonian/internal/infrastructure/api/http"
	"github.com/ppicom/newtonian/internal/infrastructure/db"
	"github.com/redis/go-redis/v9"
)

func main() {
	router := http.NewRouter()
	router.Engine().Run(":8080")

	// Initialize database connections
	conn, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/banking")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Initialize repositories and controllers
	accountRepo := db.NewAccountRepository(conn, rdb)
	transferMoneyUseCase := usecases.NewTransferMoneyUseCase(accountRepo)
	controller := http.NewController(transferMoneyUseCase)

	// Setup routes
	controller.SetupRoutes(router)
}
