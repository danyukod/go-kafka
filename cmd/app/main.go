package main

import (
	"database/sql"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/danyukod/go-kafka/internal/infra/akafka"
	"github.com/danyukod/go-kafka/internal/infra/repository"
	"github.com/danyukod/go-kafka/internal/infra/web"
	"github.com/danyukod/go-kafka/internal/usecase"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := repository.NewProductRepositoryMysql(db)
	createProductUsecase := usecase.NewCreateProductUseCase(repository)
	listProductsUsecase := usecase.NewListProductUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUsecase, listProductsUsecase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductHandler)

	go http.ListenAndServe(":8000", r)

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"product"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			// logar o erro
		}
		_, err = createProductUsecase.Execute(dto)
	}

}
