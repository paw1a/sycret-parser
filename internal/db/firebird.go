package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/nakagami/firebirdsql"
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

func ScanSelectRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	objectArray := make([]map[string]interface{}, 0)
	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}

		object := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			object[colName] = *val
		}

		objectArray = append(objectArray, object)
	}

	return objectArray, nil
}
