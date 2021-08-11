package checkers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"message_board/helpers/error_handlers"
)

func CheckUserName(name string, db *sql.DB) string {
	if len(name) == 0 {
		return "用户名不能为空"
	}
	if len(name) > 30 {
		return "用户名不能超过 30 个字符"
	}

	prep, err := db.Prepare("SELECT UID FROM USERS WHERE NAME=?")
	if err != nil {
		log.Println("error when prepare query:", err)
		return "error"
	}
	rows, err := prep.Query(name)
	defer error_handlers.CloseRows(rows)
	if err != nil {
		log.Println("error when query database:", err)
		return "error"
	}
	if rows.Next() != false {
		return "用户名已存在"
	}

	return "ok"
}