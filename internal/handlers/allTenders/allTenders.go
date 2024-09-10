package allTenders

import (
	"encoding/json"
	"net/http"
	"polina.com/m/internal/tender"
	"sort"
	"strconv"
)

type AllTenders struct {
	tenders *tender.TenderList
}

func NewAllTenders(tenders *tender.TenderList) *AllTenders {
	return &AllTenders{
		tenders: tenders,
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
	paginatedTenders := tender.NewTenderList()

	if serviceType != "" {
		for _, myTender := range aT.tenders.List() {
			if myTender.ServiceType == serviceType {
				filteredTenders.AddTender(myTender)
			}
		}
	} else {
		filteredTenders = aT.tenders
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
