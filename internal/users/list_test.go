package users

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db: db}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "John Doe", "john.doe@gmail.com", "123456", time.Now(), time.Now(), false, time.Now()).
		AddRow(2, "fulano", "fulano@gmail.com", "123456", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "deleted"=false`)).
		WillReturnRows(rows)

	h.List(recordResponse, request)

	if recordResponse.Code != http.StatusOK {
		t.Errorf("Error %v", recordResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSelectAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "John Doe", "john.doe@gmail.com", "123456", time.Now(), time.Now(), false, time.Now()).
		AddRow(2, "fulano", "fulano@gmail.com", "123456", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "deleted"=false`)).
		WillReturnRows(rows)

	_, err = SelectAll(db)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
