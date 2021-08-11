package queries

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"message_board/helpers/error_handlers"
	"message_board/utils/others"
	"time"
)

func FindNameByUID(uid int, db *sql.DB) (string, error) {
	prep, err := db.Prepare("SELECT NAME FROM USERS WHERE UID=?")
	if err != nil {
		return "", err
	}
	rows, err := prep.Query(uid)
	defer error_handlers.CloseRows(rows)
	if err != nil {
		return "", err
	}
	if rows.Next() == false {
		return "", nil
	}
	var name string
	err = rows.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func FindPowerByUID(uid int, db *sql.DB) (int, error) {
	prep, err := db.Prepare("SELECT POWER FROM USERS WHERE UID=?")
	if err != nil {
		return -1, err
	}
	rows, err := prep.Query(uid)
	defer error_handlers.CloseRows(rows)
	if err != nil {
		return -1, err
	}
	if rows.Next() == false {
		return -1, nil
	}
	var power int
	err = rows.Scan(&power)
	if err != nil {
		return -1, err
	}
	return power, nil
}

func FindUIDByName(name string, db *sql.DB) (int, error) {
	prep, err := db.Prepare("SELECT UID FROM USERS WHERE NAME=?")
	if err != nil {
		return -1, err
	}
	rows, err := prep.Query(name)
	defer error_handlers.CloseRows(rows)
	if err != nil {
		return -1, err
	}
	if rows.Next() == false {
		return -1, nil
	}
	var uid int
	err = rows.Scan(&uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

func FindLoginUID(name string, pass string, db *sql.DB) (int, error) {
	passMd5 := md5.Sum([]byte(pass))
	passEnc := hex.EncodeToString(passMd5[:])
	prep, err := db.Prepare("SELECT UID FROM USERS WHERE NAME=? AND PASS=?")
	if err != nil {
		return -1, err
	}
	rows, err := prep.Query(name, passEnc)
	defer error_handlers.CloseRows(rows)
	if err != nil {
		return -1, err
	}
	if rows.Next() == false {
		return -1, nil
	}
	var uid int
	err = rows.Scan(&uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

func FindLoginCookie(loginCookie string, db *sql.DB) (int, error) {
	prep, err := db.Prepare("SELECT UID FROM COOKIES WHERE LOGIN=?")
	if err != nil {
		return -1, err
	}
	rows, err := prep.Query(loginCookie)
	defer error_handlers.CloseRows(rows)
	if err != nil {
		return -1, err
	}
	if rows.Next() == false {
		return -1, nil
	}
	var uid int
	err = rows.Scan(&uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

type User struct {
	Uid int
	Name string
	Power string
}

func FindLowerPowerUsers(power int ,name string , db *sql.DB) ([]User, error) {
	var empty []User
	if name == "" {
		name = "%"
	}
	prep, err := db.Prepare(
		"SELECT UID, NAME, POWER FROM USERS WHERE POWER<=? AND NAME LIKE ?")
	if err != nil {
		return empty, err
	}
	rows, err := prep.Query(power, name)
	if err != nil {
		return empty, err
	}
	defer error_handlers.CloseRows(rows)
	var users []User
	for rows.Next() {
		var userData User
		var userPower int
		err = rows.Scan(&userData.Uid, &userData.Name, &userPower)
		if err != nil {
			return empty, err
		}
		userData.Power = others.PowerName[userPower]
		users = append(users, userData)
	}
	return users, nil
}

func UpdateLoginCookie(uid int, loginCookie string, db *sql.DB) error {
	prep, err := db.Prepare("UPDATE COOKIES SET LOGIN=? WHERE UID=?")
	if err != nil {
		return err
	}
	_, err = prep.Exec(loginCookie, uid)
	return err
}

func CreateNewUser(name string, pass string, db *sql.DB) error {
	prep, err := db.Prepare("INSERT INTO USERS(NAME, PASS) VALUES(?, ?)")
	if err != nil {
		return err
	}
	passMd5 := md5.Sum([]byte(pass))
	passEnc := hex.EncodeToString(passMd5[:])
	_, err = prep.Exec(name, passEnc)
	if err != nil {
		return err
	}
	uid, err := FindUIDByName(name, db)
	if err != nil {
		log.Println("error when find uid by name:", err)
		return err
	}
	prep, err = db.Prepare("INSERT INTO COOKIES(UID) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = prep.Exec(uid)
	return err
}

func AddNewMessage(uid int, content string, exam int, db *sql.DB) error {
	prep, err := db.Prepare("INSERT INTO MESSAGE(UID, CONTENT, UNIX_TIME, EXAM) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = prep.Exec(uid, content, time.Now().Unix(), exam)
	return err
}

func ChangeUserPower(uid int, power int, db *sql.DB) error {
	prep, err := db.Prepare("UPDATE USERS SET POWER=? WHERE UID=?")
	if err != nil {
		return err
	}
	_, err = prep.Exec(power, uid)
	return err
}
