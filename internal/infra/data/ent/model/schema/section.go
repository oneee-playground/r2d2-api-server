package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Section holds the schema definition for the Section entity.
type Section struct {
	ent.Schema
}

// Fields of the Section.
func (Section) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Unique(),
		field.String("title"),
		field.String("description"),
		field.Uint8("index"),
		field.Uint64("rpm"),
		field.String("type"),
		field.String("example"),
		field.UUID("taskID", uuid.New()),
	}
}

// Edges of the Section.
func (Section) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).Field("taskID").
			Ref("sections").Unique().Required(),
	}
}
