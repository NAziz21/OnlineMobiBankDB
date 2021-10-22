package model

import "time"

const RubTjs = 6_35
const UsdRub = 71_40
const UsdTjs = 11_31

type ClientLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Client struct {
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	Phone      string    `json:"phone"`
	CreateData time.Time `json:"create_data"`
	UpdateData time.Time `json:"update_data"`
}

type ClientsAccount struct {
	AccountName    string    `json:"account_name"`
	NetworkPayment string    `json:"network_payment"`
	Balance        int64     `json:"balance"`
	Currency       string    `json:"currency"`
	CreateData     time.Time `json:"create_data"`
	UpdateData     time.Time `json:"update_data"`
}

type ClientsAccountRequest struct {
	AccountName string `json:"account_name"`
}

type ClientsAccountResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type ResponseClientsAccount struct {
	ClientID       int    `json:"client_id"`
	AccountName    string `json:"account_name"`
	NetworkPayment string `json:"network_payment"`
	Balance        int    `json:"balance"`
	Currency       string `json:"currency"`
}

type ResponseToClient struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}


// AddAmount OnlineMobi

type TransactionPhone struct {
	PhoneNumber     string `json:"phone"`
	Amount          string `json:"amount"`
	SAccountName    string `json:"s_account_name"`
	SNetworkPayment string `json:"s_network_payment"`
	SCurrency       string `json:"s_currency"`
	RAccountName    string `json:"r_account_name"`
	RNetworkPayment string `json:"r_network_payment"`
	RCurrency       string `json:"r_currency"`
}

type TransactionMobiPhone struct {
	PhoneNumber     string `json:"phone"`
	Amount          string `json:"amount"`
}

type TransactionTransferStruct struct {
	SClientID int `json:"sclient_id"`
	RClientID int `json:"rclient_id"`
	SBalance int `json:"s_balance"`
	RBalance int `json:"r_balance"`
	Amount int `json:"amount"`
}
// // //

// Show Client's Bank cards

type ResponseClientsBankcards struct {
	ID 			int    `json:"№"`
	BankName    string `json:"bankName"`
	BankcardName string `json:"bankcard_name"`
	PaymentSystem string `json:"payment_system"`
	Balance        int    `json:"balance"`
	Currency       string `json:"currency"`
}
// // //

// Transfer money from cards to OnlineMobi

type TransferFromCardToOnlineMobi struct {
	Amount       int `json:"amount"`
	NumberOfCard int `json:"number_of_card"`
}

type StructTransferFromCardToOnlineMobi struct {
	OnlineMobiCurrency string
	CardCurrency	   string
	OnlineMobiBalance  int
	CardBalance		   int
}

type ResponseTransactionPhone struct {
	Status          string `json:"status"`
	Amount          int    `json:"amount"`
	SAccountName    string `json:"s_account_name"`
	SNetworkPayment string `json:"s_network_payment"`
	SCurrency       string `json:"s_currency"`
	RAccountName    string `json:"r_account_name"`
	RNetworkPayment string `json:"r_network_payment"`
	RCurrency       string `json:"r_currency"`
}

type SubscribeServiceRequest struct {
	CompanyName string `json:"company_name"`
	ServiceName string `json:"service_name"`
}

type BankCard struct {
	NumberCard string `json:"number_card"`
	CVC        string `json:"cvc"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
}

type InfoBankCard struct {
	BankName 		 string `json:"bank_name"`
	BankcardName 	 string `json:"bankcard_name"`
	BinNumber    	 string `json:"bin_number"`
	CVC          	 string `json:"cvc"`
	PaymentSystem    string `json:"payment_system"`
	Currency 		 string `json:"currency"`
	CreateData     time.Time `json:"create_data"`
	UpdateData     time.Time `json:"update_data"`
}

