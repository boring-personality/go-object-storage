package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ObjectStore interface {
	Insert(obj Object) 	error
	Read(id string)	(*Object, error)
	GetChunksFromDB(id string) ([]ChunkMetadata, error)
}

type Database struct{
	DB *sql.DB
}

type Object struct{
	Id 				string
	Original_name 	string
	Size			int64
	Created_at		string
	Chunks			[]ChunkMetadata
}

type ChunkMetadata struct {
	FileID			string
	Index			int
	Hash			string
	Path			string
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
	size BIGINT,
	created_at TIMESTAMP DEFAULT now()
	);`

	_, err := d.DB.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS chunkmetadata (
	id		UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	fileid 	TEXT REFERENCES objects(id) ON DELETE CASCADE,
	index	INT,
	hash	TEXT,
	path	TEXT,
	created_at TIMESTAMP DEFAULT now()
	);`
	_, err = d.DB.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Table created successfully!!!")
	return nil
}

func (d *Database) Insert(obj Object) error{
	query := `INSERT into objects (
	id, original_name, size) VALUES
	($1, $2, $3)`

	_, err := d.DB.Exec(query, obj.Id, obj.Original_name, obj.Size)
	if err != nil {
		log.Println(err)
		return err
	}

	query = `INSERT into chunkmetadata (
	fileid, index, hash, path) VALUES
	($1, $2, $3, $4)`

	for i := range(len(obj.Chunks)) {
		chunk := obj.Chunks[i]
		_, err = d.DB.Exec(query, chunk.FileID, chunk.Index, chunk.Hash, chunk.Path)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (d *Database) Read(id string) (*Object, error){
	query := `SELECT id, original_name, size, created_at
	from objects
	WHERE id = $1`

	row := d.DB.QueryRow(query, id)
	var obj Object
	err := row.Scan(
		&obj.Id,
		&obj.Original_name,
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

func (d *Database) GetChunksFromDB(id string) ([]ChunkMetadata, error) {
	query := `SELECT index, fileid, path
			  FROM chunkmetadata
			  WHERE fileid = $1
			  ORDER BY index ASC`

	rows, err := d.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	var chunks []ChunkMetadata
	for rows.Next() {
		var chunk ChunkMetadata
		err = rows.Scan(&chunk.Index, &chunk.FileID, &chunk.Path)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}
