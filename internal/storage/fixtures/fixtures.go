package fixtures

import "github.com/SmAlexAl/tree_server.git/internal/model"

type Fixtures struct {
	Collection []model.Object
}

func New() Fixtures {
	return Fixtures{
		Collection: []model.Object{
			{
				Id:     "a",
				Value:  "node1",
				Parent: "",
				Active: true,
			},
			{
				Id:     "b",
				Value:  "node2",
				Parent: "a",
				Active: true,
			},
			{
				Id:     "c",
				Value:  "node3",
				Parent: "a",
				Active: true,
			},
			{
				Id:     "d",
				Value:  "node4",
				Parent: "b",
				Active: true,
			},
			{
				Id:     "e",
				Value:  "node5",
				Parent: "d",
				Active: true,
			},
			{
				Id:     "tt",
				Value:  "node6",
				Parent: "e",
				Active: true,
			},
			{
				Id:     "yy",
				Value:  "node7",
				Parent: "tt",
				Active: true,
			},
			{
				Id:     "uu",
				Value:  "node8",
				Parent: "yy",
				Active: true,
			},
			{
				Id:     "ttt",
				Value:  "node9",
				Parent: "c",
				Active: true,
			},
			{
				Id:     "yyy",
				Value:  "node10",
				Parent: "ttt",
				Active: true,
			},
			{
				Id:     "uuu",
				Value:  "node11",
				Parent: "yyy",
				Active: true,
			},
		},
	}
}

func (f Fixtures) GetCollection() []model.Object {
	return f.Collection
}
