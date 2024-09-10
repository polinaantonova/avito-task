package tenderStatus

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"polina.com/m/internal/tender"
	"strconv"
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

func (tS *TenderStatus) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, _ = strconv.Atoi(ps.ByName("tenderID"))

	if r.Method == http.MethodGet {
		fmt.Println(r.URL)

	}
	w.WriteHeader(http.StatusOK)
	//w.Write(response)

}
