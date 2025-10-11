CREATE TABLE users(
      alias VARCHAR(30) unique not null,
      name varchar(80) not null,
      email VARCHAR(50) unique not null,
      password varchar(50) not null,
      signed_up timestamp default current_timestamp not null,
      last_session timestamp,
      CONSTRAINT Users_pk PRIMARY KEY(alias)
);

CREATE TABLE accounts(
      alias varchar(30),
      balance numeric(8,2) not null default 0.00,
      last_movement_type text, -- DEPOSIT, TRANSFER, WHITDRAWAL
      last_deposit timestamp,
      last_deposit_amount numeric(8,2),
      last_transfer timestamp,
      last_transfer_account varchar(30),
      last_transfer_amount numeric(8,2),
      last_whidrawal timestamp,
      last_whitdrawal_amount numeric(8,2),
      CONSTRAINT Accounts_pk PRIMARY KEY(alias)
);

ALTER TABLE accounts ADD CONSTRAINT fk_user_account
      FOREIGN KEY(alias) REFERENCES users(alias)
      NOT DEFERRABLE
      INITIALLY IMMEDIATE
;

-- Puede existir tambien una tabla Transferencia donde tenga informacion de cuenta origen, destino, fechas y montos
