package db

import "testing"

func setupDatabase(t *testing.T) *Database {
	t.Helper()

	database, err := NewDatabase()
	if err != nil {
		t.Fatal(err)
	}

	// _, err = database.DB.Exec("DROP TABLE IF EXISTS objects")
	if err != nil {
		t.Fatal(err)
	}

	err = database.Create()
	if err != nil {
		t.Fatal(err)
	}

	return database
}

func TestInsertAndRead(t *testing.T) {
	database := setupDatabase(t)

	obj := Object {
		Id: 			"1",
		Original_name:	"abc.mp4",
		Size:			15,
	}

	chunk := ChunkMetadata {
		FileID:		"1",
		Index:		1,
		Hash:		"testhash",
		Path:		"somepath",
	}

	obj.Chunks = append(obj.Chunks, chunk)

	err := database.Insert(obj)
	if err != nil {
		t.Fatal(err)
	}

	read_obj, err := database.Read("1")
	if err != nil {
		t.Fatal(err)
	}

	obj_chunks, err := database.GetChunksFromDB("1")
	if err != nil {
		t.Fatal(err)
	}

	if len(obj_chunks) != 1 {
		t.Fatal("could not retrive all the chunks")
	}
	if obj_chunks[0].Path != "somepath" {
		t.Fatal("Failed to retrive the path for chunks")
	}

	if read_obj.Id != obj.Id {
		t.Fatal("IDs do not match to the inserted value")
	}
	defer database.DB.Close()
}

func TestReadNegative(t *testing.T) {
	database := setupDatabase(t)

	obj, err := database.Read("nonexistentid")
	if err != nil {
		t.Fatal(err)
	}
	if obj != nil {
		t.Fatal("Expected nil got value", obj)
	}
	defer database.DB.Close()
}

