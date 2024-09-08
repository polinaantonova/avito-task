//- `POSTGRES_USERNAME` — имя пользователя для подключения к PostgreSQL.
//- `POSTGRES_PASSWORD` — пароль для подключения к PostgreSQL.
//- `POSTGRES_HOST` — хост для подключения к PostgreSQL (например, localhost).
//- `POSTGRES_PORT` — порт для подключения к PostgreSQL (например, 5432).
//- `POSTGRES_DATABASE` — имя базы данных PostgreSQL, которую будет использовать приложение.

package dbconnect

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type DBConnector struct{}

func (d DBConnector) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		http.Error(w, "cannot connect to database", http.StatusInternalServerError)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		http.Error(w, "cannot ping database", http.StatusInternalServerError)
	}

	fmt.Println("Successfully connected!")

	var result int
	err = db.QueryRow("SELECT 1;").Scan(result)
	if err != nil {
		http.Error(w, "cannot execute query", http.StatusInternalServerError)
	}

	if result == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(result)))
		return
	}
	http.Error(w, "cannot return response", http.StatusInternalServerError)
}

func NewDBConnector() *DBConnector {
	return &DBConnector{}
}
