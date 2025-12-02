-- name: GetAccount :one
select * from accounts where alias = $1;

-- name: InsertAccount :exec
insert into accounts (alias, balance) values ($1, default);

-- name: DeleteAccount :exec
delete from accounts where alias = $1;

-- name: ListAccounts :many
select * from accounts;

-- name: Deposit :one
update accounts set 
      balance = balance + $2, 
      last_movement_type = $3,
      last_movement_amount = $4
      where alias = $1
returning balance, last_movement_type;

-- name: Withdrawal :one
update accounts set 
      balance = balance - $2,
      last_movement_type = $3,
      last_movement_amount = $4
      where alias = $1
returning balance, last_movement_type;

-- name: GetBalance :one
select balance, last_movement_type from accounts where alias = $1;
