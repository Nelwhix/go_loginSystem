package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

const (
  Username = "Nelwhix"
  Password = "admin"
)

type User struct {
  Username string
  Password string
}

func readForm(r *http.Request) *User {
  r.ParseForm()
  user := new(User)
  decoder := schema.NewDecoder()
  decodeErr := decoder.Decode(user, r.PostForm)

  if decodeErr != nil {
    log.Printf("error mapping parsed form data to struct : ", decodeErr)
  }
  return user
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
  parsedTemplate, _ := template.ParseFiles("./templates/index.html")
  parsedTemplate.Execute(w, nil)
}

func logInUser(w http.ResponseWriter, r *http.Request) {
  user := readForm(r)
  
  if (user.Username == Username && user.Password == Password) {
    fmt.Fprintf(w, "Hello " + user.Username + "!")
  } else {
    fmt.Fprintf(w, "Bad credentials")
  }
}
func main() {
  router := mux.NewRouter()
  router.HandleFunc("/", renderTemplate).Methods("GET")
  router.HandleFunc("/login", logInUser).Methods("POST")
  router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

  log.Printf("Server starting on port %v\n", CONN_PORT)
  err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
  
  if err != nil {
    log.Fatal("error starting http server : ", err)
    return
  }
}

