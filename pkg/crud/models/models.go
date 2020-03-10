package models

// models - описание объектов предментной области
// https://github.com/golang/go/wiki/CodeReviewComments
type Burger struct {
	Id int64
	Name string
	Price int
	Removed bool
	FileName string
}

