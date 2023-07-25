package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db: db}

	u := User{
		Name:     "John Doe",
		Login:    "john.doe@golang.com",
		Password: "12345",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&u)
	if err != nil {
		t.Error(err)
	}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", &b)

	mock.ExpectExec(regexp.QuoteMeta(`insert into "users" ("name", "login", "password", "modified_at") values ($1, $2, $3, $4)`)).
		WithArgs(u.Name, u.Login, u.Password, u.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Create(recordResponse, request)

	if recordResponse.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", recordResponse.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	user, err := New("John Doe", "john.doe", "123456")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`insert into "users" ("name", "login", "password", "modified_at") values ($1, $2, $3, $4)`)).
		WithArgs(user.Name, user.Login, user.Password, user.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, user)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
