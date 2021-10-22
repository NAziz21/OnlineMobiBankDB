package service

import (
	"BankVersion3/pkg/AppBank/model"
	"BankVersion3/pkg/AppBank/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)



type MerchantService struct {
	MerchantRep *repository.MerchantRepository
}

func NewMerchantService(req *repository.MerchantRepository) *MerchantService {
	return &MerchantService{
		MerchantRep: req,
	}
}


// RegisterMerchant Создаем нового клиента!
func (clDB *MerchantService) RegisterMerchant(clConf *model.Merchant) (*model.ResponseToMerchant, error) {
	login, err := clDB.MerchantRep.CheckLoginRegister(clConf.Login)
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
		log.Printf("bcrypt error when hashing(SP)! Error:%e\n", err)
		return nil, err
	}
	clConf.Password = string(bytes)
	savedUser, err := clDB.MerchantRep.RegisterMerchant(clConf)
	if err != nil {
		log.Printf("Cannot Register merchant to DB (SP)! Error:%e\n", err)
		return nil, err
	}

	return savedUser, nil
}


// LoginMerchant Проверяем полученный Логин и Пароль с БД!
func (clDB *MerchantService) LoginMerchant(login,password string) (*model.MerchantLogin, error) {
	merchants, err := clDB.MerchantRep.LoginMerchant(login)
	if err != nil {
		log.Printf("Login of Client(SP)! Error:%e\n", err)
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(merchants.Password), []byte(password)); err != nil {
		log.Printf("Comparing of Hash Password(SP)! Error:%e\n", err)
		return nil, err
	}
	return merchants, nil
}


//GetMerchantID Получаем ID Клиента через его логин!
func (clDB *MerchantService) GetMerchantID(login string) (id int,err error){
	merchantID, err := clDB.MerchantRep.GetMerchantID(login)
	if err != nil {
		return 0, err
	}

	return merchantID,nil
}


// AddBankCard Добавление карты клиента!
func (clDB *MerchantService) AddBankCard(merchantID int, clConf *model.BankCardM) (*model.ClientsAccountResponse,error) {


	checkBinNumber, err := clDB.MerchantRep.CheckBinNumber(clConf,merchantID) // Статус проверки Бин номера и сравнение имени и фамилии с БД!
	if err != nil {
		log.Printf("Cannot check a Merchants's card bin Number(SP)! Error:%e", err)
		return nil, err
	}

	if checkBinNumber == true {
		log.Println("Error! Wrong bin Number of card(RP)", err)
		return nil, err
	}

	newStruct,err := numberBankCardM(clConf)
	if err != nil {
		log.Printf("Verification of Merchant's card(SP)! Error:%e", err)
		return nil, err
	}

	checkCardName, err := clDB.MerchantRep.CheckCardName(&newStruct,merchantID)
	if err != nil {
		log.Printf("Such card has already added(SP)! Error:%e", err)
		return nil, err
	}

	if checkCardName == true {
		log.Printf("Such card has already added(SP)! Error:%e", err)
		return nil, err
	}

	_, err = clDB.MerchantRep.AddBankCard(merchantID,&newStruct)  // Добавление карты в БД!
	if err != nil {
		log.Printf("Cannot add a merchant's card to DB (SP)! Error:%e", err)
		return nil, err
	}


	responseBody := model.ClientsAccountResponse{
		Status:  true,
		Message: "Success!",
	}

	return &responseBody,err
}


// MerchantBankcards Посмотреть банковские счета клиента!
func (clDB *MerchantService) MerchantBankcards(merchantID int) ([]model.ResponseMerchantBankAccounts, error) {
	listAccounts, err := clDB.MerchantRep.MerchantBankcards(merchantID)
	if err != nil {
		log.Printf("Showing a merchant's Accounts(SP)! Error:%e\n", err)
		return listAccounts, err
	}
	return listAccounts, nil
}


//numberBankCardM Валидация
func numberBankCardM(clConf *model.BankCardM) (model.InfoBankCardM,error){
	var infoBank model.InfoBankCardM

	infoBank.CreateData = time.Now()
	infoBank.UpdateData = time.Now()
	infoBank.IBAN = clConf.IBAN
	infoBank.ORGN = clConf.ORGN
	infoBank.INN = clConf.INN
	info,_ :=strconv.ParseInt(clConf.IBAN[:9],10,0)
	switch info {
	case 570115080:
		infoBank.BankName = "MDO HUMO"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 570115081:
		infoBank.BankName = "MDO HUMO"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 570146124:
		infoBank.BankName = "ALIF BANK"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 570146129:
		infoBank.BankName = "ALIF BANK"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 570087247:
		infoBank.BankName = "ORIEN BANK"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 570087248:
		infoBank.BankName = "ORIEN BANK"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 570359514:
		infoBank.BankName = "SPITAMEN BANK"
		infoBank.Currency = "TJS"
		return infoBank,nil
	case 570359518:
		infoBank.BankName = "SPITAMEN BANK"
		infoBank.Currency = "USD"
		return infoBank,nil
	case 408007134:
		infoBank.BankName = "SBER"
		infoBank.Currency = "RUB"
		return infoBank,nil
	case 408007137:
		infoBank.BankName = "SBER"
		infoBank.Currency = "USD"
		return infoBank,nil
	default:
		log.Println("Such IBAN doesn't exists!")
		return infoBank,errors.New("Error in Bank account(SP)")
	}
}



// ListOfATMsM Ближайшие банкоматы и филиалы!
func (clDB *MerchantService) ListOfATMsM() ([]model.ATMsListResponse, error) {
	listATMS, err := clDB.MerchantRep.ListOfATMsM()
	if err != nil {
		log.Printf("Showing a list of ATMs(SP)! Error:%e\n", err)
		return listATMS, err
	}
	return listATMS, nil
}


// NewService Добавление нового сервиса!
func (clDB *MerchantService) NewService(merchantID int, mrConf *model.ServiceStructInt) (*model.ClientsAccountResponse, error) {

	var answer *model.ClientsAccountResponse


	status,err := clDB.MerchantRep.CheckServiceName( merchantID,mrConf.ServiceName)

	if err != nil {
		log.Println("Such Service Name exists(HP)! Error:", err)
		return answer, err
	}

	if status == true {
		log.Println("Such Service Name exists(HP)! Error:", err)
		return answer, err
	}

	answer, err = clDB.MerchantRep.NewService(merchantID,mrConf)
	if err != nil {
		log.Printf("Cannot add a new service to DB(SP)!Error:%e",err)
		return answer,err
	}
	return answer,nil
}




