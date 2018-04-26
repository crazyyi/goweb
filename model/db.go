package model

type db interface {
	SelectPeople() ([]*Person, error)
	CreateNewRecord(*Person) (int, error)
	DeleteRowAt(id int64) (int64, error)
}