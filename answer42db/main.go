package main

import (
    "errors"
    "strings"
    "strconv"
    "fmt"
    "sync"
    "os"
    "io/ioutil"
    "path/filepath"
    "log"
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
                      ", name VARCHAR("+strconv.Itoa(ideaNameStrMaxLen)+") NOT NULL UNIQUE"+
                      ", json JSON HIDDEN"+
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

func insertIdea(db *sql.DB, ideaName string, jsonText string) (int64, error) {
    if len(ideaName) > ideaNameStrMaxLen {
        return -1, errors.New(fmt.Sprintf("idea name is too long, %d > %d",
                              len(ideaName), ideaNameStrMaxLen))
    }

    stmt, err := db.Prepare("INSERT INTO idea(name, json) values(?,?)")
    if err != nil {
        return -1, err
    }

    res, err := stmt.Exec(/*name:*/ideaName, /*json:*/jsonText)
    if err != nil {
        return -1, err
    }
    defer stmt.Close()

    return res.LastInsertId()
}

func readIdeas(db *sql.DB, searchDirPath string) (error) {

    fileList := []string{}
    err := filepath.Walk(searchDirPath, func(path string, f os.FileInfo, err error) error {
        if f.IsDir() != true && filepath.Ext(path) == ".json" {
            fileList = append(fileList, path)
        }
        return err
    })

    if err != nil {
        return err
    }

    var wg sync.WaitGroup
    var mutex = &sync.Mutex{}

    for _, fileName := range fileList {
        wg.Add(1)
        go func(fileName string) {
            defer wg.Done()
            ideaName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
            log.Println("Inserting", ideaName)
            jsonText, err := ioutil.ReadFile(fileName)
            if err != nil {
                return
            }
            mutex.Lock()
            _, err = insertIdea(db, ideaName, string(jsonText))
            mutex.Unlock()
            if err != nil {
                log.Println("can't insert into idea table:", err)
            }
            log.Println("Inserted", ideaName)
            //wg.Done()
        }(fileName)
    }
    log.Println("Waiting for all idea files to be processed...")
    wg.Wait()
    log.Println("Done")

    return nil
}

func showIdeaTable(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, json FROM idea")
	if err != nil {
		panic(err)//log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
                var json string
		err = rows.Scan(&id, &name, &json)
		if err != nil {
			panic(err)//log.Fatal(err)
		}
		fmt.Println(id, name, json)
	}
}

func main() {
    dbFilePath := "./"+dbFileName
    os.Remove(dbFilePath)
    if err := createDB(dbName, dbFilePath); err != nil {
        panic(err)
    }

    db, err := openDB(dbFilePath)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    readIdeas(db, "./src/github.com/curoles/answer42/dbsrc/idea")
    //showIdeaTable(db)
}
