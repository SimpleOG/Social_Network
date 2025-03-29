CREATE TABLE if not exists "users"
(
    "id"         SERIAL              NOT NULL UNIQUE primary key,
    "username"   VARCHAR(255) UNIQUE NOT NULL,
    "password"   VARCHAR(255)        NOT NULL,
    "created_at" timestamptz         NOT NULL DEFAULT (now())
);
CREATE TABLE if not exists "room"
(
    id          serial              not null primary key,
    room_unique varchar(255) unique not null
);

CREATE TABLE if not exists "messages"
(
    id              serial      not null primary key,
    message_content text        not null,
    message_owner   integer     not null references users (id),
    "created_at"    timestamptz NOT NULL DEFAULT (now()),
    was_delivered   boolean     not null
);
CREATE TABLE if not exists "room_messages"
(
    id      serial not null primary key,
    room_id varchar(255)   not null references room (room_unique),
    message_id integer not null references messages(id)
);
CREATE TABLE if not exists users_room
(
    user_id integer not null references "users" (id),
    room_id varchar(255) not null references room (room_unique)
);
CREATE TABLE if not exists images
(
    id             serial      not null primary key,
    user_id        integer     not null references users (id),
    created_at     timestamptz not null default now(),
    filestorage_id integer     not null,
    message_id     integer     not null references messages (id)
);