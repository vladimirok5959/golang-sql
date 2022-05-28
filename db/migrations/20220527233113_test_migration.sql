-- migrate:up
create table users (
  id integer,
  name varchar(255)
);
insert into users (id, name) values (1, 'Alice');
insert into users (id, name) values (2, 'Bob');

-- migrate:down
drop table users;
