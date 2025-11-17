package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-worker-service/internal/service"

	"github.com/google/uuid"
)

// APIHandler обрабатывает HTTP-запросы
type APIHandler struct {
	storage   *service.TaskStorage
	taskQueue chan<- string // Канал только для записи (отправки задач воркерам)
}

// NewAPIHandler создает новый обработчик
func NewAPIHandler(storage *service.TaskStorage, taskQueue chan<- string) *APIHandler {
	return &APIHandler{
		storage:   storage,
		taskQueue: taskQueue,
	}
}

// RegisterRoutes регистрирует эндпоинты в роутере
func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks)
}

func (h *APIHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// DTO для запроса на создание задачи
type createTaskRequest struct {
	N   int     `json:"n"`  // Сколько раз нужно добавить
	D   float64 `json:"d"`  // Число, которое добавляем
	N1  float64 `json:"n1"` // Стартовое число
	I   float64 `json:"I"`  // Интервал
	TTL float64 `json:"TTL"`
}

// @Summary      Создать новую задачу
// @Description  Принимает параметры для арифметической прогрессии и ставит задачу в очередь.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      createTaskRequest  true  "Параметры задачи"
// @Success      201  {object}  map[string]string  "Возвращает ID созданной задачи"
// @Failure      400  {object}  map[string]string  "Ошибка в теле запроса"
// @Router       /tasks [post]
func (h *APIHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task := &service.Task{
		ID:        uuid.New().String(),
		Status:    service.StatusQueued,
		N:         req.N,
		D:         req.D,
		N1:        req.N1,
		I:         req.I,
		TTL:       req.TTL,
		CreatedAt: time.Now().UTC(),
		Result:    make([]float64, 0, req.N),
	}

	h.storage.Save(task)

	// Отправляем ID задачи в очередь на обработку
	h.taskQueue <- task.ID
	log.Printf("Task %s has been queued", task.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": task.ID})
}

// @Summary      Получить список всех задач
// @Description  Возвращает отсортированный список всех задач и их статусы.
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   service.Task       "Список задач"
// @Router       /tasks [get]
func (h *APIHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.storage.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
