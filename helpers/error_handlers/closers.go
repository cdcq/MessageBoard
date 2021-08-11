package error_handlers

import (
	"database/sql"
	"log"
	"net/http"
)

func CloseHttpRequest(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		log.Println("error when close http request:", err)
	}
}

func CloseDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("error when close database:", err)
	}
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Println("error when close rows:", err)
	}
}