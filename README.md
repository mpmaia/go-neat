# go-neat library

This project implements a very simple Micro ORM for Go inspired by the [.NET Dapper Library](https://github.com/DapperLib/Dapper).

I've developed this library while learning about how to use reflection on Go. At this moment, this code is still an unstable work in progress. 

# Usage

Create a struct with the `neat` field tag describing the column name and the database column declaration. For example:

```go
type Client struct {
    Id      int    `neat:"id,INT PRIMARY KEY"`
    Name    string `neat:"name,VARCHAR(100) NOT NULL"`
    Address string `neat:"address,TEXT NOT NULL"`
}
```

Call the `neat.Open` method to get a `sql.DB` handle and use the methods prefixed by `Neat` to make the database operations:

```go

db, err := go_neat.Open("sqlite", getTempPath("neat.db"))
if err != nil {
    panic(err)
}
defer db.Close()

model := Client{Id: 1, Name: "John Doe"}

// This will create the table CLIENT with columns id, name and address
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
    fmt.Println(result.(*Client)) // Output: &{1 John Doe }
}
```