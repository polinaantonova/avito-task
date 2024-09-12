package tendersByUser

import (
	"database/sql"
	"encoding/json"
	"net/http"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	myTender := tender.NewTender()

	if username == "" {
		http.Error(w, "please specify user", http.StatusUnauthorized)
		return
	}

	rows, err := tU.db.Query("SELECT * FROM tenders WHERE creator_username = $1 ORDER BY name LIMIT $2 OFFSET $3", username, limit, offset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&myTender.Id, &myTender.Name, &myTender.Description, &myTender.ServiceType, &myTender.Status, &myTender.OrganizationId, &myTender.CreatorUsername, &myTender.CreatedAt, &myTender.Version)

		if err != nil {
			http.Error(w, "cannot select tenders from table tenders", http.StatusInternalServerError)
			return
		}
		filteredTenders.AddTender(myTender)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(filteredTenders.List()) == 0 {
		http.Error(w, "user doesn't exist or didn't create tenders", http.StatusUnauthorized)
		return
	}

	response, err := json.Marshal(filteredTenders.List())
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
