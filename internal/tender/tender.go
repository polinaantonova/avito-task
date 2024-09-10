package tender

import (
	"github.com/google/uuid"
	"time"
)

type Tender struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
	CreatedAt       string `json:"createdAt"`
	Version         int32  `json:"version"`
}

func NewTender() *Tender {
	return &Tender{
		Id:             uuid.New().String(),
		OrganizationId: uuid.New().String(),
		Version:        1,
		Status:         "Created",
		CreatedAt:      time.Now().Format(time.RFC3339),
	}
}

type TenderList struct {
	tenderList []*Tender
}

func NewTenderList() *TenderList {
	return &TenderList{
		tenderList: make([]*Tender, 0, 8),
	}
}

func (tL *TenderList) List() []*Tender {
	return tL.tenderList
}

func (tL *TenderList) AddTender(tender *Tender) {
	tL.tenderList = append(tL.tenderList, tender)
}
