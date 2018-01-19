package main

import (
    "errors"
    "strconv"
    "fmt"
    "database/sql"                    // https://golang.org/pkg/database/sql
    _ "github.com/mattn/go-sqlite3"
)

const (
    dbDriverName = "sqlite3"
    dbName = "answer42db"
    dbFileName = dbName+".sqlite3"
    ideaNameStrMaxLen = 32
)

func createDB(dbName string, dataSourceName string) (error) {

   //db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/")
   db, err := sql.Open(dbDriverName, dataSourceName)
   if err != nil {
       return err
   }
   defer db.Close()

   //_,err = db.Exec("CREATE DATABASE "+dbName)
   //if err != nil {
   //    panic(err)
   //}

   //_,err = db.Exec("USE "+dbName)
   //if err != nil {
   //    panic(err)
   //}

   createIdeaTable := "CREATE TABLE idea (" +
                      "id INTEGER PRIMARY KEY AUTOINCREMENT"+
                      ", name varchar("+strconv.Itoa(ideaNameStrMaxLen)+") NOT NULL UNIQUE"+
                      ", data varchar(32)"+
                      ")"

   _,err = db.Exec(createIdeaTable)
   if err != nil {
       panic(err)
   }

    return nil
}

func openDB(dbFilePath string) (*sql.DB, error) {
   return sql.Open(dbDriverName, dbFilePath)
}

func insertIdea(db *sql.DB, ideaName string) (int64, error) {
    if len(ideaName) > ideaNameStrMaxLen {
        return -1, errors.New(fmt.Sprintf("idea name is too long, %d > %d",
                              len(ideaName), ideaNameStrMaxLen))
    }

    stmt, err := db.Prepare("INSERT INTO idea(name, data) values(?,?)")
    if err != nil {
        return -1, err
    }

    res, err := stmt.Exec(/*name:*/ideaName, /*data:*/"test string")
    if err != nil {
        return -1, err
    }

    return res.LastInsertId()
}

func main() {
    dbFilePath := "./"+dbFileName
    if err := createDB(dbName, dbFilePath); err != nil {
        panic(err)
    }

    db, err := openDB(dbFilePath)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    insertIdea(db, "idea1")
    insertIdea(db, "idea2")
    _, err = insertIdea(db, "idea1")
    if err != nil {
            fmt.Println("can't insert into idea table:", err)
    }
}
