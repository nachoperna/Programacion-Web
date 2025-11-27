-- name: InsertRequest :one 
insert into money_requests (from_alias, to_alias, amount, message) values ($1, $2, $3, $4)
returning *;

-- name: GetRequestsTo :many
select from_alias, to_alias, amount, to_char(date, 'DD-MM-YYYY') as day, to_char(date, 'HH24:MI') as time, message 
from money_requests 
where to_alias = $1
ORDER BY
      -- Ordenamiento ascendente
      CASE WHEN sqlc.arg('sort_by') = 'from_alias' AND sqlc.arg('sort_order') = 'asc' THEN from_alias END ASC,
      CASE WHEN sqlc.arg('sort_by') = 'amount' AND sqlc.arg('sort_order') = 'asc' THEN amount END ASC,
      CASE WHEN sqlc.arg('sort_by') = 'day' AND sqlc.arg('sort_order') = 'asc' THEN date END ASC,
      -- Ordenamiento descendente
      CASE WHEN sqlc.arg('sort_by') = 'from_alias' AND sqlc.arg('sort_order') = 'desc' THEN from_alias END DESC,
      CASE WHEN sqlc.arg('sort_by') = 'amount' AND sqlc.arg('sort_order') = 'desc' THEN amount END DESC,
      CASE WHEN sqlc.arg('sort_by') = 'day' AND sqlc.arg('sort_order') = 'desc' THEN date END DESC;

-- name: GetRequestsFrom :many
select from_alias, to_alias, amount, to_char(date, 'DD-MM-YYYY') as day, to_char(date, 'HH24:MI') as time, message 
from money_requests 
where from_alias = $1
ORDER BY
      -- Ordenamiento ascendente
      CASE WHEN sqlc.arg('sort_by') = 'from_alias' AND sqlc.arg('sort_order') = 'asc' THEN from_alias END ASC,
      CASE WHEN sqlc.arg('sort_by') = 'amount' AND sqlc.arg('sort_order') = 'asc' THEN amount END ASC,
      CASE WHEN sqlc.arg('sort_by') = 'day' AND sqlc.arg('sort_order') = 'asc' THEN date END ASC,
      -- Ordenamiento descendente
      CASE WHEN sqlc.arg('sort_by') = 'from_alias' AND sqlc.arg('sort_order') = 'desc' THEN from_alias END DESC,
      CASE WHEN sqlc.arg('sort_by') = 'amount' AND sqlc.arg('sort_order') = 'desc' THEN amount END DESC,
      CASE WHEN sqlc.arg('sort_by') = 'day' AND sqlc.arg('sort_order') = 'desc' THEN date END DESC;

-- name: DeleteRequest :one
delete from money_requests where from_alias = $1 and to_alias = $2 
returning *;

