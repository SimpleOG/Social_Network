-- name: CreateMessage :exec

INSERT INTO messages (room_id,
                      message_content,
                      message_owner,
                      was_delivered)
values ($1, $2, $3, $4);

-- name: GetMessagesForRoom :many

SELECT *
from messages
where room_id = $1;

-- name: GetAllUndeliveredMessages :many

SELECT *
from messages
where room_id = $1 and message_owner !=$2
  and was_delivered is false
order by created_at  ;

-- name: ChangeDeliveryTipe :exec

UPDATE messages
set was_delivered = true where id=$1;

