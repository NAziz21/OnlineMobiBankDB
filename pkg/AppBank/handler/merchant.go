package handler

import (
	"BankVersion3/pkg/AppBank/model"
	"BankVersion3/pkg/AppBank/service"
	"BankVersion3/pkg/helpers"
	"BankVersion3/pkg/middleware"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type MerchantRA struct {
	MerchantService *service.MerchantService
}

func MerchantRestApi(clService *service.MerchantService) *MerchantRA {
	return &MerchantRA{
		MerchantService: clService,
	}
}

// RegisterMerchant Регистрация нового продавца
func (clDB *MerchantRA) RegisterMerchant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var RequestBody model.Merchant
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Merchant! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newRequest,status, err := ValidateRegisterMerchant(&RequestBody)
	if err != nil {
		log.Println("Cannot Validate Register Merchant(HP)!")
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if status != true {
		log.Println(" Validation of adding merchant was failed(HP)!")
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseBody, err := clDB.MerchantService.RegisterMerchant(&newRequest)
	if err != nil {
		log.Println("Register of Merchant is failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if responseBody == nil {
		log.Println("Register of Merchant is failed! Nil(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}


func ValidateRegisterMerchant(requestBody *model.Merchant) (model.Merchant,bool,error) {
	var newRequest model.Merchant
	newRequest.CompanyName = helpers.TittleName(requestBody.CompanyName)
	if newRequest.CompanyName == "" || len(newRequest.CompanyName) <= 3 {
		fmt.Println("Merchant's CompanyName is too short (HP)")
		return newRequest, false, errors.New("CompanyName is too short")
	}
	newRequest.CompanyAddress = helpers.TittleName(requestBody.CompanyAddress)
	if newRequest.CompanyAddress == "" || len(newRequest.CompanyAddress) <= 3 {
		fmt.Println("Merchant's CompanyAddress is too short (HP)")
		return newRequest, false, errors.New("CompanyAddress is too short")
	}
	newRequest.OwnerName = helpers.TittleName(requestBody.OwnerName)
	if newRequest.OwnerName == "" || len(newRequest.OwnerName) <= 3 {
		fmt.Println("Merchant's OwnerName is too short (HP)")
		return newRequest, false, errors.New("OwnerName is too short")
	}
	newRequest.OwnerSurname = helpers.TittleName(requestBody.OwnerSurname)
	if newRequest.OwnerSurname == "" || len(newRequest.OwnerSurname) <= 3 {
		fmt.Println("Merchant's Owner Surname is too short or nil! (HP)")
		return newRequest, false, errors.New("Owner Surname!")
	}
	newRequest.Login = requestBody.Login
	if len(newRequest.Login) <= 6 || newRequest.Login == "" {
		fmt.Println("Merchant's Login is too short or nil! (HP)")
		return newRequest, false, errors.New("Merchant's Login: Login is too short or nil!")
	}

	newRequest.Password = requestBody.Password
	if len(newRequest.Password) <= 6 || newRequest.Password == ""{
		fmt.Println("Merchant's Password is too short(HP)")
		return newRequest, false, errors.New("Password: Password is too short or nil!")
	}
	// Validate()
	newRequest.CompanyDescription = helpers.TittleName(requestBody.CompanyDescription)
	if requestBody.CompanyDescription == "" && len(requestBody.CompanyDescription) <= 2 {
		fmt.Println("Company Description")
		return newRequest, false, errors.New("Company Description:  is too short or nil!")

	}
	return newRequest,true,nil
}


// LoginMerchant Логинирование
func (clDB *MerchantRA) LoginMerchant(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var RequestBody model.MerchantLogin
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		log.Println("Cannot get information from Client! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = clDB.MerchantService.LoginMerchant(RequestBody.Login,RequestBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(130 * time.Minute)),
		ID:        RequestBody.Login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(service.MySigningKey)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Token ErroR!")
		return
	}
	err = json.NewEncoder(w).Encode(ss)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// MerchantBankcards Счет Продавца!
func (clDB *MerchantRA) MerchantBankcards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())
	merchantID, err := clDB.MerchantService.GetMerchantID(token.ID)
	if err != nil {
		fmt.Printf("Error getting MerchantID! MerchantBankcards(HP)! Error:%e", err)
		return
	}
	clientAccount, err := clDB.MerchantService.MerchantBankcards(merchantID)
	if err != nil {
		log.Println("Cannot get a Merchant's Bank cards(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(clientAccount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}


// ListOfATMsM Список Банкоматов!
func (clDB *MerchantRA) ListOfATMsM(w http.ResponseWriter, r * http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	listATMs, err := clDB.MerchantService.ListOfATMsM()
	if err != nil {
		log.Println("List of ATMS for a Merchant(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(listATMs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}


//NewService Добавление нового сервиса!
func (clDB *MerchantRA) NewService(w http.ResponseWriter, r *http.Request)  {
	token := middleware.ReadTokenFromContext(r.Context())
	merchantID, err := clDB.MerchantService.GetMerchantID(token.ID)
	if err != nil {
		fmt.Printf("Error getting ClientID(HP)! Error:%e", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var RequestBody model.ServiceStruct
	err = json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Client (HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newRequestBody,status, err := ServiceReqValidate(&RequestBody, int(merchantID))
	if err != nil {
		log.Println("Cannot validate request information from Client (HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if status != true {
		log.Println("Validation of adding new service failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var responseBody *model.ClientsAccountResponse
	responseBody,err = clDB.MerchantService.NewService(int(merchantID),&newRequestBody)
	if err != nil {
		log.Println("Cannot add a new service Merchant!(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}


// ServiceReqValidate Проверка запроса на валидацию
func ServiceReqValidate(requestBody *model.ServiceStruct, merchantID int) (model.ServiceStructInt,bool,error) {

	var newRequestBody model.ServiceStructInt

	newRequestBody.ServiceName = helpers.TittleName(requestBody.ServiceName)
	if requestBody.ServiceName == "" || len(requestBody.ServiceName) <= 4 {
		fmt.Println("Service Name")
		return newRequestBody, false, errors.New("Service Name is incorrect!")
	}

	newRequestBody.ServiceCategory = helpers.TittleName(requestBody.ServiceCategory)
	if newRequestBody.ServiceCategory == "" || len(newRequestBody.ServiceCategory) <= 1  {
		fmt.Println("Service Category")
		return newRequestBody, false, errors.New("Service Category is incorrect!")
	}
	newRequestBody.CompanyServiceName = helpers.TittleName(requestBody.CompanyServiceName)
	if newRequestBody.CompanyServiceName == "" || len(newRequestBody.CompanyServiceName) <= 4  {
		fmt.Println("Company Name")
		return newRequestBody,false, errors.New("Company Name is incorrect!")
	}

	newRequestBody.ServicePrice, _ = helpers.StrInt(requestBody.ServicePrice)

	if newRequestBody.ServicePrice == 0{
		fmt.Println("Service Price")
		return newRequestBody,false, errors.New("Service price is null!")
	}

	newRequestBody.Currency = strings.ToUpper(requestBody.Currency)

	if newRequestBody.Currency !="TJS" && newRequestBody.Currency != "RUB" && newRequestBody.Currency!= "USD" {
		fmt.Println("Currency")
		return newRequestBody, false, errors.New("Wrong Currency!")
	}
	 newRequestBody.Type = helpers.TittleName(requestBody.Type)
	if newRequestBody.Type !="Subscribe Payment" && newRequestBody.Type != "One Time Payment"  {
		fmt.Println("Type of Payment")
		return newRequestBody,false, errors.New("Wrong Type!")
	}
	return newRequestBody,true,nil
}


//AddBankCard Добавление карты продавца!
func (clDB *MerchantRA) AddBankCard(w http.ResponseWriter, r*http.Request)  {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())
	merchantID, err := clDB.MerchantService.GetMerchantID(token.ID)
	if err != nil {
		fmt.Printf("Error getting ClientID(HP)! Error:%e", err)
		return
	}
	var RequestBody model.BankCardM
	err = json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Client! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newRequest, err := ValidateBankCardM(&RequestBody)
	if err != nil {
		log.Println("Validation of Bank Card failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message, err := clDB.MerchantService.AddBankCard(merchantID,&newRequest)
	if err != nil {
		log.Println("Process of adding Bank Card to Client is failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if message == nil {
		log.Println("Status of Bank Card failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}


//ValidateBankCardM Валидация!
func ValidateBankCardM(requestBody *model.BankCardM) (model.BankCardM, error) {
	var newRequest model.BankCardM

	newRequest.IBAN = helpers.PhoneNumber(requestBody.IBAN)
	if len(newRequest.IBAN) != 20 || newRequest.IBAN == "" {
		log.Println("Client's number card is too short or nil(HP)!")
		return newRequest, errors.New("Bank number card is too short!")
	}
	newRequest.INN = helpers.PhoneNumber(requestBody.INN)
	if len(newRequest.INN) != 12 || newRequest.INN == "" {
		log.Println("Client's number card is too short or nil(HP)!")
		return newRequest, errors.New("Bank number card is too short!")
	}
	newRequest.ORGN = helpers.PhoneNumber(requestBody.ORGN)
	if len(newRequest.ORGN) != 13 || newRequest.ORGN == "" {
		log.Println("Client's number card is too short or nil(HP)!")
		return newRequest, errors.New("Bank number card is too short!")
	}
	return newRequest, nil
}

