package tender

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"

	_ "github.com/lib/pq"
)

func (t *Tender) ValidateTenderServiceType() error {
	var validServiceType = [3]string{"Construction", "Delivery", "Manufacture"}
	for _, service := range validServiceType {
		if t.ServiceType == service {
			return nil
		}
	}
	return errors.New("choose correct service type option Construction, Delivery, Manufacture")
}

func (t *Tender) ValidateStatus() error {
	var validStatus = [3]string{"Created", "Published", "Closed"}
	for _, status := range validStatus {
		if t.ServiceType == status {
			return nil
		}
	}
	return errors.New("choose correct status option: Created, Published, Closed")
}

func (t *Tender) ValidateStringFieldsLen() error {
	if len(t.Description) > 500 {
		return errors.New("tender description too long")
	}
	if len(t.Id) > 100 || len(t.Name) > 100 || len(t.OrganizationId) > 100 || len(t.CreatorUsername) > 100 {
		return errors.New("one or more too long text fields")
	}
	return nil
}

func (t *Tender) ValidateUser() error {

	if t.CreatorUsername == "" {
		return errors.New("please specify creator username")
	}

	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {

		errorText := fmt.Sprintf("empty env variables\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
		return errors.New(errorText)

		//тестирую локально
		//user = "polina"
		//password = "1234"
		//host = "localhost"
		//port = "5432"
		//dbName = "avito-task"
	}

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		errorText := fmt.Sprintf("cannot connect to database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
		return errors.New(errorText)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		errorText := fmt.Sprintf("cannot ping database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
		return errors.New(errorText)
	}

	fmt.Println("Successfully connected!")

	query := `WITH employee_id AS (
    SELECT id
    FROM employee
    WHERE username = $1
)
SELECT e.id
FROM employee_id e
JOIN organization_responsible o 
ON e.id = o.user_id
LIMIT 1;`

	/*
	   select 1
	   from employee e
	   join organization_responsible o on e.id = o.user_id
	   where e.username = $1
	   limit 1;

	*/

	var userID uuid.UUID
	err = db.QueryRow(query, t.CreatorUsername).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found or doesn't belong to org responsible")
		} else {
			return errors.New("cannot execute query")
		}
	}
	return nil
}
