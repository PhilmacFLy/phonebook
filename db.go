package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

type DBConnection struct {
	Driver     string
	Connection string
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func initialisation(dbc *DBConnection) {
	var err error
	db, err = sqlx.Open(dbc.Driver, dbc.Connection)
	if err != nil {
		log.Fatal(err)
	}
	initDB(dbc)
}

func initDB(dbc *DBConnection) {
	switch dbc.Driver {
	case "sqlite3":
		cont, err := exists(dbc.Connection)
		if err != nil {
			log.Fatal(err)
		}
		if cont {
			fmt.Println("cont")
			return
		}
		_, err = os.Create(dbc.Connection)
		if err != nil {
			log.Fatal("Could not create file "+dbc.Connection, err)
		}
		_, err = db.Exec(createSQLlitestmt)
		if err != nil {
			log.Printf("%q: %s\n", err, createSQLlitestmt)
			return
		}

	default:
		log.Fatal("DB Driver unkown. Stopping Server")
	}
}

type phonenumber struct {
	NumberID int    `xml:"-" db:"NumberID"`
	Name     string `xml:"Name" db:"Name"`
	Number   string `xml:"Telephone" db:"Number"`
}

func GetNumbersforExtension(extension string) ([]phonenumber, error) {
	var pns []phonenumber
	var ids []int
	err := db.Select(&ids, "Select NumberID from phonebooks Where Extension = ?", extension)
	if err != nil {
		return pns, errors.New("Error fetching NumberIDs: " + err.Error())
	}
	for _, id := range ids {
		var pn phonenumber
		err := db.Get(&pn, "Select * from phonenumbers Where NumberID = ?", id)
		if err != nil {
			return pns, errors.New("Error fetching phonenumberid: " + err.Error())
		}
		pns = append(pns, pn)
	}
	return pns, nil
}

const createSQLlitestmt = `
--
-- File generated with SQLiteStudio v3.0.7 on Mi. Jan. 3 14:43:55 2018
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: extensions
CREATE TABLE extensions (Extension INT NOT NULL PRIMARY KEY);

-- Table: phonebooks
CREATE TABLE phonebooks (Extension INTEGER REFERENCES extensions (extension) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL, NumberID INTEGER REFERENCES phonenumbers (NumberID) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL, PRIMARY KEY (Extension, NumberID));

-- Table: phonenumbers
CREATE TABLE phonenumbers (NumberID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, Name TEXT NOT NULL, Number TEXT NOT NULL);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;

`
