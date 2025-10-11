CREATE TABLE users (
      id serial primary key,
      name varchar(255) not null,
      email varchar(255) not null,
      created_at timestamp with time zone default current_timestamp
);

