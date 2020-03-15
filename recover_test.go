package recover

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goroute/route"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	mux := route.NewServeMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := mux.NewContext(req, rec)
	options := GetDefaultOptions()
	mw := New(
		Skipper(options.Skipper),
		StackSize(options.StackSize),
		DisableStackAll(options.DisableStackAll),
		OnError(options.OnError),
	)
	h := func(c route.Context) error {
		panic("something went wrong")
	}

	err := mw(c, h)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.EqualError(t, err, "something went wrong")
}
