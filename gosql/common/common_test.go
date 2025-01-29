package common_test

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-sql/gosql/common"
)

var _ = Describe("common", func() {
	Context("deleteRowByIDString", func() {
		It("convert struct to SQL query", func() {
			var row struct {
				ID    int64  `field:"id" table:"users"`
				Name  string `field:"name"`
				Value string `field:"value"`
			}

			Expect(common.DeleteRowByIDString(&row)).To(Equal(`DELETE FROM users WHERE id = $1`))
		})
	})

	Context("fixQuery", func() {
		It("replace param for MySQL driver", func() {
			sql := "select id, name from users where id=$1"
			Expect(common.FixQuery(sql)).To(Equal("select id, name from users where id=?"))
		})

		It("replace all params for MySQL driver", func() {
			sql := "insert into users set name=$1 where id=$2"
			Expect(common.FixQuery(sql)).To(Equal("insert into users set name=? where id=?"))
		})
	})

	Context("inArray", func() {
		It("check string in array", func() {
			arr := []string{"field1", "field2", "field3"}
			Expect(common.InArray(arr, "field1")).To(BeTrue())
			Expect(common.InArray(arr, "field2")).To(BeTrue())
			Expect(common.InArray(arr, "field3")).To(BeTrue())
			Expect(common.InArray(arr, "field4")).To(BeFalse())
			Expect(common.InArray(arr, "field5")).To(BeFalse())
		})
	})

	Context("insertRowString", func() {
		It("convert struct to SQL query", func() {
			var row struct {
				ID       int64  `field:"id" table:"users"`
				Name     string `field:"name"`
				Value    string `field:"value"`
				Position int64  `field:"position"`
			}

			row.Name = "Name"
			row.Value = "Value"
			row.Position = 59

			sql, args := common.InsertRowString(&row)

			Expect(sql).To(Equal(`INSERT INTO users (name, value, position) VALUES ($1, $2, $3)`))

			Expect(len(args)).To(Equal(3))
			Expect(args[0]).To(Equal("Name"))
			Expect(args[1]).To(Equal("Value"))
			Expect(args[2]).To(Equal(int64(59)))
		})

		It("convert struct to SQL query and populate created_at and updated_at", func() {
			var row struct {
				ID        int64  `field:"id" table:"users"`
				CreatedAt int64  `field:"created_at"`
				UpdatedAt int64  `field:"updated_at"`
				Name      string `field:"name"`
			}

			row.Name = "Name"

			sql, args := common.InsertRowString(&row)

			Expect(sql).To(Equal(`INSERT INTO users (created_at, updated_at, name) VALUES ($1, $2, $3)`))

			Expect(len(args)).To(Equal(3))
			Expect(args[0].(int64) > 0).To(BeTrue())
			Expect(args[1].(int64) > 0).To(BeTrue())
			Expect(args[2]).To(Equal("Name"))
		})
	})

	Context("log", func() {
		Context("time", func() {
			It("calculate one second", func() {
				str := common.Log(io.Discard, "Exec", time.Now().Add(time.Second*-1), nil, false, "")
				Expect(str).To(Equal("\x1b[0;33m[SQL] [func Exec] (empty) (nil) 1.000 ms\x1b[0m\n"))
			})
		})

		Context("format", func() {
			It("with func name", func() {
				str := common.Log(io.Discard, "Exec", time.Now(), nil, false, "")
				Expect(str).To(Equal("\x1b[0;33m[SQL] [func Exec] (empty) (nil) 0.000 ms\x1b[0m\n"))
			})

			It("with sql query", func() {
				str := common.Log(io.Discard, "Exec", time.Now(), nil, false, "select * from users")
				Expect(str).To(Equal("\x1b[0;33m[SQL] [func Exec] select * from users (empty) (nil) 0.000 ms\x1b[0m\n"))
			})

			It("with error message", func() {
				str := common.Log(io.Discard, "Exec", time.Now(), fmt.Errorf("Exec error"), false, "select * from users")
				Expect(str).To(Equal("\x1b[0;31m[SQL] [func Exec] select * from users (empty) (Exec error) 0.000 ms\x1b[0m\n"))
			})

			It("with transaction flag", func() {
				str := common.Log(io.Discard, "Exec", time.Now(), fmt.Errorf("Exec error"), true, "select * from users")
				Expect(str).To(Equal("\x1b[1;31m[SQL] [TX] [func Exec] select * from users (empty) (Exec error) 0.000 ms\x1b[0m\n"))
			})

			It("with sql query arguments", func() {
				str := common.Log(io.Discard, "Exec", time.Now(), fmt.Errorf("Exec error"), true, "select * from users where id=$1", 100)
				Expect(str).To(Equal("\x1b[1;31m[SQL] [TX] [func Exec] select * from users where id=$1 ([100]) (Exec error) 0.000 ms\x1b[0m\n"))
			})
		})
	})

	Context("queryRowByIDString", func() {
		It("convert struct to SQL query", func() {
			var row struct {
				ID    int64  `field:"id" table:"users"`
				Name  string `field:"name"`
				Value string `field:"value"`
			}

			Expect(common.QueryRowByIDString(&row)).To(Equal(`SELECT id, name, value FROM users WHERE id = $1 LIMIT 1`))
		})
	})

	Context("rowExistsString", func() {
		It("convert struct to SQL query", func() {
			var row struct {
				ID    int64  `field:"id" table:"users"`
				Name  string `field:"name"`
				Value string `field:"value"`
			}

			Expect(common.RowExistsString(&row)).To(Equal(`SELECT 1 FROM users WHERE id = $1 LIMIT 1`))
		})
	})

	Context("scans", func() {
		It("convert struct to array of pointers to this struct fields", func() {
			var row struct {
				ID    int64
				Name  string
				Value string
			}

			Expect(common.Scans(&row)).To(Equal([]any{
				&row.ID,
				&row.Name,
				&row.Value,
			}))
		})
	})

	Context("updateRowString", func() {
		It("convert struct to SQL query", func() {
			var row struct {
				ID       int64  `field:"id" table:"users"`
				Name     string `field:"name"`
				Value    string `field:"value"`
				Position int64  `field:"position"`
			}

			row.ID = 10
			row.Name = "Name"
			row.Value = "Value"
			row.Position = 59

			sql, args := common.UpdateRowString(&row)

			Expect(sql).To(Equal(`UPDATE users SET name = $1, value = $2, position = $3 WHERE id = $4`))

			Expect(len(args)).To(Equal(4))
			Expect(args[0]).To(Equal("Name"))
			Expect(args[1]).To(Equal("Value"))
			Expect(args[2]).To(Equal(int64(59)))
			Expect(args[3]).To(Equal(int64(10)))

			sql, args = common.UpdateRowString(&row, "name")

			Expect(sql).To(Equal(`UPDATE users SET name = $1 WHERE id = $2`))

			Expect(len(args)).To(Equal(2))
			Expect(args[0]).To(Equal("Name"))
			Expect(args[1]).To(Equal(int64(10)))

			sql, args = common.UpdateRowString(&row, "name", "value")

			Expect(sql).To(Equal(`UPDATE users SET name = $1, value = $2 WHERE id = $3`))

			Expect(len(args)).To(Equal(3))
			Expect(args[0]).To(Equal("Name"))
			Expect(args[1]).To(Equal("Value"))
			Expect(args[2]).To(Equal(int64(10)))

			sql, args = common.UpdateRowString(&row, "name", "position")

			Expect(sql).To(Equal(`UPDATE users SET name = $1, position = $2 WHERE id = $3`))

			Expect(len(args)).To(Equal(3))
			Expect(args[0]).To(Equal("Name"))
			Expect(args[1]).To(Equal(int64(59)))
			Expect(args[2]).To(Equal(int64(10)))
		})

		It("convert struct to SQL query and populate updated_at", func() {
			var row struct {
				ID        int64  `field:"id" table:"users"`
				CreatedAt int64  `field:"created_at"`
				UpdatedAt int64  `field:"updated_at"`
				Name      string `field:"name"`
			}

			row.ID = 10
			row.Name = "Name"

			sql, args := common.UpdateRowString(&row)

			Expect(sql).To(Equal(`UPDATE users SET updated_at = $1, name = $2 WHERE id = $3`))

			Expect(len(args)).To(Equal(3))
			Expect(args[0].(int64) > 0).To(BeTrue())
			Expect(args[1]).To(Equal("Name"))
			Expect(args[2]).To(Equal(int64(10)))
		})
	})

	Context("ParseUrl", func() {
		Context("Success", func() {
			It("for MySQL", func() {
				// mysql://username:password@127.0.0.1:3306/database?parseTime=true
				// mysql://username:password@/database?socket=/var/run/mysqld/mysqld.sock

				url := "mysql://username:password@127.0.0.1:3306/database?parseTime=true"
				result, err := common.ParseUrl(url)

				Expect(err).To(Succeed())
				Expect(result.Scheme).To(Equal("mysql"))
				Expect(result.User.Username()).To(Equal("username"))

				password, whether := result.User.Password()
				Expect(password).To(Equal("password"))
				Expect(whether).To(BeTrue())

				Expect(result.Host).To(Equal("127.0.0.1:3306"))
				Expect(result.Path).To(Equal("/database"))
				Expect(result.RawQuery).To(Equal("parseTime=true"))
			})

			It("for PostgreSQL", func() {
				// postgres://username:password@127.0.0.1:5432/database?sslmode=disable
				// postgresql://username:password@127.0.0.1:5432/database?sslmode=disable

				url := "postgres://username:password@127.0.0.1:5432/database?sslmode=disable"
				result, err := common.ParseUrl(url)

				Expect(err).To(Succeed())
				Expect(result.Scheme).To(Equal("postgres"))
				Expect(result.User.Username()).To(Equal("username"))

				password, whether := result.User.Password()
				Expect(password).To(Equal("password"))
				Expect(whether).To(BeTrue())

				Expect(result.Host).To(Equal("127.0.0.1:5432"))
				Expect(result.Path).To(Equal("/database"))
				Expect(result.RawQuery).To(Equal("sslmode=disable"))
			})

			It("for SQLite", func() {
				// sqlite:///data/database.sqlite
				// sqlite3:///data/database.sqlite

				url := "sqlite:///data/database.sqlite"
				result, err := common.ParseUrl(url)

				Expect(err).To(Succeed())
				Expect(result.Scheme).To(Equal("sqlite"))
				Expect(result.Host).To(Equal(""))
				Expect(result.Path).To(Equal("/data/database.sqlite"))
			})
		})

		Context("Fail", func() {
			It("for empty", func() {
				_, err := common.ParseUrl("")

				Expect(err).NotTo(Succeed())
				Expect(err.Error()).To(Equal("protocol scheme is not defined"))
			})

			It("for else", func() {
				url := "12345"
				_, err := common.ParseUrl(url)

				Expect(err).NotTo(Succeed())
				Expect(err.Error()).To(Equal("protocol scheme is not defined"))
			})

			It("for not supported", func() {
				url := "example:///some-else"
				_, err := common.ParseUrl(url)

				Expect(err).NotTo(Succeed())
				Expect(err.Error()).To(Equal("unsupported protocol scheme: example"))
			})
		})
	})

	Context("OpenDB", func() {
		var migrationsDir string

		BeforeEach(func() {
			var err error
			migrationsDir, err = filepath.Abs("../../db/migrations")
			Expect(err).To(Succeed())
		})

		Context("Success", func() {
			// // Note: you need to up MySQL server for this test case
			// It("for MySQL", func() {
			// 	databaseURL, err := url.Parse("mysql://root:root@127.0.0.1:3306/gosql")
			// 	Expect(err).To(Succeed())

			// 	db, err := common.OpenDB(databaseURL, migrationsDir, false, false)
			// 	Expect(err).To(Succeed())
			// 	Expect(db.Close()).To(Succeed())
			// })

			// // Note: you need to up PostgreSQL server for this test case
			// It("for PostgreSQL", func() {
			// 	databaseURL, err := url.Parse("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable")
			// 	Expect(err).To(Succeed())

			// 	db, err := common.OpenDB(databaseURL, migrationsDir, false, false)
			// 	Expect(err).To(Succeed())
			// 	Expect(db.Close()).To(Succeed())
			// })

			It("for SQLite", func() {
				f, err := os.CreateTemp("", "go-sqlite3-test-")
				Expect(err).To(Succeed())
				f.Close()

				databaseURL, err := url.Parse("sqlite://" + f.Name())
				Expect(err).To(Succeed())

				db, err := common.OpenDB(databaseURL, migrationsDir, false, false)
				Expect(err).To(Succeed())
				Expect(db.Close()).To(Succeed())
			})
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gosql/common")
}
