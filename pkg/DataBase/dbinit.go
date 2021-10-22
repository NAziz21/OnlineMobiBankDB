package database

import (
	"context"
	"log"
	"time"
)


const kindOfAccountsCl = `
                "Type of Account"
1. Humo Account
2. Sberbank Account
3. Alif Account
4. Raifassen Account
5. Orienbank Account`

const TypeOfCurrency = `
		"Currency"
1. TJS
2. RUB
3. USD`

const TypeOfNetworkPayment =` 
		"Network Payment"
1. Korti Milli
2. Visa 
3. Mastercard 
4. Mir`


var DataTime = time.Now()

func (dataBase *DB) DBInitialization() {

	DDLs := []string{ClientsDDL, AccountsDDL,MobiAccountDDL,BankCardDDl,BankAccountDDL, MerchantsDDL, ServicesDDL,ATMsDDL,ManagerDDL,TransactionsMerchantsDDL,TransactionsClientsDDL,DebitCreditDDL,ClientServicesDDL}
	for _, ddl := range DDLs {
		_, err := dataBase.Conn.Exec(context.Background(), ddl)
		if err != nil {
			log.Fatalf("can't create a table %e", err)
		}
	}

} // Миграция БД!
