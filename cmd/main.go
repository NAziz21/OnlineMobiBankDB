package main

import (
  "BankVersion3/pkg/AppBank/handler"
  "BankVersion3/pkg/AppBank/repository"
  "BankVersion3/pkg/AppBank/service"
  database "BankVersion3/pkg/DataBase"
  "BankVersion3/pkg/middleware"
  "github.com/gorilla/mux"
  _ "gopkg.in/natefinch/lumberjack.v2"
  "io"
  "log"
  "net/http"
  "os"
)

func main() {

  logFile, err := os.OpenFile("log.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
  if err != nil {
    panic(err)
  }
  mw := io.MultiWriter(os.Stdout, logFile)
  log.SetOutput(mw)


  db := database.NewDB("postgres://Aziz:root@localhost:5432/bankdatabase")
  db.ConnectionDB()

  ClientRepository := repository.NewClientRep(db)
  MerchantRepository := repository.NewMerchantRep(db)

  ClientService := service.NewClientService(ClientRepository)
  MerchantService := service.NewMerchantService(MerchantRepository)

  router := mux.NewRouter()
  ClientHandler :=handler.ClientRestAPI(ClientService)
  MerchantHandler :=handler.MerchantRestApi(MerchantService)

  router.HandleFunc("/client/register",ClientHandler.Register).Methods("POST")
  router.HandleFunc("/client/login",ClientHandler.Login).Methods("POST")
  router.HandleFunc("/merchant/register",MerchantHandler.RegisterMerchant).Methods("POST")
  router.HandleFunc("/merchant/login",MerchantHandler.LoginMerchant).Methods("POST")


  s := router.PathPrefix("/auth/").Subrouter()
  s.Use(middleware.Auth)
  s.HandleFunc("/client/clientsbankcards",ClientHandler.ClientsBankcards).Methods("POST")
  s.HandleFunc("/client/listOfATMs",ClientHandler.ListOfATMs)
  s.HandleFunc("/client/transferMoney",ClientHandler.TransferMoneyMobi).Methods("POST")
  s.HandleFunc("/client/ListOfServices",ClientHandler.ListOfServices)
  s.HandleFunc("/client/subService",ClientHandler.SubscribeService).Methods("POST")
  s.HandleFunc("/client/addBankcard",ClientHandler.AddBankCard).Methods("POST")
  s.HandleFunc("/client/amountFromCard",ClientHandler.AddAmountFromCards).Methods("POST")
  s.HandleFunc("/client/listOfClientsSubservices",ClientHandler.ListOfClientSubServices).Methods("POST")
  s.HandleFunc("/client/listOfHistoryOfClientActions",ClientHandler.HistoryOfClientsActions).Methods("POST")




  s.HandleFunc("/merchant/merchantsbankcards",MerchantHandler.MerchantBankcards).Methods("POST")
  s.HandleFunc("/merchant/listOfATMsM",MerchantHandler.ListOfATMsM).Methods("POST")
  s.HandleFunc("/merchant/addNewService",MerchantHandler.NewService).Methods("POST")
  s.HandleFunc("/merchant/addBankcardM",MerchantHandler.AddBankCard).Methods("POST")
  s.Use()

  err = http.ListenAndServe(":8181", router)
  if err != nil {
    log.Fatalf("can't listen and serve 8181 %e",err)
  }

}
