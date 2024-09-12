package allTenders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"polina.com/m/internal/tender"
	"strconv"
)

type AllTenders struct {
	db *sql.DB
}

func NewAllTenders(db *sql.DB) *AllTenders {
	return &AllTenders{
		db: db,
	}
}

func (aT *AllTenders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceType := r.URL.Query().Get("service_type")
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

	if serviceType == "" {
		rows, err := aT.db.Query("SELECT * FROM tenders ORDER BY name LIMIT $1 OFFSET $2", limit, offset)

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
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		rows, err := aT.db.Query("SELECT * FROM tenders WHERE service_type = $1 ORDER BY name LIMIT $2 OFFSET $3", serviceType, limit, offset)

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
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	response, err := json.Marshal(filteredTenders.List())
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
