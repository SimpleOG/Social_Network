CREATE TABLE "room"
(
    id          serial       not null primary key,
    room_unique varchar(255) not null,
    user1       integer      not null references "users" (id),
    user2       integer      not null references "users" (id)
);
Create index on room(room_unique);