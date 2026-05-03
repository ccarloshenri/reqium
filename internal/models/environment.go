package models

import "time"

type Environment struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
	Active    bool              `json:"active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
