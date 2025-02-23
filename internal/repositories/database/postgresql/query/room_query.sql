-- name: CreateRoom :one
INSERT INTO room (room_unique,
                  user1,
                  user2)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoomByUsers :one
SELECT room_unique
FROM room
WHERE (user1 = $1 and user2 = $2)
     or
      (user1 = $2 and user2 = $1);