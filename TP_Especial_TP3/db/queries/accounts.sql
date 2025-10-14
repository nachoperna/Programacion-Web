-- name: GetAccount :one
select * from accounts where alias = $1;

-- name: InsertAccount :exec
insert into accounts (alias, balance) values ($1, default);

-- name: UpdateAccount :exec
update accounts set 
      balance=balance+$2, 
      last_movement_type=$3, 
      last_deposit=$4,
      last_deposit_amount=$5,
      last_transfer=$6,
      last_transfer_account=$7,
      last_transfer_amount=$8,
      last_withdrawal=$9,
      last_withdrawal_amount=$10 
      where alias = $1;

-- name: DeleteAccount :exec
delete from accounts where alias = $1;

-- name: ListAccounts :many
select * from accounts;

-- name: Deposit :one
update accounts set 
      balance = balance + $2, 
      last_movement_type = 'Deposito',
      last_deposit = current_timestamp, 
      last_deposit_amount = $2
      where alias = $1
returning *;

-- name: Withdrawal :exec
update accounts set 
      balance = balance - $2,
      last_movement_type = 'Retiro',
      last_withdrawal = current_timestamp,
      last_withdrawal_amount = $2
      where alias = $1;

-- name: Transfer :exec
update accounts set 
      balance = balance - $3,
      last_movement_type = 'Transferencia',
      last_transfer = current_timestamp,
      last_transfer_account = $2,
      last_transfer_amount = $3
      where alias = $1;

-- name: GetBalance :one
select balance, last_movement_type from accounts where alias = $1;
