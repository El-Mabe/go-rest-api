package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type UserScore struct {
	ListScore []User
}

func GetConnection() *sql.DB {
	connStr := "postgresql://root@localhost:26257?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addUser")
	fmt.Println("Entro a addUser")

	db := GetConnection()
	var data User
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	id := strconv.Itoa(data.ID)
	score := strconv.Itoa(data.Score)
	ress := id + ",'" + data.Name + "'," + score

	fmt.Println("result: ", ress)
	if _, err := db.Query(
		fmt.Sprintf("INSERT INTO snake.users (id, name, score) VALUES (%s)",
			ress)); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	json, _ := json.Marshal(data)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-ControlAllow-Headers, Authorization, X-Requested-Widh")
	w.Header().Set("Access-Control-Allow-Credentials", "false")

	w.Write(json)
}

func getAllScores(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entro a getAllScores")
	db := GetConnection()
	rows, err := db.Query("SELECT * FROM snake.users ORDER BY score DESC")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	defer rows.Close()
	var userScores []User
	for rows.Next() {
		u := User{}
		var id int
		var score int
		var name string
		if err := rows.Scan(&id, &name, &score); err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		u.ID = id
		u.Name = name
		u.Score = score
		userScores = append(userScores, u)
	}
	json, err := json.Marshal(userScores)
	err = rows.Err()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(json)
}

func main() {
	db := GetConnection()
	fmt.Println("connecting...")
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS snake.users (id INT PRIMARY KEY, name STRING, score INT);"); err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Get("/scores", getAllScores)
	r.Post("/add-user/", addUser)
	log.Fatal(http.ListenAndServe(":8002", r))
}
