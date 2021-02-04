package db

// Book models the book table in the database
type Book struct {
	Name string `json:"name"`
	Author string `json:"author"`
	Genre string `json:"genre"`
}