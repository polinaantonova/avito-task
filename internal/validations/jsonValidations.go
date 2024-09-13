package validations

import "errors"

type StatusUpdateValidator struct {
	Username string `json:"username"`
	Status   string `json:"status"`
}

type TenderEditValidator struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"service_type"`
}

func NewTenderEditValidator() *TenderEditValidator {
	return &TenderEditValidator{}
}

func (tE *TenderEditValidator) ValidateStringFieldsLen() error {
	if len(tE.Description) > 500 {
		return errors.New("tender description too long")
	}
	if len(tE.Name) > 100 {
		return errors.New("one or more too long text fields")
	}
	return nil
}

type RollbackVersionValidator struct {
	Username string `json:"username"`
}

func NewRollbackVersionValidator() *RollbackVersionValidator {
	return &RollbackVersionValidator{}
}
