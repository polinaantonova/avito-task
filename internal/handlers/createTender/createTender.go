package createTender

import (
	"encoding/json"
	"io"
	"net/http"
	"polina.com/m/internal/handlers/structs/tender"
)

type TenderCreator struct{}

func NewTenderCreator() *TenderCreator {
	return &TenderCreator{}
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

	tender := tender.NewTender()

	err = json.Unmarshal(body, &tender)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	err = tender.ValidateTenderServiceType()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tender.ValidateStringFieldsLen()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(tender)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
