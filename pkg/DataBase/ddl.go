package database

const ClientsDDL = `create table if not exists client (
id serial primary key,
name text not null,
surname text not null,
login text not null unique,
password text not null,
phone text not null,
create_data timestamp not null,
update_data timestamp not null
);
`

const AccountsDDL = `create table if not exists accounts (
  id serial primary key,
  client_id int references client(id) on delete cascade,
  merchant_id int references merchants(id) on delete cascade,
  account_name text not null,
  network_payment text not null,	
  balance int default 0,
  currency text not null,
  create_data timestamp not null,
  update_data timestamp not null
);
`
const MobiAccountDDL = `create table if not exists account_online (
  id serial primary key,
  client_id int references client(id) on delete cascade,
  merchant_id int references merchants(id) on delete cascade,
  account_name text default 'Mobi Account',
  network_payment text default 'Mobi',	
  balance int default 0,
  currency text default 'TJS',
  create_data timestamp not null,
  update_data timestamp not null
);
`
const MerchantsDDL = `create table if not exists merchants
(
  id serial primary key,
  company_name text not null,
  company_address text not null,
  company_description text not null,
  owner_name varchar(20) not null,
  owner_surname varchar(20) not null,
  login text not null unique,
  password text not null,
  create_data timestamp not null,
  update_data timestamp not null
)`

const ServicesDDL = `create table if not exists Services
(
 id serial primary key,
 service_merchant_id  int not null references merchants(id) on delete cascade,
 service_name text not null,
 service_category text not null,	
 company_service_name text not null,
 service_Price int default 1,
 currency text not null,
 type text not null,
 service_create_data timestamp not null,	
 service_update_data timestamp not null 
)`

const BankCardDDl = `create table if not exists bank_account_client
(
 id serial primary key,
 client_id int references client(id) on delete cascade,
 bank_name text not null,
 bankcard_name text not null,	
 bin_number text not null unique,
 cvc text not null,
 balance int default '100000',
 payment_system text not null, 
 currency text not null,
 create_data timestamp not null,	
 update_data timestamp not null 
)`

const BankAccountDDL = `create table if not exists bank_account_merchant
(
 id serial primary key,
 merchant_id  int references merchants(id) on delete cascade,
 bank_name text not null,
 iban text not null unique,
 inn text not null,
 orgn text not null,
balance int default '100000',
 currency text not null,
 create_data timestamp not null,	
 update_data timestamp not null 
)`


const ClientServicesDDL = `create table if not exists Services_client
(
 id serial primary key,
 client_id  int not null references client(id) on delete cascade,
 service_name text not null,
service_category text not null,
 company_service_name text not null,
 service_Price int default 1,
 currency text not null,
 type text not null,
 state text default 'active', 
 authorized text default 'not confirmed', 
 service_create_data timestamp not null,	
 service_update_data timestamp not null 
)`

const ATMsDDL = `create table if not exists ATMs
(
  id serial primary key,
  Address text not null,
  Balance int default 0,
  create_data timestamp not null,
  update_data timestamp not null
)`
const ManagerDDL = `create table if not exists manager
(
  id serial primary key,
  name text not null,
  surname text not null,
  login text not null unique,
  password text not null,
  phone text not null,
  create_data timestamp not null,
  update_data timestamp not null
)`
const TransactionsMerchantsDDL = `create table if not exists TransactionsMerchants
(
  id serial primary key,
  sender_account int not null references client(id) on delete cascade,
  receiver_account int not null references merchants(id) on delete cascade,
  transaction_amount int default 0,
  create_data timestamp not null
)`

const TransactionsClientsDDL = `create table if not exists TransactionsClients
(
  id serial primary key,
  sender_account int not null references client(id) on delete cascade,
  receiver_account int not null references client(id) on delete cascade,
  transaction_amount int default 0,
  create_data timestamp not null
)`

const DebitCreditDDL = `create table if not exists Debit_Credit
(
  id serial primary key,
  client_id int references client(id) on delete cascade,
  merchant_id int references client(id) on delete cascade,
  transaction_type text not null, 
  account_name text not null,
  debit int default 0,
  credit int default 0,
  commission int default 0,
  total_amount_before int default 0,
  total_amount_after int default 0,
  data_transaction timestamp not null
)`
