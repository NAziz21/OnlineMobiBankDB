package handler

import (
	"testing"
	"time"
)


type ServiceStruct struct {
	ServiceName        string    `json:"service_name"`
	ServiceCategory    string    `json:"service_category"`
	CompanyServiceName string    `json:"company_service_name"`
	ServicePrice       string    `json:"service_price"`
	Currency		   string	 `json:"currency"`
	ServiceCreateData  time.Time `json:"service_create_data"`
	UpdateCreateData   time.Time `json:"update_create_data"`
}


func TestMerchantRA_NewService(t *testing.T) {

}
