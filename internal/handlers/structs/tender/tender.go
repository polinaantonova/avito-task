package tender

import (
	"errors"
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
		Id:        uuid.New().String(),
		Version:   1,
		Status:    "Created",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

func (t *Tender) ValidateTenderServiceType() error {
	var validServiceType = [3]string{"Construction", "Delivery", "Manufacture"}
	for _, service := range validServiceType {
		if t.ServiceType == service {
			return nil
		}
	}
	return errors.New("choose correct service type option Construction, Delivery, Manufacture")
}

func (t *Tender) ValidateStatus() error {
	var validStatus = [3]string{"Created", "Published", "Closed"}
	for _, status := range validStatus {
		if t.ServiceType == status {
			return nil
		}
	}
	return errors.New("choose correct status option: Created, Published, Closed")
}

func (t *Tender) ValidateStringFieldsLen() error {
	if len(t.Description) > 500 {
		return errors.New("tender description too long")
	}
	if len(t.Id) > 100 || len(t.Name) > 100 || len(t.OrganizationId) > 100 || len(t.CreatorUsername) > 100 {
		return errors.New("one or more too long text fields")
	}
	return nil
}
