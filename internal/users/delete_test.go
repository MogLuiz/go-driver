package users

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MogLuiz/go-driver/internal/utils"
)

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(`update "users" set *`).
		WithArgs(utils.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Delete(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
