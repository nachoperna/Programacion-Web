-- name: InsertRequest :one 
insert into money_requests (from_alias, to_alias, amount, message) values ($1, $2, $3, $4)
returning *;

-- name: GetRequestsTo :many
select from_alias, to_alias, amount, to_char(date, 'DD-MM-YYYY') as day, to_char(date, 'HH24:MI') as time, message 
from money_requests 
where to_alias = $1;

-- name: GetRequestsFrom :many
select from_alias, to_alias, amount, to_char(date, 'DD-MM-YYYY') as day, to_char(date, 'HH24:MI') as time, message 
from money_requests 
where from_alias = $1;
