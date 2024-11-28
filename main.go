package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var logger *log.Logger

func init() {
	// Initialize logger with timestamp
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Connect to MySQL database
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Database ping failed: %v", err)
	}
	logger.Println("Successfully connected to database")

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/current-time", getCurrentTime).Methods("GET")
	r.HandleFunc("/time-logs", getTimeLogs).Methods("GET")

	// Start the server
	logger.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getCurrentTime(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Received request for current time from %s", r.RemoteAddr)

	// Get the current time in Toronto
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		logger.Printf("Error loading timezone: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now().In(loc)

	// Log the current time to the database
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", currentTime)
	if err != nil {
		logger.Printf("Error logging time to database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Printf("Successfully logged time: %v", currentTime)

	// Return the current time in JSON format
	response := map[string]string{"current_time": currentTime.Format(time.RFC3339)}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Printf("Error encoding response: %v", err)
	}
}

func getTimeLogs(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Received request for time logs from %s", r.RemoteAddr)

	// Query all time logs from the database
	rows, err := db.Query("SELECT timestamp FROM time_log ORDER BY timestamp DESC LIMIT 100")
	if err != nil {
		logger.Printf("Error querying time logs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var timestamp time.Time
		if err := rows.Scan(&timestamp); err != nil {
			logger.Printf("Error scanning row: %v", err)
			continue
		}
		logs = append(logs, timestamp.Format(time.RFC3339))
	}

	// Return the logs in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		logger.Printf("Error encoding response: %v", err)
	}
}
