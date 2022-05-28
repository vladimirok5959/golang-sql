-- migrate:up
create table users (
  id integer,
  name varchar(255)
);
insert into users (id, name) values (1, 'alice');
insert into users (id, name) values (2, 'bob');

-- migrate:down
drop table users;
