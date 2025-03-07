CREATE TABLE if not exists "room"
(
    id          serial       not null primary key,
    room_unique varchar(255) unique not null,
    user1       integer      not null references "users" (id),
    user2       integer      not null references "users" (id)
);
Create index on room (room_unique);

CREATE TABLE if not exists  "messages"
(
    id              serial  not null primary key,
    room_id         text not null references room (room_unique),
    message_content text    not null,
    message_owner   integer not null references users (id),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    was_delivered   boolean not null
);
