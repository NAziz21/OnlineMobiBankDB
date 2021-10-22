package service

import (
	"BankVersion3/pkg/AppBank/model"
	"BankVersion3/pkg/AppBank/repository"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

var MySigningKey = []byte("MySecretCode")

/*
type TokenClaims struct {
	jwt.StandardClaims
	ClientId int `json:"user_id"`
}*/

type ClientService struct {
	ClientRep *repository.ClientRepository
}

func NewClientService(req *repository.ClientRepository) *ClientService {
	return &ClientService{
		ClientRep: req,
	}
}

// Register Создаем нового клиента!
func (clDB *ClientService) Register(clConf *model.Client) (*model.ResponseToClient, error) {
	login, err := clDB.ClientRep.CheckLoginRegister(clConf.Login)
	if err != nil {
		log.Printf("Cannot Register merchant to DB (SP)! Error:%e\n", err)
		return nil,err
	}
	if login == true {
		log.Printf("Such login exists in DB(SP)! Try another one! Error:%e\n", err)
		return nil,err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(clConf.Password), 14)
	if err != nil {
		return nil, err
	}
	clConf.Password = string(bytes)

	savedUser,err := clDB.ClientRep.Register(clConf)
	if err != nil {
		log.Printf("Register of client(SP)! Error:%e\n", err)
		return nil, err
	}

	return savedUser, nil
}


// Login Проверяем полученный Логин и Пароль с БД!
func (clDB *ClientService) Login(login, password string) (*model.ClientLogin, error) {
	client, err := clDB.ClientRep.Login(login)
	if err != nil {
		log.Printf("Login of Client(SP)! Error:%e\n", err)
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(password)); err != nil {
		log.Printf("Comparing of Hash Password(SP)! Error:%e\n", err)
		return nil, err
	}

	return client, nil
}


// GetClientID Получаем ID Клиента через его логин!
func (clDB *ClientService) GetClientID(login string) (clientID int, err error) {
	clientID, err = clDB.ClientRep.GetClientID(login)
	if err != nil {
		log.Printf("Can't Get Client ID(SP)! Error:%e\n", err)
		return clientID, err
	}
	return clientID, nil
}



// AddBankCard Добавление карты клиента!
func (clDB *ClientService) AddBankCard(clientID int, clConf *model.BankCard) (*model.ClientsAccountResponse,error) {


	checkBinNumber, err := clDB.ClientRep.CheckBinNumber(clConf,clientID) // Статус проверки Бин номера и сравнение имени и фамилии с БД!
	if err != nil {
		log.Printf("Cannot check a Client's card bin Number(SP)! Error:%e", err)
		return nil, err
	}

	if checkBinNumber == true {
		log.Println("checkBinNumber! Error! Wrong bin Number of card(RP)", err)
		return nil, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(clConf.CVC), 14) // Хэширование CVC
	if err != nil {
		return nil, err
	}
	clConf.CVC = string(bytes)


	newStruct,err := numberBankCard(clConf)  // После валидации запроса от клиента, получаем новую структуру!
	if err != nil {
		log.Printf("Verification of Client's card(SP)! Error:%e", err)
		return nil, err
	}

	checkCardName, err := clDB.ClientRep.CheckCardName(&newStruct,clientID)
	if err != nil {
		log.Printf("Such card has already added(SP)! Error:%e", err)
		return nil, err
	}

	if checkCardName == true {
		log.Printf("Such card has already added(SP)! Error:%e", err)
		return nil, err
	}

	_, err = clDB.ClientRep.AddBankCard(clientID,&newStruct)  // Добавление карты в БД!
	if err != nil {
		log.Printf("Cannot add a Client's card to DB (SP)! Error:%e", err)
		return nil, err
	}
	responseBody := model.ClientsAccountResponse{
		Status:  true,
		Message: "Success!",
	}
	return &responseBody,err
}


// numberBankCard Валидация имени карты!
func numberBankCard(clConf *model.BankCard) (model.InfoBankCard,error){
	var infoBank model.InfoBankCard
	infoBank.CVC = clConf.CVC
	infoBank.BinNumber = clConf.NumberCard
	infoBank.CreateData = time.Now()
	infoBank.UpdateData = time.Now()
	info,_ :=strconv.ParseInt(clConf.NumberCard[:6],10,0)

	//info2,_ := strconv.ParseInt(clConf.NumberCard[6:],10,0)
	switch info {
	case 450721:
		infoBank.BankName = "MDO HUMO"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "VISA"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 550721:
		infoBank.BankName = "MDO HUMO"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "MASTER CARD"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 650721:
		infoBank.BankName = "MDO HUMO"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "KORTI MILLI"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 451223:
		infoBank.BankName = "ALIF BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "VISA"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 551223:
		infoBank.BankName = "ALIF BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "MASTER CARD"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 651223:
		infoBank.BankName = "ALIF BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "KORTI MILLI"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 453122:
		infoBank.BankName = "ORIEN BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "VISA"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 553122:
		infoBank.BankName = "ORIEN BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "MASTER CARD"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 653122:
		infoBank.BankName = "ORIEN BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "KORTI MILLI"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 454929:
		infoBank.BankName = "SPITAMEN BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "VISA"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 554929:
		infoBank.BankName = "SPITAMEN BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "MASTER CARD"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 654929:
		infoBank.BankName = "SPITAMEN BANK"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "KORTI MILLI"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 235145:
		infoBank.BankName = "SBER"
		infoBank.BankcardName = "CLASSIC"
		infoBank.PaymentSystem = "MIR"
		infoBank.Currency = "RUB"
		return infoBank,nil
	default:
		log.Println("Such Card doesn't exists!")
		return infoBank,errors.New("Error in NumberBankCard(SP)")
	}
}


// ShowClientAccounts Посмотреть Счета клиента!
func (clDB *ClientService) ClientsBankcards(clientID int) ([]model.ResponseClientsBankcards, error) {
	listAccounts, err := clDB.ClientRep.ClientsBankcards(clientID)
	if err != nil {
		fmt.Printf("Showing a client's Accounts(SP)! Error:%e\n", err)
		return listAccounts, err
	}
	return listAccounts, nil
}


// ShowATMList Ближайшие банкоматы и филиалы!
func (clDB *ClientService) ShowATMList() ([]model.ATMsListResponse, error) {
	listATMS, err := clDB.ClientRep.ShowListOfATMS()
	if err != nil {
		fmt.Printf("Showing a list of ATMs(SP)! Error:%e\n", err)
		return listATMS, err
	}
	return listATMS, nil
}


// TransferMoneyMobi Транзакция: Отправка денег через телефонный номер Клиента!
func (clDB *ClientService) TransferMoneyMobi(sClientID int, clConf *model.TransactionMobiPhone) (*model.ClientsAccountResponse, error) {

	tStruct, err := clDB.ClientRep.InfoTransactionMobi(sClientID, clConf)   // Получение информаций об аккаунтах Sender и Receiver и запись их в структуру!
	if err != nil {
		fmt.Printf("Getting information about Sender and Receiver's accounts from DB (SP)! Error:%e", err)
		return nil, err
	}

	// Проверка баланса отправителя!
	if tStruct.Amount > tStruct.SBalance {
		fmt.Printf("Is not enough money to sent(SP)! Error:%e", err)
		return nil, err
	}

	response, err := clDB.ClientRep.TransferMoneyMobi(tStruct)
	if err != nil {
		fmt.Printf("Transaction Failed(SP)! Error:%e", err)
		return nil, err
	}
	return response,nil
}


//ListOfServices Список услуг!
func (clDB *ClientService) ListOfServices() ([]model.ServiceList, error) {
	listOfService, err := clDB.ClientRep.ListOfServices()
	if err != nil {
		fmt.Printf("Showing a list of ATMs(SP)! Error:%e\n", err)
		return listOfService, err
	}
	return listOfService, nil
}


// ListOfClientSubServices  Список подключенных услуг клиента
func (clDB *ClientService) ListOfClientSubServices (clientID int) ([]model.SubscribeServiceList,error){
	subServiceList,err := clDB.ClientRep.ListOfClientSubServices(clientID)
	if err != nil {
		fmt.Printf("Error in gettiog client's sub service list from DB (SP)! Error:%e", err)
		return subServiceList, err
	}
	return subServiceList, nil
}


// SubscribeServiceInf Подключение услуги
func (clDB *ClientService) SubscribeServiceInf (clientID int, clConf *model.SubscribeServiceRequest) (*model.ClientsAccountResponse,error){

	infService,merchantID,err := clDB.ClientRep.ServiceInformation(clConf) // Получение информации с БД про сервис!
	if err != nil {
		fmt.Printf("Error in getting information about service from DB(SP)! Error:%e", err)
		return nil,err
	}


	subService, err := clDB.ClientRep.SubscribeServiceInf(clientID,merchantID,infService) // Добавление подключенного сервиса в БД!
	if err != nil {
		fmt.Printf("Error in adding subscribe service to DB(SP)! Error:%e", err)
		return nil,err
	}

	return subService,nil
}


// CheckServiceName Проверяет ServiceName на уникальность
func (clDB *ClientService) CheckServiceName (serviceName,companyName string,clientID int) (bool,error){
	checkServiceName,err := clDB.ClientRep.CheckServiceName(clientID,serviceName,companyName)
	if err != nil {
		fmt.Printf("Checking of existing such Account name in DB (SP)! Error:%e", err)
		return false, err
	}
	return checkServiceName, nil
}


//AddAmountFromCards Пополнение кошелька с банковских карт клиента
func (clDB *ClientService) AddAmountFromCards (accountID int,clConf *model.TransferFromCardToOnlineMobi) ( *model.ClientsAccountResponse,error){

	newStruct, err := clDB.ClientRep.AddAmountFromCardsInfo(accountID, clConf.NumberOfCard) // получаем структуру данных с балансами и валютами карт!
	if err != nil {
		fmt.Println("Cannot get a struct of informations! AddAmountFromCards(SP)! Error:",err)
		return nil, err
	}

	structForAdding, err := LogicOfCurrencies(newStruct, clConf.Amount)
	if err != nil {
		fmt.Println("Error in matching currencies! Logic of Currencies (SP)! Error:",err)
		return nil, err
	}
		if structForAdding.CardBalance < clConf.Amount ||  clConf.Amount <= 0 {
			fmt.Println("Not Enough money! (SP)! Error:",err)
			return nil, err

		}
	checkServiceName,err := clDB.ClientRep.AddAmountFromCards(&structForAdding,accountID, clConf.NumberOfCard, clConf.Amount)
	if err != nil {
		log.Printf("Checking of existing such Account name in DB (SP)! Error:%e", err)
		return nil, err
	}
	return checkServiceName, nil
}


// LogicOfCurrencies Происходит конвертация и сравнение валют
func LogicOfCurrencies(clConf *model.StructTransferFromCardToOnlineMobi,amount int) (model.StructTransferFromCardToOnlineMobi,error) {
	var newStructForAdding model.StructTransferFromCardToOnlineMobi
	newStructForAdding.OnlineMobiCurrency = clConf.OnlineMobiCurrency
	newStructForAdding.CardCurrency = clConf.CardCurrency

	if clConf.OnlineMobiCurrency == "TJS" && clConf.CardCurrency == "TJS" {
		newStructForAdding.OnlineMobiBalance = clConf.OnlineMobiBalance + amount
		newStructForAdding.CardBalance = clConf.CardBalance - amount
		return newStructForAdding, nil
	}
	// Транзакция когда Currency равны!
	if clConf.OnlineMobiCurrency != clConf.CardCurrency {

		if clConf.OnlineMobiCurrency == "TJS" && clConf.CardCurrency == "RUB" {
			rAmount := amount * 100 / model.RubTjs
			newStructForAdding.OnlineMobiBalance = clConf.OnlineMobiBalance + rAmount
			newStructForAdding.CardBalance = clConf.CardBalance - amount
			return newStructForAdding, nil
		}
		if clConf.OnlineMobiCurrency == "TJS" && clConf.CardCurrency == "USD" {
			rAmount := amount * model.UsdTjs / 100
			newStructForAdding.OnlineMobiBalance = clConf.OnlineMobiBalance + rAmount
			newStructForAdding.CardBalance = clConf.CardBalance - amount
			return newStructForAdding, nil
		}
	}

	log.Println("Wrong!")
	err := errors.New("Wrong information!")
	return newStructForAdding, err
}

//HistoryOfClientsActions История действий клиента!
func (clDB *ClientService) HistoryOfClientsActions (clientID int) ([]model.HistoryList,error){
	HistoryList,err := clDB.ClientRep.HistoryOfClientsActions(clientID)
	if err != nil {
		fmt.Printf("Error in gettiog client's sub service list from DB (SP)! Error:%e", err)
		return HistoryList, err
	}
	return HistoryList, nil
}

















