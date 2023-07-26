package model

import gonanoid "github.com/matoous/go-nanoid"

type Object struct {
	Id     string
	Value  string
	Parent string

	//TODO посмотреть как убирать поле, либо для вывода сделать другое
	Active bool
}

func NewObject(value string, parent string) Object {
	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuwxyz1234567890", 15)

	return Object{
		Id:     id,
		Value:  value,
		Parent: parent,
		Active: true,
	}
}
