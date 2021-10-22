package model

import "time"

type MerchantLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Merchant struct {
	CompanyName        string    `json:"company_name"`
	CompanyAddress     string    `json:"company_address"`
	CompanyDescription string    `json:"company_description"`
	OwnerName          string    `json:"owner_name"`
	OwnerSurname       string    `json:"owner_surname"`
	Login              string    `json:"login"`
	Password           string    `json:"password"`
	CreateData         time.Time `json:"create_data"`
	UpdateData         time.Time `json:"update_data"`
}

type ResponseToMerchant struct {
	Status         bool   `json:"status"`
	Message        string `json:"message"`
	CompanyName    string `json:"company_name"`
	CompanyAddress string `json:"company_address"`
}

type ResponseMerchantAccount struct {
	MerchantID       int    `json:"client_id"`
	AccountName    string `json:"account_name"`
	NetworkPayment string `json:"network_payment"`
	Balance        int    `json:"balance"`
	Currency       string `json:"currency"`
}

type ServiceStruct struct {
	ServiceName        string    `json:"service_name"`
	ServiceCategory    string    `json:"service_category"`
	CompanyServiceName string    `json:"company_service_name"`
	ServicePrice       string    `json:"service_price"`
	Currency		   string	 `json:"currency"`
	Type	    	   string	 `json:"type"`
	ServiceCreateData  time.Time `json:"service_create_data"`
	UpdateCreateData   time.Time `json:"update_create_data"`
}

type ServiceStructInt struct {
	ServiceName        string    `json:"service_name"`
	ServiceCategory    string    `json:"service_category"`
	CompanyServiceName string    `json:"company_service_name"`
	ServicePrice       int   	 `json:"service_price"`
	Currency		   string	 `json:"currency"`
	Type	    	   string	 `json:"type"`
	ServiceCreateData  time.Time `json:"service_create_data"`
	UpdateCreateData   time.Time `json:"update_create_data"`
}

// Добавление нового счета для продавца
type BankCardM struct {
	IBAN		string `json:"iban"`
	INN  		string `json:"inn"`
	ORGN		string `json:"orgn"`
}

type InfoBankCardM struct {
	BankName 		 string `json:"bank_name"`
	IBAN         	 string `json:"IBAN"`
	INN		    	 string `json:"inn"`
	ORGN          	 string `json:"orgn"`
	Currency 		 string `json:"currency"`
	CreateData     time.Time `json:"create_data"`
	UpdateData     time.Time `json:"update_data"`
}
// // //

// Show merchant's accounts

type ResponseMerchantBankAccounts struct {
	BankName       string `json:"bank_name"`
	Balance        int    `json:"balance"`
	Currency       string `json:"currency"`
}