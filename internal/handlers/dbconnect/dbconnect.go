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
	"strconv"

	_ "github.com/lib/pq"
)

type DBConnector struct {
	db *sql.DB
}

func (d *DBConnector) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var result int
	err := d.db.QueryRow("SELECT 1;").Scan(&result)
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

func NewDBConnector(db *sql.DB) *DBConnector {
	return &DBConnector{db: db}
}
