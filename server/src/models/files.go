package models

import "database/sql"

type SQLFiles struct {
	Id         int64          `json:"id,omitempty"`
	Request_id int64          `json:"request_id,omitempty"`
	Time       sql.NullString `json:"time"`
}

type Files struct {
	Request_id int64  `json:"request_id,omitempty"`
	Time       string `json:"time"`
}
