package service

import (
	"log"
	"time"
)

type WorkerPool struct {
	numWorkers int
	taskQueue  <-chan string // Канал только для чтения
	storage    *TaskStorage
}

func NewWorkerPool(numWorkers int, taskQueue <-chan string, storage *TaskStorage) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		taskQueue:  taskQueue,
		storage:    storage,
	}
}

func (p *WorkerPool) Run() {
	log.Printf("Starting %d workers...", p.numWorkers)
	for i := 0; i < p.numWorkers; i++ {
		go p.worker(i + 1)
	}

	// Запускаем "сборщик мусора" для TTL в отдельной горутине
	go p.ttlManager()
}

func (p *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)
	for taskID := range p.taskQueue {
		log.Printf("Worker %d received task %s", id, taskID)

		task, ok := p.storage.Get(taskID)
		if !ok {
			log.Printf("Worker %d: task %s not found", id, taskID)
			continue
		}

		// Обновляем статус и время начала.
		now := time.Now().UTC()
		task.Status = StatusInProgress
		task.StartedAt = &now
		task.Result = append(task.Result, task.N1)
		p.storage.Save(task)

		ticker := time.NewTicker(time.Duration(task.I * float64(time.Second)))

		// Основной цикл расчета прогрессии
		for i := 1; i < task.N; i++ {
			<-ticker.C // Ждем сигнала от тикера

			lastValue := task.Result[len(task.Result)-1]
			nextValue := lastValue + task.D

			task.Result = append(task.Result, nextValue)
			task.CurrentIter = i + 1
			p.storage.Save(task)
			log.Printf("Worker %d: task %s, iteration %d, value %.2f", id, taskID, task.CurrentIter, nextValue)
		}

		ticker.Stop()

		finishedTime := time.Now().UTC()
		task.Status = StatusCompleted
		task.FinishedAt = &finishedTime
		p.storage.Save(task)

		log.Printf("Worker %d finished task %s", id, taskID)
	}
}

// ttlManager периодически проверяет и удаляет старые завершенные задачи
func (p *WorkerPool) ttlManager() {
	ticker := time.NewTicker(10 * time.Second) // Проверяем каждые 10 секунд
	defer ticker.Stop()

	for range ticker.C {
		tasks := p.storage.GetAll()
		for _, task := range tasks {
			if task.Status == StatusCompleted && task.FinishedAt != nil {
				// Проверяем, истекло ли время TTL
				if time.Since(*task.FinishedAt).Seconds() > task.TTL {
					log.Printf("TTL expired for task %s. Deleting.", task.ID)
					p.storage.Delete(task.ID)
				}
			}
		}
	}
}
