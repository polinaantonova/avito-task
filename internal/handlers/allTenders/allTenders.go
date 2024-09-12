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

	if serviceType == "" {
		query := `
WITH max_versions AS (
    SELECT 
        id, 
        MAX(version) AS max_version
    FROM 
        tenders
    GROUP BY 
        id
)
SELECT 
    t.id, 
    t.name, 
    t.description, 
    t.service_type, 
    t.status, 
    t.organization_id, 
    t.creator_username, 
    t.created_at, 
    t.version
FROM 
    tenders t
JOIN 
    max_versions mv
    ON t.id = mv.id 
    AND t.version = mv.max_version
ORDER BY 
    t.name
LIMIT 
    $1 
OFFSET 
    $2;
`

		rows, err := aT.db.Query(query, limit, offset)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			myTender := tender.NewTender()
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
		query := `SELECT t.id, t.name, t.description, t.service_type, t.status, t.organization_id, t.creator_username, t.created_at, t.version
		FROM tenders t
		JOIN (
			SELECT id, MAX(version) as max_version
		FROM tenders
		GROUP BY id
		) subquery
		ON t.id = subquery.id AND t.version = subquery.max_version
		WHERE t.service_type = $1
		ORDER BY t.name
		LIMIT $2 OFFSET $3;`

		rows, err := aT.db.Query(query, serviceType, limit, offset)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			myTender := tender.NewTender()
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
