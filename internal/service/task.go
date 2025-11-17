package service

import "time"

// Константы для статусов задач
const (
	StatusQueued     = "In Queue"
	StatusInProgress = "In Progress"
	StatusCompleted  = "Completed"
	StatusFailed     = "Failed"
)

type Task struct {
	ID          string    `json:"id"`
	Position    int       `json:"position"` // Номер в очереди для сортировки
	Status      string    `json:"status"`
	N           int       `json:"n"`
	D           float64   `json:"d"`
	N1          float64   `json:"n1"`
	I           float64   `json:"I"`   // Интервал в секундах
	TTL         float64   `json:"TTL"` // Время жизни в секундах
	CurrentIter int       `json:"current_iteration"`
	Result      []float64 `json:"result,omitempty"` // omitempty скроет поле, если оно пустое

	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}
