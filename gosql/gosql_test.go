package gosql_test

import (
	"context"
	"os"
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

			// // Reset databases
			// // Note: uncomment for MySQL and PostgreSQL tests
			// var db common.Engine

			// // MySQL
			// db, err = gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", "", true, false)
			// Expect(err).To(Succeed())
			// _, _ = db.Exec(ctx, "DROP TABLE schema_migrations, users")

			// // PostgreSQL
			// db, err = gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", "", true, false)
			// Expect(err).To(Succeed())
			// _, _ = db.Exec(ctx, "DROP TABLE schema_migrations, users")
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

		// 	It("open connection, migrate and check row", func() {
		// 		db, err := gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		Expect(db.RowExists(ctx, 1, &rowUser)).To(BeTrue())
		// 		Expect(db.RowExists(ctx, 2, &rowUser)).To(BeTrue())
		// 		Expect(db.RowExists(ctx, 3, &rowUser)).To(BeFalse())
		// 		Expect(db.RowExists(ctx, 4, &rowUser)).To(BeFalse())
		// 		Expect(db.RowExists(ctx, 5, &rowUser)).To(BeFalse())

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and delete row", func() {
		// 		db, err := gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		var size int

		// 		Expect(db.DeleteRowByID(ctx, 2, &rowUser)).To(Succeed())
		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(1))

		// 		Expect(db.DeleteRowByID(ctx, 1, &rowUser)).To(Succeed())
		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(0))

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and insert row", func() {
		// 		db, err := gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		var size int
		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(2))

		// 		rowUser.Name = "James"
		// 		Expect(db.InsertRow(ctx, &rowUser)).To(Succeed())

		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(3))

		// 		rowUser.Name = "Robert"
		// 		Expect(db.InsertRow(ctx, &rowUser)).To(Succeed())

		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(4))

		// 		names := []string{}
		// 		err = db.Each(
		// 			ctx,
		// 			"SELECT 0, name FROM users ORDER BY name ASC",
		// 			func(ctx context.Context, rows *gosql.Rows) error {
		// 				if err := rows.Scans(&rowUser); err != nil {
		// 					return err
		// 				}
		// 				names = append(names, rowUser.Name)
		// 				return nil
		// 			},
		// 		)
		// 		Expect(err).To(Succeed())
		// 		Expect(names).To(Equal([]string{"Alice", "Bob", "James", "Robert"}))

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and update row", func() {
		// 		db, err := gosql.Open("mysql://root:root@127.0.0.1:3306/gosql", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		_, err = db.Exec(ctx, "INSERT INTO users (id, name) VALUES ('3', 'Patrik')")
		// 		Expect(err).To(Succeed())

		// 		rowUser.ID = 3
		// 		rowUser.Name = "Robert"
		// 		Expect(db.UpdateRow(ctx, &rowUser)).To(Succeed())

		// 		err = db.QueryRowByID(ctx, 3, &rowUser)
		// 		Expect(err).To(Succeed())
		// 		Expect(rowUser.ID).To(Equal(int64(3)))
		// 		Expect(rowUser.Name).To(Equal("Robert"))

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

		// 	It("open connection, migrate and check row", func() {
		// 		db, err := gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		Expect(db.RowExists(ctx, 1, &rowUser)).To(BeTrue())
		// 		Expect(db.RowExists(ctx, 2, &rowUser)).To(BeTrue())
		// 		Expect(db.RowExists(ctx, 3, &rowUser)).To(BeFalse())
		// 		Expect(db.RowExists(ctx, 4, &rowUser)).To(BeFalse())
		// 		Expect(db.RowExists(ctx, 5, &rowUser)).To(BeFalse())

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and delete row", func() {
		// 		db, err := gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		var size int

		// 		Expect(db.DeleteRowByID(ctx, 2, &rowUser)).To(Succeed())
		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(1))

		// 		Expect(db.DeleteRowByID(ctx, 1, &rowUser)).To(Succeed())
		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(0))

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and insert row", func() {
		// 		db, err := gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		var size int
		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(2))

		// 		rowUser.Name = "James"
		// 		Expect(db.InsertRow(ctx, &rowUser)).To(Succeed())

		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(3))

		// 		rowUser.Name = "Robert"
		// 		Expect(db.InsertRow(ctx, &rowUser)).To(Succeed())

		// 		err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
		// 		Expect(err).To(Succeed())
		// 		Expect(size).To(Equal(4))

		// 		names := []string{}
		// 		err = db.Each(
		// 			ctx,
		// 			"SELECT 0, name FROM users ORDER BY name ASC",
		// 			func(ctx context.Context, rows *gosql.Rows) error {
		// 				if err := rows.Scans(&rowUser); err != nil {
		// 					return err
		// 				}
		// 				names = append(names, rowUser.Name)
		// 				return nil
		// 			},
		// 		)
		// 		Expect(err).To(Succeed())
		// 		Expect(names).To(Equal([]string{"Alice", "Bob", "James", "Robert"}))

		// 		Expect(db.Close()).To(Succeed())
		// 	})

		// 	It("open connection, migrate and update row", func() {
		// 		db, err := gosql.Open("postgres://root:root@127.0.0.1:5432/gosql?sslmode=disable", migrationsDir, false, false)
		// 		Expect(err).To(Succeed())

		// 		var rowUser struct {
		// 			ID   int64  `field:"id" table:"users"`
		// 			Name string `field:"name"`
		// 		}

		// 		_, err = db.Exec(ctx, "INSERT INTO users (id, name) VALUES ('3', 'Patrik')")
		// 		Expect(err).To(Succeed())

		// 		rowUser.ID = 3
		// 		rowUser.Name = "Robert"
		// 		Expect(db.UpdateRow(ctx, &rowUser)).To(Succeed())

		// 		err = db.QueryRowByID(ctx, 3, &rowUser)
		// 		Expect(err).To(Succeed())
		// 		Expect(rowUser.ID).To(Equal(int64(3)))
		// 		Expect(rowUser.Name).To(Equal("Robert"))

		// 		Expect(db.Close()).To(Succeed())
		// 	})
		// })

		Context("for SQLite", func() {
			It("open connection, migrate and select data", func() {
				f, err := os.CreateTemp("", "go-sqlite-test-")
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
				f, err := os.CreateTemp("", "go-sqlite-test-")
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

			It("open connection, migrate and check row", func() {
				f, err := os.CreateTemp("", "go-sqlite-test-")
				Expect(err).To(Succeed())
				f.Close()

				db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, false)
				Expect(err).To(Succeed())

				var rowUser struct {
					ID   int64  `field:"id" table:"users"`
					Name string `field:"name"`
				}

				Expect(db.RowExists(ctx, 1, &rowUser)).To(BeTrue())
				Expect(db.RowExists(ctx, 2, &rowUser)).To(BeTrue())
				Expect(db.RowExists(ctx, 3, &rowUser)).To(BeFalse())
				Expect(db.RowExists(ctx, 4, &rowUser)).To(BeFalse())
				Expect(db.RowExists(ctx, 5, &rowUser)).To(BeFalse())

				Expect(db.Close()).To(Succeed())
			})

			It("open connection, migrate and delete row", func() {
				f, err := os.CreateTemp("", "go-sqlite-test-")
				Expect(err).To(Succeed())
				f.Close()

				db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, false)
				Expect(err).To(Succeed())

				var rowUser struct {
					ID   int64  `field:"id" table:"users"`
					Name string `field:"name"`
				}

				var size int

				Expect(db.DeleteRowByID(ctx, 2, &rowUser)).To(Succeed())
				err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
				Expect(err).To(Succeed())
				Expect(size).To(Equal(1))

				Expect(db.DeleteRowByID(ctx, 1, &rowUser)).To(Succeed())
				err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
				Expect(err).To(Succeed())
				Expect(size).To(Equal(0))

				Expect(db.Close()).To(Succeed())
			})

			It("open connection, migrate and insert row", func() {
				f, err := os.CreateTemp("", "go-sqlite-test-")
				Expect(err).To(Succeed())
				f.Close()

				db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, false)
				Expect(err).To(Succeed())

				var rowUser struct {
					ID   int64  `field:"id" table:"users"`
					Name string `field:"name"`
				}

				var size int
				err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
				Expect(err).To(Succeed())
				Expect(size).To(Equal(2))

				rowUser.Name = "James"
				Expect(db.InsertRow(ctx, &rowUser)).To(Succeed())

				err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
				Expect(err).To(Succeed())
				Expect(size).To(Equal(3))

				rowUser.Name = "Robert"
				Expect(db.InsertRow(ctx, &rowUser)).To(Succeed())

				err = db.QueryRow(ctx, "select count(*) from users").Scan(&size)
				Expect(err).To(Succeed())
				Expect(size).To(Equal(4))

				names := []string{}
				err = db.Each(
					ctx,
					"SELECT 0, name FROM users ORDER BY name ASC",
					func(ctx context.Context, rows *gosql.Rows) error {
						if err := rows.Scans(&rowUser); err != nil {
							return err
						}
						names = append(names, rowUser.Name)
						return nil
					},
				)
				Expect(err).To(Succeed())
				Expect(names).To(Equal([]string{"Alice", "Bob", "James", "Robert"}))

				Expect(db.Close()).To(Succeed())
			})

			It("open connection, migrate and update row", func() {
				f, err := os.CreateTemp("", "go-sqlite-test-")
				Expect(err).To(Succeed())
				f.Close()

				db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, false)
				Expect(err).To(Succeed())

				var rowUser struct {
					ID   int64  `field:"id" table:"users"`
					Name string `field:"name"`
				}

				_, err = db.Exec(ctx, "INSERT INTO users (id, name) VALUES ('3', 'Patrik')")
				Expect(err).To(Succeed())

				rowUser.ID = 3
				rowUser.Name = "Robert"
				Expect(db.UpdateRow(ctx, &rowUser)).To(Succeed())

				err = db.QueryRowByID(ctx, 3, &rowUser)
				Expect(err).To(Succeed())
				Expect(rowUser.ID).To(Equal(int64(3)))
				Expect(rowUser.Name).To(Equal("Robert"))

				Expect(db.Close()).To(Succeed())
			})
		})

		It("open connection and skip migration", func() {
			f, err := os.CreateTemp("", "go-sqlite-test-")
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
