package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)



type DB struct {
	Config     string
	Conn *pgx.Conn
}

func NewDB(conf string) *DB {
	return &DB{
		Config: conf,
	}
}

func (dataBase *DB) ConnectionDB() {
	fmt.Println("Connection is started!")
	conn, err := pgx.Connect(context.Background(), dataBase.Config)
	if err != nil {
		log.Fatalf("can't connect to db %e", err)
	}
	dataBase.Conn = conn
	dataBase.DBInitialization()

}


/*var Connect = ConnectionDB()

func ConnectionDB() *pgx.Conn{
	fmt.Println("Connection to DB started!")

	//Unmarshal / marshal --> Parse from settings.json to model SettingDB

	bytes, err := ioutil.ReadFile("settings/settings.json")
	if err != nil {
		log.Fatalf("Error can't read from file %e", err)
	}
	var DataObject settings.DB
	err = json.Unmarshal(bytes, &DataObject)
	if err != nil {
		log.Fatalf("Error can't parse bytes %e", err)
	}
	// urlExample := "postgres://username:password@localhost:5432/database_name
	urlDatabase := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", DataObject.Username, DataObject.Password, DataObject.DatabaseName)
	urlDatabase := fmt.Sprintf("postgres://Aziz:root@localhost:5432/bankdatabase", DataObject.Username, DataObject.Password, DataObject.DatabaseName)
	connect, err := pgx.Connect(context.Background(), urlDatabase)
	if err != nil {
		log.Fatalf("can't connect to DB %e", err)
	}
	return connect
}
*/
