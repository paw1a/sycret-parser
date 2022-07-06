package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBConnection struct {
	DBServer string `json:"DBSERVER"`
	DBPath   string `json:"DBPATH"`
}

func NewDB(dbConnection DBConnection) (*sqlx.DB, error) {
	db, err := sqlx.Connect("firebirdsql", fmt.Sprintf("SYSDBA:masterkey@%s/%s", dbConnection.DBServer, dbConnection.DBPath))
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db is unavailable: %v", err)
	}

	return db, err
}
