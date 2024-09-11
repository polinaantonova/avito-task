package tendersByUser

import (
	"encoding/json"
	"net/http"
	"polina.com/m/internal/tender"
	"sort"
	"strconv"
)

type TendersByUser struct {
	tenders *tender.TenderList
}

func NewTendersByUser(tenders *tender.TenderList) *TendersByUser {
	return &TendersByUser{
		tenders: tenders,
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
	paginatedTenders := tender.NewTenderList()

	if username == "" {
		http.Error(w, "please specify user", http.StatusUnauthorized)
		return
	} else {
		for _, myTender := range tU.tenders.List() {
			if myTender.CreatorUsername == username {
				filteredTenders.AddTender(myTender)
			}
		}
	}

	if len(filteredTenders.List()) == 0 {
		http.Error(w, "user doesn't exist or didn't create tenders", http.StatusUnauthorized)
		return
	}

	sort.SliceStable(filteredTenders.List(), func(i, j int) bool {
		return filteredTenders.List()[i].Name < filteredTenders.List()[j].Name
	})

	start := offset
	end := offset + limit
	if start > len(filteredTenders.List()) {
		start = len(filteredTenders.List())
	}
	if end > len(filteredTenders.List()) {
		end = len(filteredTenders.List())
	}

	for _, filteredTender := range filteredTenders.List()[start:end] {
		paginatedTenders.AddTender(filteredTender)
	}

	response, err := json.Marshal(paginatedTenders.List())
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
