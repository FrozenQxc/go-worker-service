package main

import (
	"flag"
	"log"
	"net/http"

	"go-worker-service/internal/api"
	"go-worker-service/internal/service"

	_ "go-worker-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Worker Service API
// @version 1.0
// @description Сервис для расчета арифметической прогрессии в фоновом режиме.
// @host localhost:8080
// @BasePath /
func main() {
	// Конфигурация
	numWorkers := flag.Int("N", 5, "Количество параллельных воркеров")
	flag.Parse()

	//  Создание зависимостей
	storage := service.NewStorage()

	taskQueue := make(chan string, 100) // Буфер на 100 задач

	workerPool := service.NewWorkerPool(*numWorkers, taskQueue, storage)
	workerPool.Run()

	apiHandler := api.NewAPIHandler(storage, taskQueue)

	// 	Настройка роутера и запуск сервера
	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is starting on :8080")
	log.Printf("Swagger UI is available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
