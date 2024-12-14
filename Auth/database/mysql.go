package database

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB initializes MySQL database connection
func InitDB() {
    // Connect to the MySQL database
    var err error
    db, err = sql.Open("mysql", "root:IsaceL318@tcp(127.0.0.1:3306)/microservices")
    if err != nil {
        log.Fatal(err)
    }

    // Check database connection
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    log.Println("Connected to MySQL database")
}

// CreateUser inserts a new user into the database
func CreateUser(username, password string) error {
    // Prepare SQL statement
    stmt, err := db.Prepare("INSERT INTO users (username, upassword) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute SQL statement
    _, err = stmt.Exec(username, password)
    if err != nil {
        return err
    }

    return nil
}

// ValidateUser checks if the provided username and password match
func ValidateUser(username, password string) bool {
    // Prepare SQL statement
    stmt, err := db.Prepare("SELECT COUNT(*) FROM users WHERE username = ? AND upassword = ?")
    if err != nil {
        log.Println("Error preparing SQL statement:", err)
        return false
    }
    defer stmt.Close()

    // Execute SQL statement
    var count int
    err = stmt.QueryRow(username, password).Scan(&count)
    if err != nil {
        log.Println("Error querying database:", err)
        return false
    }

    return count == 1
}
