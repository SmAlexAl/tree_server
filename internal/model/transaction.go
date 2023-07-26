package model

const ADD = "add"
const DELETE = "delete"
const UPDATE = "update"

type Transaction struct {
	Command string
	Object  Object
}

func NewTransaction(command string, ob Object) Transaction {
	return Transaction{
		Command: command,
		Object:  ob,
	}
}
