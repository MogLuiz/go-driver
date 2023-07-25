package users

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MogLuiz/go-driver/internal/utils"
	"github.com/go-chi/chi"
)

func TestModify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db: db}

	u := User{
		ID:   1,
		Name: "John Doe",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&u)
	if err != nil {
		t.Error(err)
	}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", strconv.Itoa(int(u.ID)))

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" set "name"=$1, "modified_at"=$2  where "id"=$3`)).
		WithArgs(u.Name, utils.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "John Doe", "john.doe@gmail.com", "123456", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	h.Modify(recordResponse, request)

	if recordResponse.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", recordResponse.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" set "name"=$1, "modified_at"=$2  where "id"=$3`)).
		WithArgs("Luiz", utils.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, 1, &User{Name: "Luiz"})
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
