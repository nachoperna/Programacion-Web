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
update users set
      name = coalesce(nullif(sqlc.arg('name'),''), name), 
      email = coalesce(nullif(sqlc.arg('email'),''), email), 
      password = coalesce(nullif(sqlc.arg('password'),''), password)
where alias = $1
returning *;

-- name: DeleteAllUsers :exec
delete from users;

