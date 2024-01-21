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

type Employee struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

func main() {
	//connect to database
	DB_URL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", "5432", "postgres", "inventory", "postgres")
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connection successful")
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS employee (id SERIAL PRIMARY KEY, name TEXT, profile TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	//create router
	router := mux.NewRouter()
	router.HandleFunc("/employee", getEmployees(db)).Methods("GET")
	router.HandleFunc("/employee/{id}", getEmployee(db)).Methods("GET")
	router.HandleFunc("/employee", createEmployee(db)).Methods("POST")
	router.HandleFunc("/employee/{id}", updateEmployee(db)).Methods("PUT")
	router.HandleFunc("/employee/{id}", deleteEmployee(db)).Methods("DELETE")

	fmt.Println("Listening to port 8000")

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	//allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token, Access-Control-Allow-Origin"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			//w.Header().Set("Access-Control-Allow-Origin", allowedHeaders)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Origin", allowedHeaders)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all employee
func getEmployees(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM employee")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		employee := []Employee{}
		for rows.Next() {
			var u Employee
			if err := rows.Scan(&u.ID, &u.Name, &u.Profile); err != nil {
				log.Fatal(err)
			}
			employee = append(employee, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(employee)
	}
}

// get user by id
func getEmployee(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u Employee
		err := db.QueryRow("SELECT * FROM employee WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Profile)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// create user
func createEmployee(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Employee
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO employee (name, profile) VALUES ($1, $2) RETURNING id", u.Name, u.Profile).Scan(&u.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

// update user
func updateEmployee(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Employee
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE employee SET name = $1, profile = $2 WHERE id = $3", u.Name, u.Profile, id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

// delete user
func deleteEmployee(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u Employee
		err := db.QueryRow("SELECT * FROM employee WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Profile)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM employee WHERE id = $1", id)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Employee deleted")
		}
	}
}
