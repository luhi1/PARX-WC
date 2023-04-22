package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

type DisplayError struct {
	ErrorDescription string
}

// TeacherPageHandlers Consider creating a generic interface for both teacher and student to implement.
type TeacherPageHandlers interface {
	GETHandler(writer http.ResponseWriter, request *http.Request)
	POSTHandler(writer http.ResponseWriter, request *http.Request)
	valHandler(writer http.ResponseWriter, request *http.Request)
	dataVal(requestMethod string) bool
}

var db *sql.DB

// Start server run, files, and other shit.
func main() {
	userInfo := UserData{}
	eventInfo := EventInfo{}
	prize := Prize{}
	winners := Winners{}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login", userInfo.GETHandler)

	http.HandleFunc("/signup", userInfo.POSTHandler)

	http.HandleFunc("/userValidation/", userInfo.valHandler)

	http.HandleFunc("/teacherEvents", eventInfo.GETHandler)

	http.HandleFunc("/teacherCreateEvent", eventInfo.POSTHandler)

	http.HandleFunc("/winners", winners.GETHandler)
	http.HandleFunc("/prizes", prize.GETHandler)

	http.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		userInfo = UserData{}
		http.Redirect(writer, request, "./login", 307)
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			err := tplExec(writer, "error.gohtml", nil)
			//@TODO: REMOVE
			if err != nil {
				return
			}
		} else {
			http.Redirect(writer, request, "./login", 301)
		}
	})

	initdb, err := sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/fbla")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db = initdb

	fmt.Println("Connected to DB")

	/*@todo: Add this to the setup wizard eventually */
	fmt.Println("Server is running on port 8082")

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("Error starting server, aborting tasks")
		panic(err)
	}
}

func tplExec(w http.ResponseWriter, filename string, information any) error {
	temp := template.Must(template.ParseFiles("./WebPages/" + filename))

	err := temp.Execute(w, information)
	//@TODO: REMOVE
	if err != nil {
		return err
	}
	return nil
}

func hashPswd(pwd string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pwd))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
