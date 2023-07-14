package api

import (
	"bar/autogen"
	"bar/internal/db"
)

type Server struct {
	db.DBackend
}

func NewServer(db db.DBackend) autogen.StrictServerInterface {
	return &Server{
		db,
	}
}
