package model

import gonanoid "github.com/matoous/go-nanoid"

type Object struct {
	Id     string
	Value  string
	Parent string

	Active bool
}

func NewObject(value string, parent string) Object {
	//TODO по хорошему надо вынести, оставил здесь тк задание тестовое
	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuwxyz1234567890", 15)

	return Object{
		Id:     id,
		Value:  value,
		Parent: parent,
		Active: true,
	}
}
