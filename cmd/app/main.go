package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"my_app/internal/handler"
	"my_app/internal/repository"
	"my_app/internal/service"
)

func main() {
	repos := repository.NewRepository(establishConnection())
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	engine := handlers.InitRoutes()

	if err := engine.Run(":8081"); err != nil {
		panic(err)
	}

}

func establishConnection() *sqlx.DB {
	conn, err := sqlx.Connect("postgres", "host=localhost port=5433 user=user dbname=app_database password=password sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err := conn.Ping(); err != nil {
		panic(err)
	}

	return conn
}
