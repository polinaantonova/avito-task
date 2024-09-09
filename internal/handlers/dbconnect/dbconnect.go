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

	_ "github.com/lib/pq"
)

type DBConnector struct{}

func (d *DBConnector) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")

	//	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
	//		user = "polina"
	//		password = "1234"
	//		host = "localhost"
	//		port = "5432"
	//		dbName = "avito-task"
	//	}

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	errorText := fmt.Sprintf("cannot connect to database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)

	if err != nil {
		http.Error(w, errorText, http.StatusInternalServerError)
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		http.Error(w, "cannot ping database", http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully connected!")

	var result int
	err = db.QueryRow("SELECT 1;").Scan(&result)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot execute query, %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if result == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(result)))
		return
	}
	http.Error(w, "cannot return response", http.StatusInternalServerError)
	return
}

func NewDBConnector() *DBConnector {
	return &DBConnector{}
}
