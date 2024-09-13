package tendersByUser

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
	"strconv"
)

type TendersByUser struct {
	db *sql.DB
}

func NewTendersByUser(db *sql.DB) *TendersByUser {
	return &TendersByUser{
		db: db,
	}
}

func (tU *TendersByUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage.SendErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 5
	offset := 0

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	filteredTenders := tender.NewTenderList()

	if username == "" {
		errorMessage.SendErrorMessage(w, "please specify user", http.StatusUnauthorized)
		return
	}

	query := `SELECT t.id, t.name, t.description, t.service_type, t.status, t.organization_id, t.creator_username, t.created_at, t.version
FROM tenders t
JOIN (
    SELECT id, MAX(version) as max_version
    FROM tenders
    GROUP BY id
) subquery
ON t.id = subquery.id AND t.version = subquery.max_version
WHERE t.creator_username = $1
ORDER BY t.name
LIMIT $2 OFFSET $3;`

	rows, err := tU.db.Query(query, username, limit, offset)

	if err != nil {
		errorMessage.SendErrorMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	for rows.Next() {
		myTender := tender.NewTender()
		err = rows.Scan(&myTender.Id, &myTender.Name, &myTender.Description, &myTender.ServiceType, &myTender.Status, &myTender.OrganizationId, &myTender.CreatorUsername, &myTender.CreatedAt, &myTender.Version)

		if err != nil {
			errorMessage.SendErrorMessage(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filteredTenders.AddTender(myTender)
	}

	if err = rows.Err(); err != nil {
		errorMessage.SendErrorMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(filteredTenders.List()) == 0 {
		errorMessage.SendErrorMessage(w, "user doesn't exist or didn't create tenders", http.StatusUnauthorized)
		return
	}

	response, err := json.Marshal(filteredTenders.List())
	if err != nil {
		errorMessage.SendErrorMessage(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
