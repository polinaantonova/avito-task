package createTender

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
)

type TenderCreator struct {
	db *sql.DB
}

func NewTenderCreator(db *sql.DB) *TenderCreator {
	return &TenderCreator{
		db: db,
	}
}

func (tC *TenderCreator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorMessage.SendErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	myTender := tender.NewTender()
	myTender.Status = "Created"

	//проверить json
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&myTender); err != nil {
		errorMessage.SendErrorMessage(w, "check your JSON schema", http.StatusBadRequest)
		return
	}

	err := tender.ValidateTenderServiceType(myTender.ServiceType)
	if err != nil {
		errorMessage.SendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = myTender.ValidateStringFieldsLen()
	if err != nil {
		errorMessage.SendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = myTender.ValidateUserCreation(tC.db)
	if err != nil {
		errorMessage.SendErrorMessage(w, err.Error(), http.StatusForbidden)
		return
	}

	sqlStatement := `
WITH UserId AS (
    SELECT id
    FROM employee
    WHERE username = $1
),

OrgId AS (
    SELECT organization_id
    FROM organization_responsible
    WHERE user_id = (SELECT id FROM UserId)
)

INSERT INTO tenders (creator_username, id, name, description, service_type, organization_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
	$5,
    (SELECT organization_id FROM OrgId)
);
`
	_, err = tC.db.Exec(sqlStatement, myTender.CreatorUsername, myTender.Id, myTender.Name, myTender.Description, myTender.ServiceType)
	if err != nil {
		errorMessage.SendErrorMessage(w, "cannot insert query in tenders table", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(myTender)
	if err != nil {
		errorMessage.SendErrorMessage(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	w.Write([]byte("\ntender added to database"))
}
