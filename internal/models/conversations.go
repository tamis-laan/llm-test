package models

import "github.com/pgvector/pgvector-go"

// Report database model
type Conversation struct {
	Id          uint            `db:"id"          json:"id"`
	Title       string          `db:"title"       json:"title"`
	Description string          `db:"description" json:"description"`
	Embedding   pgvector.Vector `db:"embedding"   json:"-"           pg:"type:vector(1024)"`
}
