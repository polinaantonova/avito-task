package tender

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"polina.com/m/internal/validations"
)

func ValidateTenderServiceType(service string) error {
	var validServiceType = [3]string{"Construction", "Delivery", "Manufacture"}
	for _, validService := range validServiceType {
		if service == validService {
			return nil
		}
	}
	return errors.New("choose correct service type option Construction, Delivery, Manufacture")
}

func ValidateStatus(status string) error {
	validStatuses := [3]string{"Created", "Published", "Closed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
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

func (t *Tender) ValidateUserCreation(db *sql.DB) error {

	if t.CreatorUsername == "" {
		return errors.New("please specify creator username")
	}

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
	err := db.QueryRow(query, t.CreatorUsername).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found or doesn't belong to org responsible")
		} else {
			return errors.New("cannot execute query")
		}
	}
	return nil
}

func NewStatusUpdateValidator() *validations.StatusUpdateValidator {
	return &validations.StatusUpdateValidator{}
}
