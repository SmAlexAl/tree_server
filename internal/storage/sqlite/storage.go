package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB

	tx *sql.Tx
}

const PREPARE_ERROR = "prepare error"
const SCAN_ERROR = "scan error"
const EXEC_ERROR = "exec error"

// TODO возвращается полностью собранный объект, либо список объектов. В некоторых местах есть валидация.
// Обычно создание объекта и валидацию нужно выносить в прослойку между бд и бизнесом
// сделал так, потому что задание тестовое, для ускорения разработки
func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	st := Storage{db: db}

	err = st.initDb()

	if err != nil {
		return nil, fmt.Errorf("init error: %w", err)
	}

	return &st, nil
}

func (s *Storage) BeginTransaction() error {
	tx, err := s.db.Begin()

	s.tx = tx

	return err
}

func (s *Storage) Commit() error {
	if s.tx != nil {
		err := s.tx.Commit()

		if err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		s.tx = nil
	} else {
		return fmt.Errorf("transaction doesnt init")
	}

	return nil
}

func (s *Storage) Rollback() error {
	if s.tx != nil {
		err := s.tx.Rollback()

		if err != nil {
			return fmt.Errorf("rollback error: %w", err)
		}
		s.tx = nil
	} else {
		return fmt.Errorf("transaction doesnt init")
	}

	return nil
}

func (s Storage) prepare(sql string) (*sql.Stmt, error) {
	if s.tx != nil {
		return s.tx.Prepare(sql)
	} else {
		return s.db.Prepare(sql)
	}
}

func (s Storage) GetTree() (map[string]model.Object, error) {
	stmt, err := s.prepare(`SELECT id, leaf, parentId, active FROM tree`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	res := make(map[string]model.Object)
	var parentId *string

	rows, err := stmt.Query()

	for rows.Next() {
		ob := model.Object{}

		err = rows.Scan(&ob.Id, &ob.Value, &parentId, &ob.Active)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", SCAN_ERROR, err)
		}

		if parentId != nil {
			ob.Parent = *parentId
		}

		res[ob.Id] = ob
	}

	return res, nil
}

func (s *Storage) GetLeaf(id string) (model.Object, error) {
	stmt, err := s.prepare(`SELECT id, leaf, parentId, active FROM tree WHERE id = ? and active = true`)

	if err != nil {
		return model.Object{}, fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	var ob model.Object
	var parentId *string

	err = stmt.QueryRow(id).Scan(&ob.Id, &ob.Value, &parentId, &ob.Active)

	if err != nil {
		return model.Object{}, fmt.Errorf("%s: %w", SCAN_ERROR, err)
	}

	if parentId != nil {
		ob.Parent = *parentId
	}

	return ob, nil
}

func (s *Storage) SaveLeaf(object model.Object) error {
	stmt, err := s.prepare(`INSERT INTO tree(id, leaf, parentId, active) VALUES (?, ?, ?, ?)`)

	if err != nil {
		return fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	if object.Parent == "" {
		_, err = stmt.Exec(object.Id, object.Value, nil, 1)
	} else {
		_, err = stmt.Exec(object.Id, object.Value, object.Parent, 1)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", EXEC_ERROR, err)
	}

	return nil
}

func (s Storage) UpdateLeaf(val model.Object) error {
	stmt, err := s.prepare(`UPDATE tree SET leaf = ? WHERE id = ?;`)

	if err != nil {
		return fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	_, err = stmt.Exec(val.Value, val.Id)

	if err != nil {
		return fmt.Errorf("%s: %w", EXEC_ERROR, err)
	}

	return nil
}

func (s Storage) GetLeafsByActive(id string, active bool) (bool, error) {
	stmt, err := s.prepare(`
	WITH treeView AS (
    	select id, leaf, parentId, active FROM tree WHERE id = ?
    	UNION ALL
    	select t.id, t.leaf, t.parentId, t.active
    	FROM tree as t
    	JOIN treeView as tv ON tv.parentId = t.id
	) SELECT id From treeView WHERE active = ? LIMIT 1
	`)

	if err != nil {
		return false, fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	var res string

	err = stmt.QueryRow(id, active).Scan(&res)

	//TODO решение не оптимальное, принято такое решение, чтобы сократить время выполнения задания
	if err != nil && err.Error() == "sql: no rows in result set" {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("%s: %w", SCAN_ERROR, err)
	}

	return true, nil

}

func (s Storage) DeleteLeaf(val model.Object) error {
	stmt, err := s.prepare(`
	WITH treeView AS (
    	select id, leaf, parentId, active FROM tree WHERE id = ?
    	UNION ALL
    	select t.id, t.leaf, t.parentId, t.active
    	FROM tree as t
    	JOIN treeView as tv ON tv.id = t.parentId
	) UPDATE tree SET active = 0 FROM treeView WHERE tree.id = treeView.id;
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	_, err = stmt.Exec(val.Id)

	if err != nil {
		return fmt.Errorf("%s: %w", EXEC_ERROR, err)
	}

	return nil
}

func (s *Storage) initDb() error {
	stmt, err := s.prepare(`
	CREATE TABLE IF NOT EXISTS tree(
		id STRING PRIMARY KEY,
		leaf STRING NOT NULL,
		active BOOLEAN NOT NULL,
		parentId STRING,
		FOREIGN KEY (parentId) REFERENCES tree (id));
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return fmt.Errorf("%s: %w", EXEC_ERROR, err)
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
	stmt, err := s.prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	_, err = stmt.Exec(rows...)

	if err != nil {
		return fmt.Errorf("%s: %w", EXEC_ERROR, err)
	}

	return nil
}

func (s *Storage) TruncateTree() error {
	stmt, err := s.prepare("DELETE FROM tree")

	if err != nil {
		return fmt.Errorf("%s: %w", PREPARE_ERROR, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return fmt.Errorf("%s: %w", EXEC_ERROR, err)
	}

	return nil
}
