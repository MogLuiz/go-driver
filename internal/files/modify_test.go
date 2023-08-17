package files

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

	f := File{
		ID:   1,
		Name: "image.png",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&f)
	if err != nil {
		t.Error(err)
	}

	recordResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", strconv.Itoa(int(f.ID)))

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"}).
		AddRow(1, 1, 1, "image.png", "image/png", "/tmp/image.png", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "files" where "id"=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" set "name"=$1, "modified_at"=$2, "deleted"=$3  where "id"=$4`)).
		WithArgs(f.Name, utils.AnyTime{}, false, f.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Modify(recordResponse, request)

	if recordResponse.Code != http.StatusOK {
		t.Errorf("Error: %v", recordResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
