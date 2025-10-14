-- name: GetUser :one
select * from users where alias = $1;

-- name: DeleteUser :one
delete from users where alias = $1
returning *;

-- name: ListUsers :many 
select * from users;

-- name: InsertUser :one
insert into users 
      (alias, name, email, password)
      values ($1, $2, $3, $4)
returning *;

-- name: UpdateUser :one
update users set name = $2, email = $3, password = $4 where alias = $1
returning *;

-- name: DeleteAllUsers :exec
delete from users;

