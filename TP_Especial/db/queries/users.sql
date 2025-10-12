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

-- name: DeleteAllUsers :exec
delete from users;

