package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MogLuiz/go-driver/internal/utils"
	"github.com/go-chi/chi"
)

func TestDeleteHTTP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db: db}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/{id}", nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	mock.ExpectExec(`update "users" set *`).
		WithArgs(utils.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Delete(recordResponse, request)

	if recordResponse.Code != http.StatusNoContent {
		t.Errorf("Error: %v", recordResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

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
