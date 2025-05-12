package domain

import "fmt"

type ErrAuthorNotFound struct {
	AuthorID int
}

func (e *ErrAuthorNotFound) Error() string {
	return fmt.Sprintf("author with ID %d not found", e.AuthorID)
}

type ErrBookNotFound struct {
	BookID int
}

func (e *ErrBookNotFound) Error() string {
	return fmt.Sprintf("book with ID %d not found", e.BookID)
}
