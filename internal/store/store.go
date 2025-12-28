package store

import (
	"EcomTechGo/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Store struct {
	Todos    map[int]models.Todo
	mu       sync.Mutex
	NextID   int
	FilePath string
}

var ErrNotFound = errors.New("todo not found")

func New() *Store {
	return &Store{
		Todos: make(map[int]models.Todo),
	}
}

func (s *Store) save() error {
	if s.FilePath == "" {
		return nil
	}
	// Lock is already held by caller
	data, err := json.MarshalIndent(s.Todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0644)
}

func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.save()
}

func (s *Store) Load(filename string) error {
	s.FilePath = filename
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := json.Unmarshal(data, &s.Todos); err != nil {
		return err
	}
	// Recalculate NextID
	maxID := 0
	for id := range s.Todos {
		if id > maxID {
			maxID = id
		}
	}
	s.NextID = maxID
	return nil
}

func (s *Store) Create(ctx context.Context, todo *models.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if todo == nil {
		return errors.New("todo cannot be nil")
	}

	s.NextID++
	todo.ID = s.NextID
	s.Todos[s.NextID] = *todo

	return s.save()
}

func (s *Store) GetAll(ctx context.Context) ([]models.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newSlice := make([]models.Todo, 0)
	for i := range s.Todos {
		newSlice = append(newSlice, s.Todos[i])
	}

	return newSlice, nil
}

func (s *Store) GetByID(ctx context.Context, ID int) (models.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.Todos[ID]
	if !ok {
		return models.Todo{}, fmt.Errorf("cannot get todo with id %d: %w", ID, ErrNotFound)
	}
	return todo, nil
}

func (s *Store) Delete(ctx context.Context, ID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Todos[ID]; !ok {
		return fmt.Errorf("cannot delete todo with id %d: %w", ID, ErrNotFound)
	}

	delete(s.Todos, ID)
	return s.save()
}

func (s *Store) Update(ctx context.Context, ID int, todo *models.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if todo == nil {
		return errors.New("update data cannot be nil")
	}

	if _, ok := s.Todos[ID]; !ok {
		return fmt.Errorf("cannot update todo with id %d: %w", ID, ErrNotFound)
	}

	todo.ID = ID
	s.Todos[ID] = *todo
	return s.save()
}
