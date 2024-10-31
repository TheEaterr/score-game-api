package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// album represents data about a record album.
type score struct {
	Player string `json:"player"`
	Score  int64  `json:"score"`
}

type scoreServer struct {
	db *sql.DB
}

func NewScoreServer() *scoreServer {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "./data/score.db")
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Connected to the SQLite database successfully.")

	// Create the scores table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scores (player TEXT, score INTEGER)")
	if err != nil {
		log.Println(err)
		return nil
	}

	return &scoreServer{db}
}

func (s *scoreServer) getAllScores(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	// query scores from the database
	rows, err := s.db.Query("SELECT * FROM scores ORDER BY score DESC LIMIT 10")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scores []score
	for rows.Next() {
		var score score
		if err := rows.Scan(&score.Player, &score.Score); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		scores = append(scores, score)
	}

	renderJSON(w, scores)
}

func (s *scoreServer) createScore(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling create task at %s\n", req.URL.Path)
	log.Printf("%s\n", req.Body)

	var score score
	if err := json.NewDecoder(req.Body).Decode(&score); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check player and score
	if score.Player == "" || score.Score == 0 {
		http.Error(w, "player and score are required", http.StatusBadRequest)
		return
	}

	// insert score into the database
	_, err := s.db.Exec("INSERT INTO scores (player, score) VALUES (?, ?)", score.Player, score.Score)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, score)
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	mux := http.NewServeMux()
	server := NewScoreServer()
	mux.HandleFunc("POST /score/add/", server.createScore)
	mux.HandleFunc("GET /score/list/", server.getAllScores)

	port := os.Getenv("SERVERPORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe("localhost:"+port, mux))
}
