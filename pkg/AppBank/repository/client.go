package repository

import (
	"BankVersion3/pkg/AppBank/model"
	database "BankVersion3/pkg/DataBase"
	"context"
	"fmt"
	"log"
	"strconv"
)

const Type = `Between Clients`
const Subscribe = `Client / Merchant`
const Account = `ONLINE MOBI`
const Service = "Subscribe"
const TypeA = "From Card"

type ClientRepository struct {
	DB *database.DB
}

func NewClientRep(db *database.DB) *ClientRepository {
	return &ClientRepository{
		db,
	}
}


// Register Создание нового клиента
func (dbClient *ClientRepository) Register(client *model.Client) (*model.ResponseToClient,error) {
	var clientID int
	tx, err := dbClient.DB.Conn.Begin(context.Background())
	if err != nil {
		fmt.Printf("transaction failed %e", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				log.Println("Transaction Failed! Error:",err)
			}
		}
		err := tx.Commit(context.Background())
		if err != nil {
			return
		}
		if err != nil {
			log.Println("Transaction Failed! Commit has failed! Error:",err)
		}
	}()
	_, err = tx.Exec(context.Background(), `insert into client(name, surname, login, password, phone,create_data,update_data)
    values ($1,$2,$3,$4,$5,$6,$7) returning id`, client.Name, client.Surname, client.Login, client.Password, client.Phone, database.DataTime, database.DataTime)
	if err != nil {
		fmt.Print("Cannot insert Client's data to DB(RP)!Error:", err)
		return nil,err
	}

	err = tx.QueryRow(context.Background(),`select id from client where login = $1`, client.Login).Scan(&clientID)
	if err != nil {
		fmt.Print("Cannot get Client's id from DB(RP)!Error:", err)
		return nil,err
	}

	_, err = tx.Exec(context.Background(), `insert into account_online(client_id, create_data, update_data) 
      values ($1,$2,$3)`, clientID, database.DataTime, database.DataTime)
	if err != nil {
		fmt.Println("Cannot add a new client(RP)! Error:", err)
		return nil, err
	}

	responseBody := model.ResponseToClient{
		Status:  true,
		Message: "Success!",
		Name:    client.Name,
		Surname: client.Surname,
	}

	return &responseBody, nil
}

// CheckLoginRegister Проверка логина на дубликат
func (dbClient *ClientRepository)CheckLoginRegister(login string) (bool, error) {
	var answer bool
	err := dbClient.DB.Conn.QueryRow(context.Background(), "select exists (select login from client where login=$1)", login).Scan(&answer)
	if err != nil {
		log.Printf("Can't get client's login(RP)! %e", err)
		return answer, err
	}
	if answer == true {
		return true,err
	}
	return false,nil
}


// Login Получаем Логин и ХЭШ Пароль с БД!
func (dbClient *ClientRepository) Login(login string) (*model.ClientLogin, error) {
	var client model.ClientLogin

	err := dbClient.DB.Conn.QueryRow(context.Background(), "select login, password from client where login=$1", login).Scan(&client.Login, &client.Password)
	if err != nil {
		log.Println("Cannot get client's login and password from DB(RP)!Error:", err)
		return nil, err
	}
	return &client, nil
}


// GetClientID Получаем ID клиента!
func (dbClient *ClientRepository) GetClientID(login string) (id int, err error) {
	err = dbClient.DB.Conn.QueryRow(context.Background(), "select id from client where login=$1", login).Scan(&id)
	if err != nil {
		log.Println("GetClientID! Error!", err)
		return 0, err
	}
	return id, nil
}


// AddBankCard Добавление карты клиента!
func (dbClient *ClientRepository) AddBankCard(clientID int, clConf *model.InfoBankCard) (bool, error) {
	_, err := dbClient.DB.Conn.Exec(context.Background(),
		`insert into bank_account_client(client_id,bank_name,bankcard_name,bin_number,cvc,payment_system,currency, create_data, update_data) 
      values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, clientID, clConf.BankName, clConf.BankcardName, clConf.BinNumber,clConf.CVC,clConf.PaymentSystem,
		clConf.Currency, database.DataTime, database.DataTime)
	if err != nil {
		fmt.Println("Cannot add a new client's bank card to DB(RP)! Error:", err)
		return false, err
	}
	return true, nil
}


// CheckBinNumber Проверка BIN Number, name and surname!
func (dbClient *ClientRepository) CheckBinNumber(clConf *model.BankCard,clientID int) (bool, error) {
	var answerBinNumber bool
	err := dbClient.DB.Conn.QueryRow(context.Background(),`select exists(select bin_number from bank_account_client where bin_number =$1)`,clConf.NumberCard).Scan(&answerBinNumber)
	if err != nil {
		fmt.Println("Cannot check a bin number for unique(RP)! Error:", err)
		return true, err
	}

	if answerBinNumber == true {
		fmt.Println("Error! Wrong bin Number of card(RP)", err)
		return true,err
	}
	return false, nil
}


// CheckCardName Проверка на существование такого аккаунта в БД!
func (dbClient *ClientRepository) CheckCardName(clConf *model.InfoBankCard,clientID int) (bool, error) {
	var accountName bool
	err := dbClient.DB.Conn.QueryRow(context.Background(),`select exists(select * from bank_account_client where bank_name =$1 and currency = $2 and payment_system = $3 and client_id =$4)`,clConf.BankName,clConf.Currency,clConf.PaymentSystem,clientID).Scan(&accountName)
	if err != nil {
		fmt.Println("Cannot check a bin number for unique(RP)! Error:", err)
		return true, err
	}

	if accountName == true {
		fmt.Println("Error! Wrong bin Number of card(RP)", err)
		return true,err
	}
	return false, nil
}


// TransferMoneyMobi Транзакция между клиентами банков по номеру телефону!
func (dbClient *ClientRepository) TransferMoneyMobi(clConf *model.TransactionTransferStruct) (message *model.ClientsAccountResponse, err error) {

	tx, err := dbClient.DB.Conn.Begin(context.Background())
	if err != nil {
		fmt.Printf("transaction failed %e", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				fmt.Println(err)
			}
		}
		err := tx.Commit(context.Background())
		if err != nil {
			return
		}
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = tx.Exec(context.Background(), `update account_online set balance = $1 where client_id = $2`, clConf.SBalance- clConf.Amount, clConf.SClientID)
	if err != nil {
		fmt.Println("Cannot update Sender's amount(RP)! Error:", err)
		return nil, err
	}

	_, err = tx.Exec(context.Background(), `update account_online set balance = $1 where client_id = $2`,clConf.RBalance + clConf.Amount, clConf.RClientID)
	if err != nil {
		fmt.Println("Cannot update Receiver's amount(RP)! Error:", err)
		return nil, err
	}

	_, err = tx.Exec(context.Background(), `insert into transactionsclients (sender_account,receiver_account,transaction_amount,create_data) values ($1,$2,$3,$4)`, clConf.SClientID, clConf.RClientID, clConf.Amount, database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	}
	_, err = tx.Exec(context.Background(), `insert into debit_credit (client_id,transaction_type,account_name,credit,total_amount_before,total_amount_after,data_transaction) values ($1,$2,$3,$4,$5,$6,$7)`, clConf.SClientID, Type, Account,clConf.Amount,clConf.SBalance,clConf.SBalance-clConf.Amount,database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	}
	_, err = tx.Exec(context.Background(), `insert into debit_credit (client_id,transaction_type,account_name,debit,total_amount_before,total_amount_after,data_transaction) values ($1,$2,$3,$4,$5,$6,$7)`, clConf.RClientID, Type, Account,clConf.Amount,clConf.RBalance,clConf.RBalance+clConf.Amount,database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	}

	response := model.ClientsAccountResponse{
		Status:  true,
		Message: "Success!",
	}

	return &response, nil
} // Transfer money to another client!


// InfoTransactionMobi получаем данные аккаунтов Sender и receiver
func (dbClient *ClientRepository) InfoTransactionMobi(sClientID int, clConf *model.TransactionMobiPhone) (*model.TransactionTransferStruct, error) {
	var sBalance, rBalance, rClientID int

	err := dbClient.DB.Conn.QueryRow(context.Background(), `select balance from account_online where client_id =$1`, sClientID).Scan(&sBalance)
	if err != nil {
		fmt.Println("Error GetAccountSender(RP)", err)
		return nil, err
	}

	err = dbClient.DB.Conn.QueryRow(context.Background(), `select id from client where phone =$1`, clConf.PhoneNumber).Scan(&rClientID)
	if err != nil {
		fmt.Println("Cannot get RClientID(RP)! Error:", err)
		return nil, err
	}

	err = dbClient.DB.Conn.QueryRow(context.Background(), `select balance from account_online where client_id =$1 `, rClientID).Scan(&rBalance)
	if err != nil {
		fmt.Println("Error GetAccountSender(RP)", err)
		return nil, err
	}

	amount,_ := strconv.ParseInt(clConf.Amount,10,0)

	tStruct := model.TransactionTransferStruct{
		SClientID: sClientID,
		RClientID: rClientID,
		SBalance:  sBalance,
		RBalance:  rBalance,
		Amount: int(amount),
	}
	return &tStruct, nil
}


// ClientsBankcards Получение списка  банковских карт клиента!
func (dbClient *ClientRepository) ClientsBankcards(clientID int) ([]model.ResponseClientsBankcards, error) {
	var clientList []model.ResponseClientsBankcards
	rows, err := dbClient.DB.Conn.Query(context.Background(), "select id, bank_name, bankcard_name, balance,payment_system,currency from bank_account_client where client_ID=$1", clientID)
	if err != nil {
		fmt.Printf("can't get clients %e", err)
		return clientList, err
	}
	for rows.Next() {
		Client := model.ResponseClientsBankcards{}
		err := rows.Scan(&Client.ID,&Client.BankName, &Client.BankcardName,&Client.Balance,&Client.PaymentSystem, &Client.Currency)
		if err != nil {
			fmt.Printf("can't scan %e", err)
		}
		clientList = append(clientList, Client)
	}
	if rows.Err() != nil {
		fmt.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return clientList, err
}


// ShowListOfATMS Показывает список банкоматов!
func (dbClient *ClientRepository) ShowListOfATMS() ([]model.ATMsListResponse, error) {
	var ATMList []model.ATMsListResponse
	rows, err := dbClient.DB.Conn.Query(context.Background(), "select address from atms")
	if err != nil {
		fmt.Printf("can't get clients %e", err)
		return ATMList, err
	}
	for rows.Next() {
		NewList := model.ATMsListResponse{}
		err := rows.Scan(&NewList.Address)
		if err != nil {
			fmt.Printf("can't scan %e", err)
		}
		ATMList = append(ATMList, NewList)
	}
	if rows.Err() != nil {
		fmt.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return ATMList, err
}


// ListOfServices Показывает список всех сервисов!
func (dbClient *ClientRepository) ListOfServices() ([]model.ServiceList, error) {
	var ServiceList []model.ServiceList
	rows, err := dbClient.DB.Conn.Query(context.Background(), "select service_name,service_category,company_service_name,service_price,currency,type from services")
	if err != nil {
		fmt.Printf("can't get clients %e", err)
		return ServiceList, err
	}
	for rows.Next() {
		NewList := model.ServiceList{}
		err := rows.Scan(&NewList.ServiceName,&NewList.ServiceCategory,&NewList.CompanyServiceName,&NewList.ServicePrice,&NewList.Currency,&NewList.Type)
		if err != nil {
			fmt.Printf("can't scan %e", err)
		}
		//Astring := strconv.Itoa(NewList.ServicePrice)
		ServiceList = append(ServiceList, NewList)
	}
	if rows.Err() != nil {
		fmt.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return ServiceList, err
}


// ListOfServices Показывает список подключенных сервисов клиента!
func (dbClient *ClientRepository) ListOfClientSubServices(clientID int) ([]model.SubscribeServiceList, error) {
	var subServiceList []model.SubscribeServiceList
	rows, err := dbClient.DB.Conn.Query(context.Background(), "select id,service_name,service_category,company_service_name,service_price,currency,type,service_create_data from services_client where client_id =$1",clientID)
	if err != nil {
		fmt.Printf("can't get clients %e", err)
		return subServiceList, err
	}
	for rows.Next() {
		NewList := model.SubscribeServiceList{}
		err := rows.Scan(&NewList.ID,&NewList.ServiceName,&NewList.ServiceCategory,&NewList.CompanyServiceName,&NewList.ServicePrice,&NewList.Currency,&NewList.Type,&NewList.CreateData)
		if err != nil {
			fmt.Printf("can't scan %e", err)
		}
		subServiceList = append(subServiceList, NewList)
	}
	if rows.Err() != nil {
		fmt.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return subServiceList, err
}


// SubscribeServiceInf Подключение сервиса
func (dbClient *ClientRepository) SubscribeServiceInf(clientID,merchantID int,clConf *model.ServiceStructInt,) (*model.ClientsAccountResponse,error){
	var message model.ClientsAccountResponse
	var cBalance, mBalance int


	tx, err := dbClient.DB.Conn.Begin(context.Background())
	if err != nil {
		fmt.Printf("transaction failed %e", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				fmt.Println(err)
			}
		}
		err := tx.Commit(context.Background())
		if err != nil {
			return
		}
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = dbClient.DB.Conn.QueryRow(context.Background(), `select balance from account_online where client_id =$1`, clientID).Scan(&cBalance)
	if err != nil {
		fmt.Println("Error cannot get client's balance(RP)", err)
		return nil, err
	} // Происходит взятие баланса Покупателя!


	err = dbClient.DB.Conn.QueryRow(context.Background(), `select balance from account_online where merchant_id =$1`, merchantID).Scan(&mBalance)
	if err != nil {
		fmt.Println("Error cannot get merchant's balance (RP)", err)
		return nil, err
	} // Происходит взятие баланса Продавца!


	if cBalance < clConf.ServicePrice {
		fmt.Println("Not Enough money to purchase a subscribe!(RP)", err)
		return nil, err
	}


	_, err = tx.Exec(context.Background(), `update account_online set balance = $1 where client_id = $2`, cBalance- clConf.ServicePrice, clientID)
	if err != nil {
		fmt.Println("Cannot update Sender's amount(RP)! Error:", err)
		return nil, err
	} // Update Client's Account


	_, err = tx.Exec(context.Background(), `update account_online set balance = $1 where merchant_id = $2`,mBalance + clConf.ServicePrice, merchantID)
	if err != nil {
		fmt.Println("Cannot update Receiver's amount(RP)! Error:", err)
		return nil, err
	} // Update Merchant's Account


	_, err = tx.Exec(context.Background(), `insert into transactionsmerchants (sender_account,receiver_account,transaction_amount,create_data) values ($1,$2,$3,$4)`,
		clientID, merchantID, clConf.ServicePrice, database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	} // Добавление в таблицу историй транзакций


	_, err = tx.Exec(context.Background(), `insert into debit_credit (client_id,transaction_type,account_name,credit,total_amount_before,total_amount_after,data_transaction) values ($1,$2,$3,$4,$5,$6,$7)`,
		clientID, Type, Account,clConf.ServicePrice,cBalance,cBalance-clConf.ServicePrice,database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	} // Добавление в историю Дебит/Кредита "История" Клиента!


	_, err = tx.Exec(context.Background(), `insert into debit_credit (client_id,transaction_type,account_name,debit,total_amount_before,total_amount_after,data_transaction) values ($1,$2,$3,$4,$5,$6,$7)`,
		merchantID, Subscribe, Service,clConf.ServicePrice,mBalance,mBalance+clConf.ServicePrice,database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	} // Добавление в историю Дебит/Кредита "История" Торговца!


	_, err = tx.Exec(context.Background(), `insert into services_client(client_id,service_name,service_category,
company_service_name,service_price,currency,type, service_create_data, service_update_data) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
clientID, clConf.ServiceName, clConf.ServiceCategory,clConf.CompanyServiceName, clConf.ServicePrice, clConf.Currency,clConf.Type,database.DataTime, database.DataTime)
	if err != nil {
		fmt.Println("Cannot add a new client(RP)! Error:", err)
		return &message, err
	} // Добавление сервиса в таблицу сервисов Клиентов!

	message = model.ClientsAccountResponse{
		Status:  true,
		Message: "Success!",
	}
	return &message, nil
}


// ServiceInformation Получение информации о сервисе!
func (dbClient *ClientRepository) ServiceInformation(clConf *model.SubscribeServiceRequest) (*model.ServiceStructInt,int,error){
	var infService model.ServiceStructInt
	var merchantID int
	err := dbClient.DB.Conn.QueryRow(context.Background(), `select service_merchant_id,service_name,service_category,company_service_name,service_price,currency,type from services where
 service_name =$1 and company_service_name =$2`,clConf.ServiceName,clConf.CompanyName).Scan(&merchantID,&infService.ServiceName,&infService.ServiceCategory,&infService.CompanyServiceName,&infService.ServicePrice,&infService.Currency,&infService.Type)
	if err != nil {
		fmt.Println("Error GetAccountSender(RP)", err)
		return &infService,merchantID,err
	}
	return &infService, merchantID,nil
}


//CheckServiceName проверка Сервиса на уникальность!
func (dbClient *ClientRepository) CheckServiceName(clientID int,serviceName, companyName string) (bool, error) {
	var existsInDB, subscribed bool
	err := dbClient.DB.Conn.QueryRow(context.Background(), `select exists (select service_name,company_service_name from services where
service_name =$1 and company_service_name =$2)`, serviceName, companyName).Scan(&existsInDB)
	if err != nil {
		fmt.Println("Error GetAccountSender(RP)", err)
		return false,err
	}
	err = dbClient.DB.Conn.QueryRow(context.Background(), `select exists (select client_id,service_name,company_service_name from services_client where
service_name =$1 and company_service_name =$2 and client_id =$3)`,serviceName, companyName,clientID).Scan(&subscribed)
	if err != nil {
		fmt.Println("Error GetAccountSender(RP)", err)
		return false,err
	}

	if existsInDB == true {
		if subscribed == false {
			return true,nil
		}
	}
	return false ,nil
}


// AddAmountFromCards Пополнение кошелька с банковской карты!
func (dbClient *ClientRepository) AddAmountFromCards(clConf *model.StructTransferFromCardToOnlineMobi,accountID,cardID,amount int) (*model.ClientsAccountResponse,error){
	var message model.ClientsAccountResponse
	tx, err := dbClient.DB.Conn.Begin(context.Background())
	if err != nil {
		fmt.Printf("transaction failed %e", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				fmt.Println(err)
			}
		}
		err := tx.Commit(context.Background())
		if err != nil {
			return
		}
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = tx.Exec(context.Background(), `update account_online set balance = $1 where client_id = $2`, clConf.OnlineMobiBalance,accountID )
	if err != nil {
		fmt.Println("Cannot update Sender's amount(RP)! Error:", err)
		return nil, err
	} // Update Client's Account


	_, err = tx.Exec(context.Background(), `update bank_account_client set balance = $1 where id = $2`,clConf.CardBalance, cardID)
	if err != nil {
		fmt.Println("Cannot update Receiver's amount(RP)! Error:", err)
		return nil, err
	} // Update Merchant's Account


	_, err = tx.Exec(context.Background(), `insert into transactionsclients (sender_account,receiver_account,transaction_amount,create_data) values ($1,$2,$3,$4)`,
		accountID, accountID, amount, database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	} // Добавление в таблицу историй транзакций


	_, err = tx.Exec(context.Background(), `insert into debit_credit (client_id,transaction_type,account_name,credit,total_amount_before,total_amount_after,data_transaction) values ($1,$2,$3,$4,$5,$6,$7)`,
		accountID, TypeA, Account,amount,clConf.OnlineMobiBalance - amount,clConf.OnlineMobiBalance,database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	} // Добавление в историю Дебит/Кредита "История" Клиента!


	_, err = tx.Exec(context.Background(), `insert into debit_credit (client_id,transaction_type,account_name,debit,total_amount_before,total_amount_after,data_transaction) values ($1,$2,$3,$4,$5,$6,$7)`,
		accountID, TypeA, Account,amount,clConf.CardBalance+amount,clConf.CardBalance,database.DataTime)
	if err != nil {
		fmt.Println("Insert data to Transaction table failed", err)
	} // Добавление в историю Дебит/Кредита "История" Торговца!


	message = model.ClientsAccountResponse{
		Status:  true,
		Message: "Success!",
	}
	return &message, nil
}


// AddAmountFromCardsGettingBalance
func (dbClient *ClientRepository) AddAmountFromCardsInfo(accountID,cardID int) (*model.StructTransferFromCardToOnlineMobi,error) {

	var cardBalance, onlineAccountBalance int
	var cardCurrency, accountCurrency string

	err := dbClient.DB.Conn.QueryRow(context.Background(), `select balance from account_online where client_id =$1`, accountID).Scan(&onlineAccountBalance)
	if err != nil {
		fmt.Println("Error cannot get client's balance from onlineMobi(RP)", err)
		return nil, err
	} // Происходит взятие баланса Покупателя!


	err = dbClient.DB.Conn.QueryRow(context.Background(), `select balance from bank_account_client where id =$1`, cardID).Scan(&cardBalance)
	if err != nil {
		fmt.Println("Error cannot get client's balance from Card table DB! (RP)", err)
		return nil, err
	} // Происходит взятие баланса Продавца!

	err = dbClient.DB.Conn.QueryRow(context.Background(), `select currency from account_online where client_id =$1`, accountID).Scan(&accountCurrency)
	if err != nil {
		fmt.Println("Error cannot get OnlineMobi's currency(RP)", err)
		return nil, err
	} // Происходит взятие валюты электронного кошелька

	err = dbClient.DB.Conn.QueryRow(context.Background(), `select currency from bank_account_client where id =$1 and client_id =$2`, cardID,accountID).Scan(&cardCurrency)
	if err != nil {
		fmt.Println("Error cannot get Card's currency(RP)", err)
		return nil, err
	} // Происходит взятие валюты электронного кошелька

	getStruct := model.StructTransferFromCardToOnlineMobi{
		OnlineMobiCurrency: accountCurrency,
		CardCurrency:       cardCurrency,
		OnlineMobiBalance:  onlineAccountBalance,
		CardBalance:        cardBalance,
	}

	return &getStruct, nil
}



// ListOfServices Показывает список подключенных сервисов клиента!
func (dbClient *ClientRepository) HistoryOfClientsActions(clientID int) ([]model.HistoryList, error) {
	var HistoryOfActionsList []model.HistoryList
	rows, err := dbClient.DB.Conn.Query(context.Background(), "select account_name, transaction_type,debit,credit, data_transaction from debit_credit where client_id = $1",clientID)
	if err != nil {
		fmt.Printf("can't get clients %e", err)
		return HistoryOfActionsList, err
	}
	for rows.Next() {
		NewList := model.HistoryList{}
		err := rows.Scan(&NewList.AccountName,
			&NewList.TransactionType,
			&NewList.Debit,
			&NewList.Credit,
			&NewList.DataOfTransaction)
		if err != nil {
			fmt.Printf("can't scan HistroyOfActions %e", err)
		}
		HistoryOfActionsList = append(HistoryOfActionsList, NewList)
	}
	if rows.Err() != nil {
		fmt.Printf("rows err %e", err)
		return nil, rows.Err()
	}
	return HistoryOfActionsList, err
}



