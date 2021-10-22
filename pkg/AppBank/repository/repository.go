package repository

import (
	"BankVersion3/pkg/AppBank/model"
)

type ClientAuthorization interface {
	Register(client *model.Client) (*model.Client, error)
	CheckLoginRegister(login string) (bool, error)                                                                                             //Добавление клиента в БД!
	Login(login string) (*model.ClientLogin, error)                                                                                            //Получаем Логин и ХЭШ Пароль с БД!
	AddBankCard(clientID int, clConf *model.InfoBankCard) (bool, error)                                                                        //Добавление карты клиента!
	ClientsBankcards(clientID int) ([]model.ResponseClientsBankcards, error)                                                                   // Получение списка банковских карт клиента!
	TransferMoneyMobi(clConf *model.TransactionTransferStruct) (message *model.ClientsAccountResponse, err error)                              //Транзакция между клиентами банков по номеру телефону!
	ShowListOfATMS() ([]model.ATMsListResponse, error)                                                                                         //Показывает список банкоматов!
	ListOfServices() ([]model.ServiceList, error)                                                                                              //Показывает список всех сервисов!
	ListOfClientSubServices(clientID int) ([]model.SubscribeServiceList, error)                                                                // Показывает список подключенных сервисов клиента!
	SubscribeServiceInf(clientID, merchantID int, clConf *model.ServiceStructInt) (*model.ClientsAccountResponse, error)                       // Подключение сервиса
	AddAmountFromCards(clConf *model.StructTransferFromCardToOnlineMobi, accountID, cardID, amount int) (*model.ClientsAccountResponse, error) // Пополнение кошелька с банковской карты!
	HistoryOfClientsActions(clientID int) ([]model.HistoryList, error)                                                                         //Показывает Историю клиента!

}
type MerchantAuthorization interface {
	RegisterMerchant(merchant *model.Merchant) (*model.Merchant, error)                                         //Добавление Торговца
	LoginMerchant(login string) (*model.MerchantLogin, error)                                                   //Логинирование
	AddBankCard(merchantID int, clConf *model.InfoBankCardM) (bool, error)                                      // Добавление карты торговца!
	MerchantBankcards(clientID int) ([]model.ResponseMerchantBankAccounts, error)                               // Получение списка банковских карт  продавца!
	ListOfATMsM() ([]model.ATMsListResponse, error)                                                             // Показывает список банкоматов!
	CheckAccountNameM(accountName, networkPayment, currency string, merchantID int) (answer bool, err error)    //Проверяет не существует ли выбранный счет аккаунта в базе данных
	NewService(merchantID int, service *model.ServiceStruct) (message *model.ClientsAccountResponse, err error) //Добавление нового сервиса
}

type Repository struct {
	ClientAuthorization
	MerchantAuthorization
}
