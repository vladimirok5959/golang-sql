# golang-sql

Go bindings library for MySQL, PostgreSQL and SQLite

Used [amacneil/dbmate](https://github.com/amacneil/dbmate) inside the project for creating connection (please review dbmate docs), so next schemes is supported:

```sh
# MySQL TCP connection:
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
