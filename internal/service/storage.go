package service

import (
	"sort"
	"sync"
)

type TaskStorage struct {
	mu      sync.RWMutex
	tasks   map[string]*Task
	counter int // Счетчик для определения Position
}

func NewStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]*Task),
	}
}

// сохраняет задачу и присваивает ей номер в очереди
func (s *TaskStorage) Save(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.Position == 0 { // Присваиваем номер только новым задачам
		s.counter++
		task.Position = s.counter
	}
	s.tasks[task.ID] = task
}

// Get возвращает задачу по ID
func (s *TaskStorage) Get(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

// GetAll возвращает все задачи, отсортированные по номеру в очереди
func (s *TaskStorage) GetAll() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	// Сортируем по полю Position
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Position < tasks[j].Position
	})

	return tasks
}

// Delete удаляет задачу по ID
func (s *TaskStorage) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}
