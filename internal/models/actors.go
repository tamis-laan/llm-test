package models

import "github.com/pgvector/pgvector-go"

// Actor database model
type Actor struct {
	Id        uint            `db:"id"          json:"id"`
	Name      string          `db:"name"        json:"name"`
	Embedding pgvector.Vector `db:"embedding"   json:"-"    pg:"type:vector(1024)"`
}
