package queries

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"message_board/helpers/error_handlers"
	"message_board/utils/others"
	"time"
)

type Message struct {
	Mid     int
	Name    string
	Content string
	Time 	string
	Power 	string
	Exam    int
}

func FindMessageByeExam(exam int, db *sql.DB) ([]Message, error) {
	var empty []Message
	prep, err := db.Prepare("SELECT * FROM MESSAGE WHERE EXAM=? ORDER BY MID DESC")
	if err != nil {
		return empty, err
	}
	rows, err := prep.Query(exam)
	if err != nil {
		return empty, err
	}
	defer error_handlers.CloseRows(rows)
	var ret []Message
	for rows.Next() {
		var row Message
		var uid int
		var unixTime int64
		err = rows.Scan(&row.Mid, &uid, &row.Content, &unixTime, &row.Exam)
		if err != nil {
			return empty, err
		}
		name, err := FindNameByUID(uid, db)
		if err != nil {
			return empty, err
		}
		power, err := FindPowerByUID(uid, db)
		if err != nil {
			return empty, err
		}
		powerS := others.PowerName[power]
		row.Name = name
		row.Time = time.Unix(unixTime, 0).Format("2006-01-02 15:04:05")
		row.Power = powerS
		ret = append(ret, row)
	}
	return ret, nil
}

func AccessMessage(mid int, db *sql.DB) error {
	prep, err := db.Prepare("UPDATE MESSAGE SET EXAM=? WHERE MID=?")
	if err != nil {
		return err
	}
	_, err = prep.Exec(1, mid)
	if err != nil {
		return err
	}
	return nil
}
