package tenderStatus

import (
	"fmt"
	"net/http"
	"polina.com/m/internal/tender"
)

type TenderStatus struct {
	tenders       *tender.TenderList
	currentTender *tender.Tender
}

func (tC *TenderStatus) CurrentTender() *tender.Tender {
	return tC.currentTender
}

func NewTenderStatus(tenders *tender.TenderList) *TenderStatus {
	return &TenderStatus{
		tenders: tenders,
	}
}

func (tS *TenderStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		fmt.Println(r.URL)

	}
	w.WriteHeader(http.StatusOK)
	//w.Write(response)

}
