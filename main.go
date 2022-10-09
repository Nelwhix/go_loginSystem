package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
  "github.com/asaskevich/govalidator"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

const (
  USERNAME = "Nelwhix"
  PASSWORD = "admin"
  USERNAME_ERROR_MESSAGE = "Please enter a valid Username"
  PASSWORD_ERROR_MESSAGE = "Please enter a valid Password"
  GENERIC_ERROR_MESSAGE = "Validation Error"
)

type User struct {
  Username string `valid:"alpha,required"`
  Password string `valid:"alpha,required"`
}

func readLoginForm(r *http.Request) *User {
  r.ParseForm()
  user := new(User)
  decoder := schema.NewDecoder()
  decodeErr := decoder.Decode(user, r.PostForm)

  if decodeErr != nil {
    log.Printf("error mapping parsed form data to struct : ", decodeErr)
  }
  return user
}

func validateUser(w http.ResponseWriter, r *http.Request, user *User) (bool, string) {
  valid, validationError := govalidator.ValidateStruct(user)

  if !valid {
    usernameError := govalidator.ErrorByField(validationError, "Username")
    passwordError := govalidator.ErrorByField(validationError, "Password")

    if usernameError != "" {
      log.Printf("Username validation error : ", usernameError)
      return valid, USERNAME_ERROR_MESSAGE
    }

    if passwordError != "" {
      log.Printf("password validation error : ", passwordError)
      return valid, PASSWORD_ERROR_MESSAGE
    }
  }
  return valid, GENERIC_ERROR_MESSAGE
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
  parsedTemplate, _ := template.ParseFiles("./templates/index.html")
  parsedTemplate.Execute(w, nil)
}

func renderRegisterForm(w http.ResponseWriter, r *http.Request) {
  parsedTemplate, _ := template.ParseFiles("./templates/register.html")
  parsedTemplate.Execute(w, nil)
}

func logInUser(w http.ResponseWriter, r *http.Request) {
  user := readLoginForm(r)
  valid, validationErrorMessage := validateUser(w, r, user)

  if !valid {
    fmt.Fprint(w, validationErrorMessage)
    return
  }

  if (user.Username == USERNAME && user.Password == PASSWORD) {
    fmt.Fprintf(w, "Hello " + user.Username + "!")
  } else {
    fmt.Fprintf(w, "Bad credentials")
  }
}

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/", renderTemplate).Methods("GET")
  router.HandleFunc("/login", logInUser).Methods("POST")
  router.HandleFunc("/signup", renderRegisterForm).Methods("GET")
  // router.HandleFunc("/signup", signUpUser).Methods("POST")
  router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

  log.Printf("Server starting on port %v\n", CONN_PORT)
  err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
  
  if err != nil {
    log.Fatal("error starting http server : ", err)
    return
  }
}

