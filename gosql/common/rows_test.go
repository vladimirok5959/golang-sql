package common_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-sql/gosql/common"
)

var _ = Describe("common", func() {
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
})
