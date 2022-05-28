# golang-sql

Go bindings library for MySQL, PostgreSQL and SQLite

## Features

- Multiple storages: MySQL, PostgreSQL, SQLite
- Debug mode for tracking SQL queries
- Transactions is supported
- Transactions helpers
- Database migrations
- Simple usage

This library allow you to faster make your development if you need to use/support multiple database engines such MySQL, PostgreSQL or/and SQLite. For example you can make web service and give ability to choose storage type. Or you can use SQLite for tests or for demo version and MySQL or PostgreSQL for production. Migrations (thanks to dbmate) is suported out of box and full SQL messages for debugging. Note: please use PostgreSQL parameter placeholders even if MySQL is used, it will be automatically replaced with `?`

Used [amacneil/dbmate](https://github.com/amacneil/dbmate) inside the project for creating connections (please review dbmate docs) and for migrations, so next schemes is supported:

```sh
# MySQL connection:
mysql://username:password@127.0.0.1:3306/database?parseTime=true

# MySQL through the socket:
mysql://username:password@/database?socket=/var/run/mysqld/mysqld.sock

# PostgreSQL connection:
postgres://username:password@127.0.0.1:5432/database?sslmode=disable
postgresql://username:password@127.0.0.1:5432/database?sslmode=disable

# SQLite file:
sqlite:///data/database.sqlite
sqlite3:///data/database.sqlite
```

## Examples

```sh
$ go run main.go 
Insert some data to users table
[SQL] [func Exec] INSERT INTO users (id, name) VALUES ($1, $2) ([5 John]) (nil) 0.005 ms
Select all rows from users table
[SQL] [func Query] SELECT id, name FROM users ORDER BY id ASC (empty) (nil) 0.000 ms
ID: 1, Name: Alice
ID: 2, Name: Bob
ID: 5, Name: John
Update inside transaction
[SQL] [TX] [func Begin] (empty) (nil) 0.000 ms
[SQL] [TX] [func Exec] UPDATE users SET name=$1 WHERE id=$2 ([John 1]) (nil) 0.000 ms
[SQL] [TX] [func Exec] UPDATE users SET name=$1 WHERE id=$2 ([Alice 5]) (nil) 0.000 ms
[SQL] [TX] [func Commit] (empty) (nil) 0.004 ms
Select all rows from users again
[SQL] [func Query] SELECT id, name FROM users ORDER BY id ASC (empty) (nil) 0.000 ms
ID: 1, Name: John
ID: 2, Name: Bob
ID: 5, Name: Alice
[SQL] [func Close] (empty) (nil) 0.000 ms
```
