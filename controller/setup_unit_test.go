package controller

import (
	"testing"

	"URL_SHORTENER/storage"

	"github.com/golang/mock/gomock"
)

type Resources struct {
	ctl    *gomock.Controller
	MockDb *storage.MockURLOperations
}

func SetupTestDB(t *testing.T) *Resources {
	t.Helper()
	r := new(Resources)
	r.ctl = gomock.NewController(t)
	r.MockDb = storage.NewMockURLOperations(r.ctl)
	Init(r.MockDb)
	return r
}

// teardown test resources.
func (r *Resources) TearDown() {
	r.ctl.Finish()
}
