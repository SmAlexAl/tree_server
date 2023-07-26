package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.sqlite.new"

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	st := Storage{db: db}

	err = st.initDb()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &st, nil
}

func (s Storage) GetTree() (map[string]model.Object, error) {
	stmt, err := s.db.Prepare(`SELECT id, leaf, parentId, active FROM tree`)

	if err != nil {
		return nil, err
	}

	res := make(map[string]model.Object)
	var parentId interface{}

	rows, err := stmt.Query()

	for rows.Next() {
		ob := model.Object{}

		err = rows.Scan(&ob.Id, &ob.Value, &parentId, &ob.Active)

		if err != nil {
			return nil, err
		}

		if parentId != nil {
			ob.Parent = parentId.(string)
		}

		res[ob.Id] = ob
	}

	return res, nil
}

func (s *Storage) GetLeaf(id string) (model.Object, error) {
	stmt, err := s.db.Prepare(`SELECT id, leaf, parentId, active FROM tree WHERE id = ? and active = true`)

	if err != nil {
		return model.Object{}, err
	}

	//подумать над адаптером
	var ob model.Object
	var parentId interface{}

	err = stmt.QueryRow(id).Scan(&ob.Id, &ob.Value, &parentId, &ob.Active)

	if err != nil {
		return model.Object{}, err
	}

	if parentId != nil {
		ob.Parent = parentId.(string)
	}

	return ob, nil
}

func (s *Storage) SaveLeaf(object model.Object) error {
	stmt, err := s.db.Prepare(`INSERT INTO tree(id, leaf, parentId, active) VALUES (?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	if object.Parent == "" {
		_, err = stmt.Exec(object.Id, object.Value, nil, 1)
	} else {
		_, err = stmt.Exec(object.Id, object.Value, object.Parent, 1)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) UpdateLeaf(val model.Object) error {
	stmt, err := s.db.Prepare(`UPDATE tree SET leaf = ? WHERE id = ?;`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(val.Value, val.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) GetLeafsByActive(id string, active bool) error {
	stmt, err := s.db.Prepare(`
	WITH treeView AS (
    	select id, leaf, parentId, active FROM tree WHERE id = ?
    	UNION ALL
    	select t.id, t.leaf, t.parentId, t.active
    	FROM tree as t
    	JOIN treeView as tv ON tv.parentId = t.id
	) SELECT id From treeView WHERE active = ? LIMIT 1
	`)

	if err != nil {
		return err
	}

	var res string

	err = stmt.QueryRow(id, active).Scan(&res)

	if err != nil {
		return err
	}

	return nil

}

func (s Storage) DeleteLeaf(val model.Object) error {
	stmt, err := s.db.Prepare(`
	WITH treeView AS (
    	select id, leaf, parentId, active FROM tree WHERE id = ?
    	UNION ALL
    	select t.id, t.leaf, t.parentId, t.active
    	FROM tree as t
    	JOIN treeView as tv ON tv.id = t.parentId
	) UPDATE tree SET active = 0 FROM treeView WHERE tree.id = treeView.id;
	`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(val.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) initDb() error {
	const op = "storage.sqlite.new"

	stmt, err := s.db.Prepare(`
	CREATE TABLE IF NOT EXISTS tree(
		id STRING PRIMARY KEY,
		leaf STRING NOT NULL,
		active BOOLEAN NOT NULL,
		parentId STRING,
		FOREIGN KEY (parentId) REFERENCES tree (id));
	`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) SaveLeafs(fixtures []model.Object) error {
	sqlStr := "INSERT OR IGNORE INTO tree(id, leaf, parentId, active) VALUES "

	var rows []interface{}
	var p interface{}

	for _, row := range fixtures {
		sqlStr += "(?, ?, ?, ?),"

		p = row.Parent
		if row.Parent == "" {
			p = nil
		}

		rows = append(rows, row.Id, row.Value, p, row.Active)
	}
	sqlStr = sqlStr[:len(sqlStr)-1]
	stmt, err := s.db.Prepare(sqlStr)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(rows...)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) TruncateTree() error {
	stmt, err := s.db.Prepare("DELETE FROM tree")

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}
