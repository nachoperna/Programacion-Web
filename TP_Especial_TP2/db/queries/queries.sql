-- name: GetUser :one
select * from users where alias = $1;

-- name: DeleteUser :exec
delete from users where alias = $1;

-- name: ListUsers :many 
select * from users;

-- name: InsertUser :exec
insert into users 
      (alias, name, email, password)
      values ($1, $2, $3, $4);

-- name: UpdateUser :exec
update users set name = $2, email = $3, password = $4 where alias = $1;

-- name: GetAccount :one
select * from accounts where alias = $1;

-- name: InsertAccount :exec
insert into accounts (alias, balance) values ($1, default);

-- name: UpdateAccount :exec
update accounts set 
      balance=$2, 
      last_movement_type=$3, 
      last_deposit=$4,
      last_deposit_amount=$5,
      last_transfer=$6,
      last_transfer_account=$7,
      last_transfer_amount=$8,
      last_whidrawal=$9,
      last_whitdrawal_amount=$10 
      where alias = $1;

-- name: DeleteAccount :exec
delete from accounts where alias = $1;

-- name: ListAccounts :many
select * from accounts;

-- name: DeleteAll :exec
delete from users;
