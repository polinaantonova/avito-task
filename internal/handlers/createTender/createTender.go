package createTender

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}

	myTender := tender.NewTender()

	err = json.Unmarshal(body, &myTender)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	err = myTender.ValidateTenderServiceType()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = myTender.ValidateStringFieldsLen()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = myTender.ValidateUserCreation(tC.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sqlStatement := `
WITH UserId AS (
    SELECT id
    FROM employee
    WHERE username = $1
),

-- Шаг 2: Найти organization_id из таблицы organization_responsible
OrgId AS (
    SELECT organization_id
    FROM organization_responsible
    WHERE user_id = (SELECT id FROM UserId)
)

-- Шаг 3: Вставить данные в таблицу tenders
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
		http.Error(w, "cannot insert query in tenders table", http.StatusBadRequest)
	}
	response, err := json.Marshal(myTender)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	w.Write([]byte("\ntender added to database"))
}
