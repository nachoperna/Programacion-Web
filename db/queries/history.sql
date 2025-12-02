-- name: UpdateHistory :exec
insert into movements_history 
      (alias, type, amount, time)
      values ($1, $2, $3, now());

-- name: GetHistory :many
select 
      type,
      amount,
      to_char(time, 'DD-MM-YYYY') as day,
      to_char(time, 'HH24:MI') as time
      from movements_history where alias = $1
      order by time desc;

