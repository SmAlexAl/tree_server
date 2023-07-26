package model

import gonanoid "github.com/matoous/go-nanoid"

const ACTIVE_STATE = ""
const DELETE_STATE = " (delete)"
const UNKNOW_STATE = " (state object unknow)"

type Object struct {
	Id     string
	Value  string
	Parent string

	//TODO посмотреть как убирать поле, либо для вывода сделать другое
	Active bool
	State  string
}

func NewObject(value string, parent string) Object {
	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuwxyz1234567890", 15)

	state := ACTIVE_STATE

	return Object{
		Id:     id,
		Value:  value,
		Parent: parent,
		Active: true,
		State:  state,
	}
}
