package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "Nam12345"
	DB_NAME     = "mydb"
)

type User struct {
	Id       int    `json:id`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	fmt.Println("Getting users...")

	rows, err := db.Query("SELECT * FROM users")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var users []User

	for rows.Next() {
		var id int
		var name string
		var birthday string
		var gender string
		var email string

		err = rows.Scan(&id, &name, &birthday, &gender, &email)

		// check errors
		checkErr(err)

		users = append(users, User{Id: id, Name: name, Birthday: birthday, Gender: gender, Email: email})
	}

	var response = JsonResponse{Type: "success", Data: users}

	json.NewEncoder(w).Encode(response)

	fmt.Println("Get users successfully!!")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["id"]

	var response = JsonResponse{}

	if userId == "" {
		response = JsonResponse{Type: "error", Message: "You are missing userId parameter."}
	} else {
		db := setupDB()

		fmt.Println("getting user " + userId + " from DB...")

		rows, err := db.Query("SELECT * FROM users where id = " + userId)

		// check errors
		checkErr(err)

		var users []User

		for rows.Next() {
			var id int
			var name string
			var birthday string
			var gender string
			var email string

			err = rows.Scan(&id, &name, &birthday, &gender, &email)

			// check errors
			checkErr(err)

			users = append(users, User{Id: id, Name: name, Birthday: birthday, Gender: gender, Email: email})
		}

		response = JsonResponse{Type: "success", Data: users, Message: "Get user " + userId + " successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	/*
		w.Header().Set("Content-Type", "application/json")
		name := r.FormValue("name")
		birthday := r.FormValue("birthday")
		gender := r.FormValue("gender")
		email := r.FormValue("email")

		var response = JsonResponse{}


			if name == "" {
				response = JsonResponse{Type: "error", Message: "You are missing name parameters."}
			} else if birthday == "" {
				response = JsonResponse{Type: "error", Message: "You are missing birthday parameters."}
			} else if gender == "" {
				response = JsonResponse{Type: "error", Message: "You are missing gender parameters."}
			} else if email == "" {
				response = JsonResponse{Type: "error", Message: "You are missing email parameters."}
			} else {

		db := setupDB()

		fmt.Println("Inserting user into DB")

		fmt.Println("Inserting new user with name: " + name + ", birthday: " + birthday + ", gender: " + gender + ", email: " + email)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO users(name, birthday, gender, email) VALUES($1, $2, $3, $4) returning id;", name, birthday, gender, email).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The user has been inserted successfully!"}
		}

		json.NewEncoder(w).Encode(response)
	*/

	w.Header().Set("Content-Type", "application/json")

	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)

	var response = JsonResponse{}

	db := setupDB()

	fmt.Println("Inserting user into DB")

	fmt.Println("Inserting new user with name: " + user.Name + ", birthday: " + user.Birthday + ", gender: " + user.Gender + ", email: " + user.Email)

	var lastInsertID int
	err := db.QueryRow("INSERT INTO users(name, birthday, gender, email) VALUES($1, $2, $3, $4) returning id;", user.Name, user.Birthday, user.Gender, user.Email).Scan(&lastInsertID)

	// check errors
	checkErr(err)

	response = JsonResponse{Type: "success", Message: "The user has been inserted successfully!"}

	json.NewEncoder(w).Encode(response)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userId := params["id"]

	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)

	var response = JsonResponse{}

	db := setupDB()

	fmt.Println("updating user " + userId + " from DB...")

	// create the update sql query
	sqlStatement := `UPDATE users SET name=$2, birthday=$3, gender=$4, email=$5 WHERE id=$1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, userId, user.Name, user.Birthday, user.Gender, user.Email)

	// check errors
	checkErr(err)

	response = JsonResponse{Type: "success", Message: "Update user " + userId + " successfully!"}

	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["id"]

	var response = JsonResponse{}

	if userId == "" {
		response = JsonResponse{Type: "error", Message: "You are missing userId parameter."}
	} else {
		db := setupDB()

		fmt.Println("Deleting user from DB")

		_, err := db.Exec("DELETE FROM users where id = $1", userId)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	fmt.Println("Server at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
