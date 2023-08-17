package folders

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/MogLuiz/go-driver/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db: db}

	f := Folder{
		Name: "Folder 1",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&f)
	if err != nil {
		t.Error(err)
	}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", &b)

	mock.ExpectExec(regexp.QuoteMeta(`insert into "folders" ("parent_id", "name", "modified_at") values ($1, $2, $3)`)).
		WithArgs(f.ParentID, f.Name, utils.AnyTime{}).
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

	folder, err := New("Folder 1", 1)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`insert into "folders" ("parent_id", "name", "modified_at") values ($1, $2, $3)`)).
		WithArgs(folder.ParentID, folder.Name, utils.AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
