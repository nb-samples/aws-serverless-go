package sample

import "time"

// Item structure
type Item struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Details   Details    `json:"details,omitempty"`
}

// Details of the item
type Details struct {
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
}
