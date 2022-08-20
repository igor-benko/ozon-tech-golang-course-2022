package postgres

import (
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	storage "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
)

type personRepoFixture struct {
	ctrl *gomock.Controller
	pool *pgxpoolmock.MockPgxPool
	repo storage.PersonRepo
}

func setUpPersonRepoFixture(t *testing.T) personRepoFixture {
	ctrl := gomock.NewController(t)
	pool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewPersonRepo(pool)

	return personRepoFixture{
		ctrl: ctrl,
		pool: pool,
		repo: repo,
	}
}

func (f *personRepoFixture) tearDownPersonRepoFixture() {
	f.ctrl.Finish()
}
