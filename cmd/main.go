package main

import (
	"fmt"
	"log"
	"net/http"

	_ "go-worker-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Worker Service API
// @version 1.0
// @description Тестовый сервис.
// @host localhost:8080
// @BasePath /

// testHandler демонстрирует простой эндпоинт
// @Summary Тестовый эндпоинт
// @Description Возвращает подтверждение работы сервиса
// @Tags Test
// @Success 200 {string} string "ok"
// @Router /test [get]
func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service works!")
}

func main() {
	http.HandleFunc("/test", testHandler)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	log.Println("Swagger UI: http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
