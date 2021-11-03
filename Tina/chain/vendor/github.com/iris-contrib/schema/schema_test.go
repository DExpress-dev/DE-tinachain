// package schema_test black box testing
package schema_test

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

type testRequest struct {
	Username string     `form:"test_username" json:"test_username"`
	Data     []testData `form:"test_data" json:"test_data"`
	Other    int        `form:"test_other" json:"test_other"`
}

type testData struct {
	Name string `form:"test_data_name" json:"test_data_name"`
}

func TestSchemaSliceFieldTag(t *testing.T) {
	app := iris.New()
	app.Post("/", func(ctx iris.Context) {
		var p testRequest
		if err := ctx.ReadForm(&p); err != nil && iris.IsErrPath(err) {
			t.Fatal(err)
		}

		ctx.JSON(p)
	})

	payload := testRequest{
		Username: "test username",
		Data: []testData{
			{Name: "test data name 1"},
			{Name: "test data name 2"},
			{Name: "test data name 3"},
		},
		Other: 42,
	}

	e := httptest.New(t, app)
	e.POST("/").WithForm(payload).Expect().Status(iris.StatusOK).JSON().Equal(payload)
}
