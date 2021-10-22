package model

import (
	"time"
)

type ATMsList struct {
	Address    string
	Balance    int64
	CreateData time.Time
	UpdateData time.Time
}

type ATMsListResponse struct {
	Address string
}

type ServiceList struct {
	ServiceName        string `json:"service_name"`
	ServiceCategory    string `json:"service_category"`
	CompanyServiceName string `json:"company_service_name"`
	ServicePrice       int    `json:"service_price"`
	Currency           string `json:"currency"`
	Type               string `json:"type"`
}

type SubscribeServiceList struct {
	ID                 int       `json:"№"`
	ServiceName        string    `json:"service_name"`
	ServiceCategory    string    `json:"service_category"`
	CompanyServiceName string    `json:"company_service_name"`
	ServicePrice       int       `json:"service_price"`
	Currency           string    `json:"currency"`
	Type               string    `json:"type"`
	CreateData         time.Time `json:"create_data"`
}

type HistoryList struct {
	AccountName       string `json:"account_name"`
	TransactionType   string `json:"transaction_type"`
	Debit             int `json:"debit"`
	Credit            int `json:"credit"`
	DataOfTransaction time.Time    `json:"data_of_transaction"`
}
