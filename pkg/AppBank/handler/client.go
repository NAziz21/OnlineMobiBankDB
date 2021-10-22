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
	"time"
)

type ClientRA struct {
	ClientService *service.ClientService
}

func ClientRestAPI(clService *service.ClientService) *ClientRA {
	return &ClientRA{
		ClientService: clService,
	}
}

// Register Регистрация: Получение данных пользователя!
func (clDB *ClientRA) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var RequestBody model.Client
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Client (HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newRequest, err := ValidateRegisterClient(&RequestBody)
	if err != nil {
		log.Println("Validation of Client's request is failed(HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseBody, err := clDB.ClientService.Register(&newRequest)
	if err != nil {
		log.Println("Register of a Client Failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if responseBody == nil {
		log.Println("Registering of client failed!(HP)")
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

// ValidateRegisterClient Проверка полученных данных от клиента!
func ValidateRegisterClient(requestBody *model.Client) (model.Client, error) {
	var newRequest model.Client
	newRequest.Name = helpers.TittleName(requestBody.Name)
	if newRequest.Name == "" || len(newRequest.Name) <= 2 {
		fmt.Println("Client's name is too short (HP)")
		return newRequest, errors.New("Name is too short")
	}
	newRequest.Surname = helpers.TittleName(requestBody.Surname)
	if newRequest.Surname == "" || len(newRequest.Surname) <= 2 {
		fmt.Println("Client's surname is too short (HP)")
		return newRequest, errors.New("Surname is too short")
	}

	newRequest.Login = requestBody.Login
	if newRequest.Login == "" || len(newRequest.Login) <= 3 {
		fmt.Println("Client's login is too short (HP)")
		return newRequest, errors.New("Login is too short")
	}
	newRequest.Password = requestBody.Password
	if len(newRequest.Password) <= 6 || newRequest.Password == "" {
		fmt.Println("Client's Password is too short(HP)")
		return newRequest, errors.New("Password: Password is too short or nil!")
	}
	newRequest.Phone = helpers.PhoneNumber(requestBody.Phone)
	if len(requestBody.Phone) < 12 || requestBody.Phone == "" {
		fmt.Println("Client's phone is too short (HP)")
		return newRequest, errors.New("Phone number is too short")
	}
	return newRequest, nil
}

// Login Логин: Получение Логина и Пароля Клиента!
func (clDB *ClientRA) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var RequestBody model.ClientLogin
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Client! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = clDB.ClientService.Login(RequestBody.Login, RequestBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Login of a Client Failed(HP)! Error:", err)
		log.SetOutput(w)
		return
	}

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

//AddBankCard Добавление карты клиента!
func (clDB *ClientRA) AddBankCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	token := middleware.ReadTokenFromContext(r.Context())

	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		fmt.Printf("Error getting ClientID(HP)! Error:%e", err)
		return
	}

	var RequestBody model.BankCard
	err = json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Client! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newRequest, err := ValidateBankCard(&RequestBody)
	if err != nil {
		log.Println("Validation of Bank Card failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message, err := clDB.ClientService.AddBankCard(clientID, &newRequest)
	if err != nil {
		log.Println("Process of adding Bank Card to Client is failed(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	/*if message == nil {
		log.Println("Status of Bank Card failed(HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

//ValidateBankCard Валидация!
func ValidateBankCard(requestBody *model.BankCard) (model.BankCard, error) {
	var newRequest model.BankCard

	newRequest.NumberCard = helpers.PhoneNumber(requestBody.NumberCard)
	if len(newRequest.NumberCard) != 16 || newRequest.NumberCard == "" {
		fmt.Println("Client's number card is too short or nil(HP)!")
		return newRequest, errors.New("Bank number card is too short!")
	}
	newRequest.Name = helpers.TittleName(requestBody.Name)
	if len(newRequest.Name) <= 2 || newRequest.Name == "" {
		fmt.Println("Client's name is too short or nil(HP)!")
		return newRequest, errors.New("Name is too short!")
	}
	newRequest.Surname = helpers.TittleName(requestBody.Surname)
	if len(newRequest.Surname) <= 2 || newRequest.Surname == "" {
		fmt.Println("Client's surname is too short or nil(HP)!")
		return newRequest, errors.New("Surnmae is too short!")
	}
	newRequest.CVC = helpers.PhoneNumber(requestBody.CVC)
	if len(newRequest.CVC) != 3 || newRequest.CVC == "" {
		fmt.Println("CVC Client is too short or too long or nil(HP)!")
		return newRequest, errors.New("CVC is too short or too long!")
	}
	return newRequest, nil
}

// TransferMoneyMobi Транзакция между клиентами банков по номеру телефону!

func (clDB *ClientRA) TransferMoneyMobi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	token := middleware.ReadTokenFromContext(r.Context())
	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		fmt.Printf("Error getting ClientID(HP)! Error:%e", err)
		return
	}

	var RequestBody model.TransactionMobiPhone
	err = json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		fmt.Println("Cannot get information from Client(HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newRequest, err := ValidateTransferMobi(&RequestBody) // Валидация запроса клиента
	if err != nil {
		log.Println("Validation of Transfer Mobi: Cannot get information from Client(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	saved, err := clDB.ClientService.TransferMoneyMobi(clientID, &newRequest) //
	if err != nil {
		log.Println("Cannot Transfer Money! (HP)")
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if saved == nil {
		log.Println("Cannot Transfer Money! (HP)")
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(saved)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

//ValidateTransferMobi трансфер по мобильному телефону!
func ValidateTransferMobi(requestBody *model.TransactionMobiPhone) (model.TransactionMobiPhone, error) {
	var newRequest model.TransactionMobiPhone

	newRequest.PhoneNumber = helpers.PhoneNumber(requestBody.PhoneNumber)
	if len(newRequest.PhoneNumber) != 12 || newRequest.PhoneNumber == "" {
		fmt.Println("Client's number card is too short or nil(HP)!")
		return newRequest, errors.New("Bank number card is too short!")
	}

	newRequest.Amount = requestBody.Amount

	return newRequest, nil
}

// ClientsBankcards Счет Клиента!
func (clDB *ClientRA) ClientsBankcards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())
	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		log.Printf("Error getting ClientID(HP)! Error:%e", err)

		return
	}
	clientAccount, err := clDB.ClientService.ClientsBankcards(clientID)
	if err != nil {
		log.Println("Cannot get a Client's Bank cards! (HP)")
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

// ListOfATMs Список Банкоматов!
func (clDB *ClientRA) ListOfATMs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	listATMs, err := clDB.ClientService.ShowATMList()
	if err != nil {
		log.Println("Cannot get a list of ATMS! (HP)")
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

// ListOfServices Список всех услуг!
func (clDB *ClientRA) ListOfServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	listOfServices, err := clDB.ClientService.ListOfServices()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error in watching all lists!", err)
		log.SetOutput(w)
		return
	}

	err = json.NewEncoder(w).Encode(listOfServices)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

//ListOfClientSubServices Список подключенных сервисов клиента
func (clDB *ClientRA) ListOfClientSubServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())

	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		log.Printf("Error getting ClientID(HP)! Error:%e", err)
		log.SetOutput(w)
		return
	}

	listSubService, err := clDB.ClientService.ListOfClientSubServices(clientID)
	if err != nil {
		log.Printf("Cannot getting a client's sub service list(HP)! Error:%e", err)
		log.SetOutput(w)
		return
	}

	err = json.NewEncoder(w).Encode(listSubService)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// SubscribeService Добавить услугу клиенту
func (clDB *ClientRA) SubscribeService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())
	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		log.Printf("Error getting ClientID(HP)! Error:%e", err)
		log.SetOutput(w)
		return
	}
	var RequestBody model.SubscribeServiceRequest
	err = json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		log.Println("Cannot get information from Client! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newRequest, err := validateSubscribeService(clDB, clientID, &RequestBody)
	if err != nil {
		log.Println("Validation Subscribe service(HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := clDB.ClientService.SubscribeServiceInf(clientID, &newRequest)
	if err != nil {
		log.Println("Cannot add a service to a Client(HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if message == nil {
		log.Println("Purchasing Subscribe - Failed (HP)! Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func validateSubscribeService(clDB *ClientRA, clientID int, requestBody *model.SubscribeServiceRequest) (model.SubscribeServiceRequest, error) {
	var newRequest model.SubscribeServiceRequest
	newRequest.ServiceName = helpers.TittleName(requestBody.ServiceName)
	newRequest.CompanyName = helpers.TittleName(requestBody.CompanyName)

	if newRequest.ServiceName == "" || len(newRequest.ServiceName) <= 3 {
		fmt.Println("Service name is too short or None!")
		return newRequest, errors.New("Error! Service name is too short or None!")
	}

	if newRequest.CompanyName == "" || len(newRequest.CompanyName) <= 3 {
		log.Println("Company name is too short or None!")
		return newRequest, errors.New("Error! Company name is too short or None!")
	}

	status, err := clDB.ClientService.CheckServiceName(newRequest.ServiceName, newRequest.CompanyName, clientID)
	if err != nil {
		log.Println("Cannot check service name in DB(HP)!")
		return newRequest, err
	}
	if status != true {
		log.Println("Such a Service name doesn't exists In DB(HP)!")
		return newRequest, errors.New("Such login exists!")
	}
	return newRequest, nil
}

// AddAmountFromCards пополняем кошелек с банковских карт!
func (clDB *ClientRA) AddAmountFromCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())
	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		log.Printf("Error getting ClientID(HP)! Error:%e", err)
		log.SetOutput(w)
		return
	}

	var RequestBody model.TransferFromCardToOnlineMobi
	err = json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		log.Println("Cannot get information from Client(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	responseBody, err := clDB.ClientService.AddAmountFromCards(clientID, &RequestBody)
	if err != nil {
		log.Println("Cannot add a service to a Client(HP)! Error:", err)
		log.SetOutput(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

//HistoryOfClientsActions История действий клиента!
func (clDB *ClientRA) HistoryOfClientsActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	token := middleware.ReadTokenFromContext(r.Context())

	clientID, err := clDB.ClientService.GetClientID(token.ID)
	if err != nil {
		log.Printf("Error getting ClientID(HP)! Error:%e", err)
		log.SetOutput(w)
		return
	}

	historyOfActions, err := clDB.ClientService.HistoryOfClientsActions(clientID)
	if err != nil {
		log.Printf("Cannot getting a client's sub service list(HP)! Error:%e", err)
		log.SetOutput(w)
		return
	}

	err = json.NewEncoder(w).Encode(historyOfActions)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
