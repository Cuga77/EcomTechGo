package store

import (
	"EcomTechGo/internal/models"
	"context"
	"errors"
	"testing"
)

func TestStore_Create(t *testing.T) {
	s := New()
	todo := &models.Todo{
		Title:       "Test Task",
		Description: "Desc",
	}

	err := s.Create(context.Background(), todo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}

	if len(s.Todos) != 1 {
		t.Errorf("expected map length 1, got %d", len(s.Todos))
	}
}

func TestStore_GetByID(t *testing.T) {
	s := New()
	existingTodo := &models.Todo{Title: "Exists", Description: "D"}
	if err := s.Create(context.Background(), existingTodo); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	targetID := existingTodo.ID

	found, err := s.GetByID(context.Background(), targetID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if found.Title != "Exists" {
		t.Errorf("expected title 'Exists', got '%s'", found.Title)
	}

	_, err = s.GetByID(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_GetAll(t *testing.T) {
	s := New()
	if err := s.Create(context.Background(), &models.Todo{Title: "Task 1"}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := s.Create(context.Background(), &models.Todo{Title: "Task 2"}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	todos, err := s.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("expected 2 todos, got %d", len(todos))
	}
}

func TestStore_Update(t *testing.T) {
	s := New()
	todo := &models.Todo{Title: "Old Title"}
	if err := s.Create(context.Background(), todo); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	id := todo.ID

	newTodo := &models.Todo{Title: "New Title", Description: "Updated"}
	err := s.Update(context.Background(), id, newTodo)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	saved, _ := s.GetByID(context.Background(), id)
	if saved.Title != "New Title" {
		t.Errorf("expected title updated to 'New Title', got '%s'", saved.Title)
	}

	err = s.Update(context.Background(), 999, newTodo)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound for invalid ID, got %v", err)
	}
}

func TestStore_Delete(t *testing.T) {
	s := New()
	todo := &models.Todo{Title: "To Delete"}
	if err := s.Create(context.Background(), todo); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	id := todo.ID

	err := s.Delete(context.Background(), id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(s.Todos) != 0 {
		t.Errorf("expected map to be empty, got length %d", len(s.Todos))
	}

	err = s.Delete(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_Create_DuplicateID(t *testing.T) {
	s := New()
	todo1 := &models.Todo{Title: "First", ID: 1}
	if err := s.Create(context.Background(), todo1); err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	todo2 := &models.Todo{Title: "Second", ID: 1}
	if err := s.Create(context.Background(), todo2); err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	if todo1.ID == todo2.ID {
		t.Errorf("IDs should be unique, got duplicate %d", todo1.ID)
	}
	if todo2.ID != 2 {
		t.Errorf("expected second ID to be 2, got %d", todo2.ID)
	}
}
