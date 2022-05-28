package common_test

import (
	"io/ioutil"
	"net/url"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-sql/gosql/common"
)

var _ = Describe("common", func() {
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

			// 	db, err := common.OpenDB(databaseURL, migrationsDir)
			// 	Expect(err).To(Succeed())
			// 	Expect(db.Close()).To(Succeed())
			// })

			// // Note: you need to up PostgreSQL server for this test case
			// It("for PostgreSQL", func() {
			// 	databaseURL, err := url.Parse("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable")
			// 	Expect(err).To(Succeed())

			// 	db, err := common.OpenDB(databaseURL, migrationsDir)
			// 	Expect(err).To(Succeed())
			// 	Expect(db.Close()).To(Succeed())
			// })

			It("for SQLite", func() {
				f, err := ioutil.TempFile("", "go-sqlite3-test-")
				Expect(err).To(Succeed())
				f.Close()

				databaseURL, err := url.Parse("sqlite://" + f.Name())
				Expect(err).To(Succeed())

				db, err := common.OpenDB(databaseURL, migrationsDir)
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
