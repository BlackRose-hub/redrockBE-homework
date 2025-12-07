package models

import "time"

type Course struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Teacher   string    `json:"teacher"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
}
