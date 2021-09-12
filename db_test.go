package go_neat

import (
	"fmt"
	_ "modernc.org/sqlite"
)

type Client struct {
	Id      int    `neat:"id,INT PRIMARY KEY"`
	Name    string `neat:"name,VARCHAR(100) NOT NULL"`
	Address string `neat:"address,TEXT NOT NULL"`
}

func ExampleNeatCreateTable() {
	db, err := Open("sqlite", getTempPath("neat.db"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	model := Client{Id: 1, Name: "John Doe"}

	if _, err = db.NeatCreateTable(model); err != nil {
		panic(err)
	}

	if _, err := db.NeatInsert(model); err != nil {
		panic(err)
	}

	factory := func() interface{} {
		return &Client{}
	}

	if result, err := db.NeatSelectOne("SELECT id, name FROM CLIENT WHERE id=?", factory, 1); err == nil {
		fmt.Println(result.(*Client))
	} else {
		fmt.Println(err)
	}
	// Output: &{1 John Doe }
}

func ExampleNeatUpdate() {
	db, err := Open("sqlite", getTempPath("neat.db"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	model := Client{Id: 1, Name: "John Doe"}

	if _, err = db.NeatCreateTable(model); err != nil {
		panic(err)
	}

	if _, err := db.NeatInsert(model); err != nil {
		panic(err)
	}

	factory := func() interface{} {
		return &Client{}
	}

	model.Name = "Mary Jane"
	if _, err := db.NeatUpdate(model, "Id", model.Id); err == nil {
		if result, err := db.NeatSelectOne("SELECT id, name FROM CLIENT WHERE id=?", factory, 1); err == nil {
			fmt.Println(result.(*Client))
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	// Output: &{1 Mary Jane }
}