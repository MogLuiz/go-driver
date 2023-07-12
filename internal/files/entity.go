package files

import (
	"errors"
	"time"
)

var (
	ErrNameRequired    = errors.New("name is required and can't be empty")
	ErrIDRequired      = errors.New("id is required and can't be empty")
	ErrTypeRequired    = errors.New("type is required and can't be empty")
	ErrPathRequired    = errors.New("path is required and can't be empty")
	ErrOwnerIDRequired = errors.New("owner_id is required and can't be empty")
	ErrFileNotFound    = errors.New("file not found")
)

type File struct {
	ID         int64     `json:"id"`
	FolderID   int64     `json:"-"`
	OwnerID    int64     `json:"owner_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Path       string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
}

func New(ownerID int64, name, fileType, path string) (*File, error) {
	file := File{
		OwnerID: ownerID,
		Name:    name,
		Type:    fileType,
		Path:    path,
	}

	err := file.Validate()
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (file *File) Validate() error {
	if file.Name == "" {
		return ErrNameRequired
	}

	if file.Type == "" {
		return ErrTypeRequired
	}

	if file.Path == "" {
		return ErrPathRequired
	}

	if file.OwnerID == 0 {
		return ErrOwnerIDRequired
	}

	return nil
}
