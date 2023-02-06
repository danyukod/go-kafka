package main

import (
	"database/sql"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/danyukod/go-kafka/internal/infra/akafka"
	"github.com/danyukod/go-kafka/internal/infra/repository"
	"github.com/danyukod/go-kafka/internal/usecase"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	msgChan := make(chan *kafka.Message)
	go akafka.Consumer([]string{"products"}, "host.docker.internal:9094", msgChan)

	repository := repository.NewProductRepositoryMysql(db)
	createProductUseCase := usecase.NewCreateProductUsecase(repository)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			continue
		}
		createProductUseCase.Execute(dto)
	}

}
