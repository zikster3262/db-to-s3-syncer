package models

import (
	"database/sql"
)

type SQLRequest struct {
	Id   int64          `json:"id,omitempty"`
	Uuid string         `json:"uuid"`
	Time sql.NullString `json:"time"`
}

type DbRequest struct {
	Id   int64          `db:"id"`
	Uuid string         `db:"uuid"`
	Time sql.NullString `db:"time"`
}
