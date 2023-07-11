package files

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type handler struct {
	db *sql.DB
}

func SetRoutes(r chi.Router, db *sql.DB) {
	h := &handler{db: db}

	r.Put("/{id}", h.Modify)
	r.Post("/", h.Create)
}
