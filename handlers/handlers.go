package handlers

import (
	"github.com/Orphei104/database"
	

	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

var h = New(database.Connect())
var r = mux.NewRouter()

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world")
}

func HandleRequests(){
	r.HandleFunc("/", mainPage)
	fmt.Println("server is started")
	http.ListenAndServe(":2004", r)
}