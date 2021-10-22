package repository

import (
	"BankVersion3/pkg/AppBank/model"
	database "BankVersion3/pkg/DataBase"
	"context"
	"fmt"
	"log"
)

type MerchantRepository struct {
	DB *database.DB
}

func NewMerchantRep(db *database.DB) *MerchantRepository {
	return &MerchantRepository{
		db,
	}
}

// AddMerchant Добавление Торговца
func (dbMerchant *MerchantRepository) RegisterMerchant(merchant *model.Merchant) (*model.ResponseToMerchant, error) {
	var merchantID int

	tx, err := dbMerchant.DB.Conn.Begin(context.Background())
	if err != nil {
		fmt.Printf("transaction failed %e", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				log.Println("Transaction Failed! Error:", err)
			}
		}
		err = tx.Commit(context.Background())
		if err != nil {
			return
		}
		if err != nil {
			log.Println("Transaction Failed! Commit has failed! Error:", err)
		}
	}()
	_, err = dbMerchant.DB.Conn.Exec(context.Background(), `insert into merchants(company_name,company_address,company_description,owner_name,owner_surname,login,password,create_data,update_data) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		merchant.CompanyName, merchant.CompanyAddress, merchant.CompanyDescription, merchant.OwnerName, merchant.OwnerSurname, merchant.Login, merchant.Password, database.DataTime, database.DataTime)
	if err != nil {
		log.Println("Error in Adding Merchant is:", err)
		return nil, err
	}
	err = tx.QueryRow(context.Background(), `select id from merchants where login = $1`, merchant.Login).Scan(&merchantID)
	if err != nil {
		log.Println("Cannot get Client's id from DB(RP)!Error:", err)
		return nil, err
	}

	_, err = tx.Exec(context.Background(), `insert into account_online(merchant_id, create_data, update_data) 
      values ($1,$2,$3)`, merchantID, database.DataTime, database.DataTime)
	if err != nil {
		log.Println("Cannot add a new client(RP)! Error:", err)
		return nil, err
	}

	responseBody := model.ResponseToMerchant{
		Status:         true,
		Message:        "Success!",
		CompanyName:    merchant.CompanyName,
		CompanyAddress: merchant.CompanyDescription,
	}

	return &responseBody, nil
}

// CheckLoginRegister  Проверка логина на дубликат
func (dbMerchant *MerchantRepository) CheckLoginRegister(login string) (bool, error) {
	var answer bool
	err := dbMerchant.DB.Conn.QueryRow(context.Background(), "select exists (select login from merchants where login=$1)", login).Scan(&answer)
	if err != nil {
		log.Printf("Can't get merchant login(RP)! %e", err)
		return answer, err
	}
	if answer == true {
		return true, err
	}
	return false, nil
}

// LoginMerchant Логинирование
func (dbMerchant *MerchantRepository) LoginMerchant(login string) (*model.MerchantLogin, error) {
	var merchant model.MerchantLogin
	err := dbMerchant.DB.Conn.QueryRow(context.Background(), "select login, password from merchants where login=$1", login).Scan(&merchant.Login, &merchant.Password)
	if err != nil {
		log.Printf("can't get merchant login and password %e", err)
		return nil, err
	}
	return &merchant, nil
}

// GetMerchantID Берем ID торговца!
func (dbMerchant *MerchantRepository) GetMerchantID(login string) (id int, err error) {
	err = dbMerchant.DB.Conn.QueryRow(context.Background(), "select id from merchants where login=$1", login).Scan(&id)
	if err != nil {
		log.Printf("Cannot get A Merchant ID(RP)! Error:%e", err)
		return 0, err
	}
	return id, nil
}

// AddBankCard Добавление карты торговца!
func (dbMerchant *MerchantRepository) AddBankCard(merchantID int, clConf *model.InfoBankCardM) (bool, error) {
	_, err := dbMerchant.DB.Conn.Exec(context.Background(),
		`insert into bank_account_merchant(merchant_id,bank_name,iban,inn,orgn,currency, create_data, update_data) 
      values ($1,$2,$3,$4,$5,$6,$7,$8)`, merchantID, clConf.BankName, clConf.IBAN, clConf.INN, clConf.ORGN,
		clConf.Currency, database.DataTime, database.DataTime)
	if err != nil {
		log.Println("Cannot add a new merchant's bank account to DB(RP)! Error:", err)
		return false, err
	}
	return true, nil
}

// CheckBinNumber Проверка BIN Number, name and surname!
func (dbMerchant *MerchantRepository) CheckBinNumber(clConf *model.BankCardM, merchantID int) (bool, error) {
	var answer bool
	err := dbMerchant.DB.Conn.QueryRow(context.Background(), `select exists(select iban from bank_account_merchant where iban =$1 and inn =$2 and orgn =$3)`, clConf.IBAN, clConf.INN, clConf.ORGN).Scan(&answer)
	if err != nil {
		log.Println("Cannot check a bin number for unique(RP)! Error:", err)
		return true, err
	}

	if answer == true {
		log.Println("Error! Wrong information of card(RP)", err)
		return true, err
	}
	return false, nil
}

// CheckCardName Проверка нового аккаунта на уникальность перед добавлением, name and surname!
func (dbMerchant *MerchantRepository) CheckCardName(clConf *model.InfoBankCardM, merchantID int) (bool, error) {
	var accountName bool
	err := dbMerchant.DB.Conn.QueryRow(context.Background(), `select exists(select * from bank_account_merchant where bank_name =$1 and currency = $2 and merchant_id =$3)`, clConf.BankName, clConf.Currency, merchantID).Scan(&accountName)
	if err != nil {
		log.Println("Cannot check a bin number for unique(RP)! Error:", err)
		return true, err
	}

	if accountName == true {
		log.Println("Error! Wrong bin Number of card(RP)", err)
		return true, err
	}
	return false, nil
}

// MerchantBankcards Получение списка банковских карт  продавца!
func (dbMerchant *MerchantRepository) MerchantBankcards(clientID int) ([]model.ResponseMerchantBankAccounts, error) {
	var merchantList []model.ResponseMerchantBankAccounts
	rows, err := dbMerchant.DB.Conn.Query(context.Background(), "select bank_name,balance,currency from bank_account_merchant where merchant_ID=$1", clientID)
	if err != nil {
		log.Printf("can't get merchants account %e", err)
		return merchantList, err
	}
	for rows.Next() {
		Merchant := model.ResponseMerchantBankAccounts{}
		err := rows.Scan(&Merchant.BankName, &Merchant.Balance, &Merchant.Currency)
		if err != nil {
			log.Printf("can't scan %e", err)
		}
		merchantList = append(merchantList, Merchant)
	}
	if rows.Err() != nil {
		log.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return merchantList, err
}

// ListOfATMsM Показывает список банкоматов!
func (dbMerchant *MerchantRepository) ListOfATMsM() ([]model.ATMsListResponse, error) {
	var ATMList []model.ATMsListResponse
	rows, err := dbMerchant.DB.Conn.Query(context.Background(), "select address from atms")
	if err != nil {
		log.Printf("can't get clients %e", err)
		return ATMList, err
	}
	for rows.Next() {
		NewList := model.ATMsListResponse{}
		err := rows.Scan(&NewList.Address)
		if err != nil {
			log.Printf("can't scan %e", err)
		}
		ATMList = append(ATMList, NewList)
	}
	if rows.Err() != nil {
		log.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return ATMList, err
}

// NewService Добавление нового сервиса
func (dbMerchant *MerchantRepository) NewService(merchantID int, service *model.ServiceStructInt) (message *model.ClientsAccountResponse, err error) {
	_, err = dbMerchant.DB.Conn.Exec(context.Background(), `insert into services(service_name, service_category, company_service_name, service_merchant_id, service_price, currency,type,service_create_data, service_update_data) 
    values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, service.ServiceName, service.ServiceCategory, service.CompanyServiceName, merchantID, service.ServicePrice, service.Currency, service.Type, database.DataTime, database.DataTime)
	if err != nil {
		log.Printf("Error in inserting service %e", err)
		return message, err
	}

	response := model.ClientsAccountResponse{
		Status:  true,
		Message: "Success!",
	}

	return &response, nil
}

// CheckServiceName Проверяет, существует ли такое имя в базе данных!
func (dbMerchant *MerchantRepository) CheckServiceName(serviceMerchantID int, serviceName string) (answer bool, err error) {

	err = dbMerchant.DB.Conn.QueryRow(context.Background(), `select exists (select service_name, company_service_name,services from  services where service_name = $1 and service_merchant_id = $2)`, serviceName, serviceMerchantID).Scan(&answer)
	if err != nil {
		log.Printf("Error: \"cannot check the merchantID in DT %e\"", err)
		return true, err
	}
	if answer == true {
		return true, err
	}
	return false, nil
}
