package settings

import (
	"encoding/json"
	"net/http"
	"time"
)

type DB struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
}


type ClientFromFile struct {
	Clients []Clients `json:"clients"`
}


type Clients struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Phone    string  `json:"phone"`
	CreateData time.Time `json:"create_data"`
	UpdateData time.Time `json:"update_data"`
}

type MerchantFromFile struct {
	Merchants []Merchants `json:"merchants"`
}

type Merchants struct {
	CompanyName        string `json:"company_name"`
	CompanyAddress     string `json:"company_address"`
	CompanyDescription string `json:"company_description"`
	OwnerName          string `json:"owner_name"`
	OwnerSurname       string `json:"owner_surname"`
	Login              string `json:"login"`
	Password           string `json:"password"`
	CreateData         string `json:"create_data"`
	UpdateData         string `json:"update_data"`
}

// **********    Для Аккаунта Клиента и Продавца ***********////

type MerchantsAccountFromFile struct {
	Merchants []MerchantsAccount `json:"accounts"`
}

type MerchantsAccount struct {
	AccountName string `json:"account_name"`
	ClientID    *int64  `json:"client_id"`
	MerchantID  *int64  `json:"merchant_id"`
	Balance     int64  `json:"balance"`
	CreateData  string `json:"create_data"`
	UpdateData  string `json:"update_data"`
}
type ATMsFromFile struct {
	ATMs []ATMsStruct `json:"ATMs"`
}

type ATMsStruct struct {
	Address     string `json:"address"`
	Balance     int64  `json:"balance"`
	CreateData  string `json:"create_data"`
	UpdateData  string `json:"update_data"`
}

type MessageInfo struct {
	Status  bool
	Message	string
}


func Message(status string, message string) (map[string]interface{}) {
	return map[string]interface{}{"status": status,"message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}