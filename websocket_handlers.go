package main

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/websocket"
	"log"
	"message_board/utils/checkers"
	cookies2 "message_board/utils/cookies"
	"message_board/utils/others"
	"message_board/utils/queries"
	"net/http"
	"strconv"
	"syscall"
)

func addWebSocketHandlers() {
	http.Handle("/register/check-user-name", websocket.Handler(checkUserName))
	http.Handle("/login/check-login-status", websocket.Handler(checkLoginStatus))
	http.Handle("/user/logout", websocket.Handler(userLogout))
	http.Handle("/watch/get_messages", websocket.Handler(getMessages))
	http.Handle("/examine/get_examine_messages", websocket.Handler(getExamineMessages))
	http.Handle("/examine/access_message", websocket.Handler(accessMessage))
	http.Handle("/manage/get_lower_power_users", websocket.Handler(getLowerPowerUser))
	http.Handle("/manage/change_user_power", websocket.Handler(changeUserPower))
}

type loginStatus struct {
	Name  string
	Power int
}

func checkLoginStatus(ws *websocket.Conn) {
	cookies := ws.Request().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "login" {
			loginCookie := cookie.Value
			uid, err := queries.FindLoginCookie(loginCookie, db)
			if err != nil {
				log.Println("error when find login cookie:", err)
				break
			}
			if uid == -1 {
				break
			}
			name, err := queries.FindNameByUID(uid, db)
			if err != nil {
				log.Println("error when find name by uid:", err)
				break
			}
			power, err := queries.FindPowerByUID(uid, db)
			if err != nil {
				log.Println("error when find power by uid:", err)
				break
			}
			status := loginStatus{Name: name, Power: power}
			message, err := json.Marshal(status)
			if err != nil {
				log.Println("error when marshal json:", err)
				break
			}
			err = websocket.Message.Send(ws, string(message))
			if err != nil {
				log.Println("error when send websocket message:", err)
				break
			}
			return
		}
	}
	err := websocket.Message.Send(ws, "no")
	if err != nil {
		log.Println("error when send message to websocket:", err)
	}
}

func checkUserName(ws *websocket.Conn) {
	var err error
	var data string
	for {
		err = websocket.Message.Receive(ws, &data)
		if err == syscall.EAGAIN {
			err = nil
			continue
		}
		if err != nil {
			log.Println("error when websocket receive:", err)
			break
		}
		err = websocket.Message.Send(ws, checkers.CheckUserName(data, db))
		if err != nil {
			log.Println("error when websocket send:", err)
		}
		return
	}
}

func userLogout(ws *websocket.Conn) {
	cookies := ws.Request().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "login" {
			loginCookie := cookie.Value
			uid, err := queries.FindLoginCookie(loginCookie, db)
			if err != nil {
				log.Println("error when find login cookie:", err)
				break
			}
			if uid == -1 {
				break
			}
			err = queries.UpdateLoginCookie(uid, "", db)
			if err != nil {
				log.Println("error when update login cookie:", err)
				break
			}
			return
		}
	}
}

type messageData struct {
	Rows    []queries.Message
	PageCnt int
}

func getMessages(ws *websocket.Conn) {
	var pageS string
	for {
		err := websocket.Message.Receive(ws, &pageS)
		if err == syscall.EAGAIN {
			err = nil
			continue
		}
		if err != nil {
			log.Println("error when receive websocket message:", err)
			return
		}
		break
	}
	page, err := strconv.Atoi(pageS)
	if err != nil {
		log.Println("error when string to int:", err)
		return
	}
	rows, err := queries.FindMessageByeExam(1, db)
	if err != nil {
		log.Println("error when find visible message:", err)
		return
	}
	rowsPerPage := 5
	pageCnt := (len(rows) + rowsPerPage - 1) / rowsPerPage
	if page <= 0 || page > pageCnt {
		return
	}
	if len(rows) < (page-1)*rowsPerPage {
		return
	}
	rows = rows[page*rowsPerPage-rowsPerPage : others.Min(len(rows), page*rowsPerPage)]
	data := messageData{Rows: rows, PageCnt: pageCnt}
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Println("error when marshal rows:", err)
		return
	}
	dataS := string(dataJson[:])
	err = websocket.Message.Send(ws, dataS)
	if err != nil {
		log.Println("error when send websocket message:", err)
		return
	}
}

func getExamineMessages(ws *websocket.Conn) {
	cookies := ws.Request().Cookies()
	uid, err := cookies2.GetLoginUid(cookies, db)
	if err != nil {
		log.Println("error when get login uid:", err)
		return
	}
	if uid < 1 {
		return
	}
	power, err := queries.FindPowerByUID(uid, db)
	if err != nil {
		log.Println("error when find power by uid:", err)
		return
	}
	if power < 1 {
		return
	}
	var pageS string
	for {
		err = websocket.Message.Receive(ws, &pageS)
		if err == syscall.EAGAIN {
			err = nil
			continue
		}
		if err != nil {
			log.Println("error when receive websocket message:", err)
			return
		}
		break
	}
	page, err := strconv.Atoi(pageS)
	if err != nil {
		log.Println("error when string to int:", err)
		return
	}
	rows, err := queries.FindMessageByeExam(0, db)
	if err != nil {
		log.Println("error when find visible message:", err)
		return
	}
	rowsPerPage := 5
	pageCnt := (len(rows) + rowsPerPage - 1) / rowsPerPage
	if page <= 0 || page > pageCnt {
		return
	}
	if len(rows) < (page-1)*rowsPerPage {
		return
	}
	rows = rows[page*rowsPerPage-rowsPerPage : others.Min(len(rows), page*rowsPerPage)]
	data := messageData{Rows: rows, PageCnt: pageCnt}
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Println("error when marshal rows:", err)
		return
	}
	dataS := string(dataJson[:])
	err = websocket.Message.Send(ws, dataS)
	if err != nil {
		log.Println("error when send websocket message:", err)
		return
	}
}

func accessMessage(ws *websocket.Conn) {
	cookies := ws.Request().Cookies()
	uid, err := cookies2.GetLoginUid(cookies, db)
	if err != nil {
		log.Println("error when get login uid:", err)
		return
	}
	if uid < 1 {
		return
	}
	power, err := queries.FindPowerByUID(uid, db)
	if err != nil {
		log.Println("error when find power by uid:", err)
		return
	}
	if power < 1 {
		return
	}
	var data string
	for {
		err = websocket.Message.Receive(ws, &data)
		if err == syscall.EAGAIN {
			err = nil
			continue
		}
		if err != nil {
			log.Println("error when websocket receive:", err)
			break
		}
		mid, err := strconv.Atoi(data)
		if err != nil {
			log.Println("error when turn string to int:", err)
			break
		}
		err = queries.AccessMessage(mid, db)
		if err != nil {
			log.Println("error when access message in database:", err)
			break
		}
		return
	}
}

type searchNameAndPage struct {
	Name string
	Page int
}

type usersData struct {
	Rows    []queries.User
	PageCnt int
}

func getLowerPowerUser(ws *websocket.Conn) {
	cookies := ws.Request().Cookies()
	uid, err := cookies2.GetLoginUid(cookies, db)
	if err != nil {
		log.Println("error when get login uid:", err)
		return
	}
	if uid < 1 {
		return
	}
	power, err := queries.FindPowerByUID(uid, db)
	if err != nil {
		log.Println("error when find power by uid:", err)
		return
	}
	var dataS string
	for {
		err = websocket.Message.Receive(ws, &dataS)
		if err == syscall.EAGAIN {
			err = nil
			continue
		}
		if err != nil {
			log.Println("error when websocket receive:", err)
			break
		}
		var data searchNameAndPage
		err = json.Unmarshal([]byte(dataS), &data)
		if err != nil {
			log.Println("error when unmarshal json:", err)
			break
		}
		name := data.Name
		page := data.Page
		rows, err := queries.FindLowerPowerUsers(power, name, db)
		if err != nil {
			log.Println("error when find lower power users:", err)
			break
		}
		rowsPerPage := 10
		pageCnt := (len(rows) + rowsPerPage - 1) / rowsPerPage
		if page <= 0 || page > pageCnt {
			break
		}
		if len(rows) < (page-1)*rowsPerPage {
			break
		}
		rows = rows[page*rowsPerPage-rowsPerPage : others.Min(len(rows), page*rowsPerPage)]

		resData := usersData{Rows: rows, PageCnt: pageCnt}
		resDataJson, err := json.Marshal(resData)
		if err != nil {
			log.Println("error when marshal json:", err)
			break
		}
		err = websocket.Message.Send(ws, string(resDataJson))
		if err != nil {
			log.Println("error when send websocket message:", err)
			break
		}
		return
	}
}

type changePowerData struct {
	Uid   int
	Power int
}

func changeUserPower(ws *websocket.Conn) {
	cookies := ws.Request().Cookies()
	uid, err := cookies2.GetLoginUid(cookies, db)
	if err != nil {
		log.Println("error when get login uid:", err)
		return
	}
	if uid < 1 {
		return
	}
	power, err := queries.FindPowerByUID(uid, db)
	if err != nil {
		log.Println("error when find power by uid:", err)
		return
	}
	if power < 1 {
		return
	}
	var data changePowerData
	var dataS string
	for {
		err = websocket.Message.Receive(ws, &dataS)
		if err == syscall.EAGAIN {
			err = nil
			continue
		}
		if err != nil {
			log.Println("error when websocket receive:", err)
			break
		}
		err = json.Unmarshal([]byte(dataS), &data)
		if err != nil {
			log.Println("error when unmarshal json data:", err)
			break
		}
		if data.Power > power {
			break
		}
		err = queries.ChangeUserPower(data.Uid, data.Power, db)
		if err != nil {
			log.Println("error when change user poser:", err)
			break
		}
		return
	}
}
