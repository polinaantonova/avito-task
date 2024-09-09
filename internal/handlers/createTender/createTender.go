package createTender

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"polina.com/m/internal/handlers/structs/tender"
)

type TenderCreator struct {
	tenders *tender.TenderList
}

func NewTenderCreator(tenders *tender.TenderList) *TenderCreator {

	return &TenderCreator{
		tenders: tenders,
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

	err = myTender.ValidateUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tC.tenders.AddTender(myTender)
	fmt.Println(tC.tenders.TenderList)

	response, err := json.Marshal(myTender)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
