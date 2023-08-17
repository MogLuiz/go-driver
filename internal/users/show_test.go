package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db: db}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/{id}", nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "John Doe", "john.doe@gmail.com", "123456", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	h.GetByID(recordResponse, request)

	if recordResponse.Code != http.StatusOK {
		t.Errorf("Error: %v", recordResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestShowByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "John Doe", "john.doe@gmail.com", "123456", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = ShowByID(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
