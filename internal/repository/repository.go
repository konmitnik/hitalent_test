package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
	ctx  context.Context
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{
		conn: conn,
		ctx:  context.Background(),
	}
}
