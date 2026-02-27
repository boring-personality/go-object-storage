package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ObjectStore interface {
	Insert(obj Object) 	error
	Read(id string)	(*Object, error)
}

type Database struct{
	DB *sql.DB
}

type Object struct{
	Id 				string
	Original_name 	string
	Disk_path 		string
	Size			int64
	Created_at		string
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("pgx", CONN_STRING)

	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

const CONN_STRING = "postgres://root:root@localhost:5432/objectstore"

func (d *Database) Create() error{
	query := `CREATE TABLE IF NOT EXISTS objects (
	id TEXT PRIMARY KEY,
	original_name TEXT,
	disk_path TEXT NOT NULL,
	size BIGINT,
	created_at TIMESTAMP DEFAULT now()
	);`

	_, err := d.DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Table created successfully!!!")
	return nil
}

func (d *Database) Insert(obj Object) error{
	query := `INSERT into objects (
	id, original_name, disk_path, size) VALUES
	($1, $2, $3, $4)`

	_, err := d.DB.Exec(query, obj.Id, obj.Original_name, obj.Disk_path, obj.Size)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Read(id string) (*Object, error){
	query := `SELECT id, original_name, disk_path, size, created_at
	from objects
	WHERE id = $1`

	row := d.DB.QueryRow(query, id)
	var obj Object
	err := row.Scan(
		&obj.Id,
		&obj.Original_name,
		&obj.Disk_path,
		&obj.Size,
		&obj.Created_at,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}
