# golang-sql

## Features

- Multiple storages MySQL, PostgreSQL, SQLite
- Debug mode for SQL queries
- Transactions is supported
- Transactions helpers
- Database migrations
- Simple usage

Go bindings library for MySQL, PostgreSQL and SQLite

This library allow you to faster make your development if you need to use/support multiple database engines such MySQL, PostgreSQL and SQLite. For example you can make web service and give ability to choose storage type. Or you can use SQLite for tests or for demo version and MySQL or PostgreSQL for production. Migrations (thanks to dbmate) is suported out of box and full SQL messages for debugging. Note: please use PostgreSQL parameter placeholders even if MySQL is used, it will be automatically replaced with `?`

Used [amacneil/dbmate](https://github.com/amacneil/dbmate) inside the project for creating connection (please review dbmate docs) and for migrations, so next schemes is supported:

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
