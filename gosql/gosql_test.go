package gosql_test

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-sql/gosql"
)

var _ = Describe("gosql", func() {
	Context("Open", func() {
		var migrationsDir string
		var ctx = context.Background()
		var sql = "select id, name from users where id=$1"

		var (
			id   int
			name string
		)

		BeforeEach(func() {
			var err error
			migrationsDir, err = filepath.Abs("../db/migrations")
			Expect(err).To(Succeed())
		})

		// // Note: you need to up MySQL server for this test case
		// Context("for MySQL", func() {
		// 	It("open connection, migrate and select data", func() {
		// 		db, err := gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		err = db.QueryRow(ctx, sql, 1).Scan(&id, &name)
		// 		Expect(err).To(Succeed())
		// 		Expect(id).To(Equal(1))
		// 		Expect(name).To(Equal("Alice"))

		// 		err = db.QueryRow(ctx, sql, 2).Scan(&id, &name)
		// 		Expect(err).To(Succeed())
		// 		Expect(id).To(Equal(2))
		// 		Expect(name).To(Equal("Bob"))

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and select by ID", func() {
		// 		db, err := gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		err = db.QueryRowByID(ctx, 1, &rowUser)
		// 		Expect(err).To(Succeed())
		// 		Expect(rowUser.ID).To(Equal(int64(1)))
		// 		Expect(rowUser.Name).To(Equal("Alice"))

		// 		Expect(db.Close()).To(Succeed())
		// 	})
		// })

		// // Note: you need to up PostgreSQL server for this test case
		// Context("for PostgreSQL", func() {
		// 	It("open connection, migrate and select data", func() {
		// 		db, err := gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		err = db.QueryRow(ctx, sql, 1).Scan(&id, &name)
		// 		Expect(err).To(Succeed())
		// 		Expect(id).To(Equal(1))
		// 		Expect(name).To(Equal("Alice"))

		// 		err = db.QueryRow(ctx, sql, 2).Scan(&id, &name)
		// 		Expect(err).To(Succeed())
		// 		Expect(id).To(Equal(2))
		// 		Expect(name).To(Equal("Bob"))

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and select by ID", func() {
		// 		db, err := gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		err = db.QueryRowByID(ctx, 1, &rowUser)
		// 		Expect(err).To(Succeed())
		// 		Expect(rowUser.ID).To(Equal(int64(1)))
		// 		Expect(rowUser.Name).To(Equal("Alice"))

		// 		Expect(db.Close()).To(Succeed())
		// 	})
		// })

		Context("for SQLite", func() {
			It("open connection, migrate and select data", func() {
				f, err := ioutil.TempFile("", "go-sqlite-test-")
				Expect(err).To(Succeed())
				f.Close()

				db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, false)
				Expect(err).To(Succeed())

				err = db.QueryRow(ctx, sql, 1).Scan(&id, &name)
				Expect(err).To(Succeed())
				Expect(id).To(Equal(1))
				Expect(name).To(Equal("Alice"))

				err = db.QueryRow(ctx, sql, 2).Scan(&id, &name)
				Expect(err).To(Succeed())
				Expect(id).To(Equal(2))
				Expect(name).To(Equal("Bob"))

				Expect(db.Close()).To(Succeed())
			})

			It("open connection, migrate and select by ID", func() {
				f, err := ioutil.TempFile("", "go-sqlite-test-")
				Expect(err).To(Succeed())
				f.Close()

				db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, false)
				Expect(err).To(Succeed())

				var rowUser struct {
					ID   int64  `field:"id" table:"users"`
					Name string `field:"name"`
				}

				err = db.QueryRowByID(ctx, 1, &rowUser)
				Expect(err).To(Succeed())
				Expect(rowUser.ID).To(Equal(int64(1)))
				Expect(rowUser.Name).To(Equal("Alice"))

				Expect(db.Close()).To(Succeed())
			})
		})

		It("open connection and skip migration", func() {
			f, err := ioutil.TempFile("", "go-sqlite-test-")
			Expect(err).To(Succeed())
			f.Close()

			db, err := gosql.Open("sqlite://"+f.Name(), "", true, false)
			Expect(err).To(Succeed())
			Expect(db.Ping(ctx)).To(Succeed())

			var size int
			err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
			Expect(err.Error()).To(Equal("no such table: users"))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gosql")
}
