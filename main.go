package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"message_board/helpers/error_handlers"
	"message_board/utils/checkers"
	"message_board/utils/cookies"
	"message_board/utils/queries"
	"net/http"
	"os"
	"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("./logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	rand.Seed(time.Now().Unix())
}

func main() {
	defer error_handlers.CloseDatabase(db)

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/index", index)
	http.HandleFunc("/watch", watch)
	http.HandleFunc("/write", write)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/user", user)
	http.HandleFunc("/examine", examine)
	http.HandleFunc("/manage", manage)
	addWebSocketHandlers()

	_ = http.ListenAndServe(":8080", nil)
}

func index(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)
	tmpl, err := template.ParseFiles("./view/index.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}
	err = tmpl.Execute(rw, "")
	if err != nil {
		log.Println("error when execute template:", err)
		return
	}
}

func watch(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/watch.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}
	err = tmpl.Execute(rw, "")
	if err != nil {
		log.Println("error when execute template:", err)
		return
	}
}

func write(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/write.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println("error when parse form", err)
		return
	}
	content := r.PostForm.Get("content")
	if len(content) == 0 {
		err = tmpl.Execute(rw, "")
		if err != nil {
			log.Println("error when execute template:", err)
			rw.WriteHeader(500)
		}
		return
	}

	submitStatus := "提交成功"
	cookieCollection := r.Cookies()
	for _, cookie := range cookieCollection {
		if cookie.Name == "login" {
			loginCookie := cookie.Value
			uid, err := queries.FindLoginCookie(loginCookie, db)
			if err != nil {
				log.Println("error when find login cookie:", err)
				submitStatus = "提交失败"
				break
			}
			if uid == -1 {
				submitStatus = "提交失败"
				break
			}
			exam := 0
			power, err := queries.FindPowerByUID(uid, db)
			if err != nil {
				log.Println("error when find poser by uid:", err)
				submitStatus = "提交失败"
				break
			}
			if power >= 1 {
				exam = 1
			} else {
				submitStatus += "，等待管理员审核"
			}
			err = queries.AddNewMessage(uid, content, exam, db)
			if err != nil {
				log.Println("error when add new message:", err)
				submitStatus = "提交失败"
				break
			}
			break
		}
	}
	err = tmpl.Execute(rw, submitStatus)
	if err != nil {
		log.Println("error when execute template:", err)
		rw.WriteHeader(500)
	}
}

func login(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/login.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("error when parse form", err)
		return
	}

	name := r.PostForm.Get("name")
	pass := r.PostForm.Get("pass")
	if len(name) == 0 && len(pass) == 0 {
		err = tmpl.Execute(rw, "")
		if err != nil {
			log.Println("error when execute template:", err)
			rw.WriteHeader(500)
		}
		return
	}

	uid, err := queries.FindLoginUID(name, pass, db)
	if err != nil {
		log.Println("error when find login uid:", err)
		rw.WriteHeader(500)
		return
	}
	if uid == -1 {
		err = tmpl.Execute(rw, "用户名或密码错误")
		if err != nil {
			log.Println("error when execute template:", err)
			rw.WriteHeader(500)
		}
		return
	}

	loginCookie := cookies.GetLoginCookie(uid)
	err = queries.UpdateLoginCookie(uid, loginCookie, db)
	if err != nil {
		log.Println("error when update login cookie:", err)
		rw.WriteHeader(500)
		return
	}
	rw.Header().Add("set-cookie", "login="+loginCookie)
	err = tmpl.Execute(rw, "")
	if err != nil {
		log.Println("error when execute template:", err)
		rw.WriteHeader(500)
	}
	return
}

func register(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/register.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		rw.WriteHeader(500)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("error when parse form:", err)
		rw.WriteHeader(500)
		return
	}

	name := r.PostForm.Get("name")
	pass := r.PostForm.Get("pass")
	if len(name) == 0 && len(pass) == 0 {
		err = tmpl.Execute(rw, "")
		if err != nil {
			log.Println("error when execute template:", err)
			rw.WriteHeader(500)
		}
		return
	}
	checkRes := checkers.CheckUserName(name, db)
	if checkRes != "ok" {
		log.Println("User name to be registered is invalid.")
		rw.WriteHeader(500)
		return
	}
	err = queries.CreateNewUser(name, pass, db)
	if err != nil {
		log.Println("error when create new user:", err)
		rw.WriteHeader(500)
		return
	}

	uid, err := queries.FindUIDByName(name, db)
	if err != nil {
		log.Println("error when find uid by name:", err)
	} else {
		loginCookie := cookies.GetLoginCookie(uid)
		err = queries.UpdateLoginCookie(uid, loginCookie, db)
		if err != nil {
			log.Println("error when update login cookie:", err)
		}
		rw.Header().Add("set-cookie", "login="+loginCookie)
	}

	tmpl, err = template.ParseFiles(
		"./view/register_success.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}
	err = tmpl.Execute(rw, "")
	if err != nil {
		log.Println("error when execute template:", err)
		return
	}
}

func user(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/user.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}

	err = tmpl.Execute(rw, "")
	if err != nil {
		log.Println("error when execute template:", err)
		rw.WriteHeader(500)
	}
	return
}

func examine(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/examine.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}
	err = tmpl.Execute(rw, "")
	if err != nil {
		rw.WriteHeader(500)
		return
	}
}

func manage(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)

	tmpl, err := template.ParseFiles(
		"./view/manage.html", "./view/header.html", "./view/footer.html")
	if err != nil {
		log.Println("error when parse template:", err)
		return
	}
	err = tmpl.Execute(rw, "")
	if err != nil {
		rw.WriteHeader(500)
		return
	}
}
