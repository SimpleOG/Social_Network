CREATE TABLE if not exists images(
  id serial not null primary key,
  user_id integer not null references users(id),
  created_at timestamptz not null default now()
);