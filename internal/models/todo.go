package models

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Completed   bool   `json:"completed"`
}
