-- name: CreateRoom :one
INSERT INTO room (room_unique,
                  user1,
                  user2)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoomByUsers :one
SELECT room_unique
FROM room
WHERE user_id = ANY ($1::int[]);

-- name: GetAllExistingRooms :many

SELECT * FROM room ;

-- name: DeleteAllRooms :exec
DELETE from room ;