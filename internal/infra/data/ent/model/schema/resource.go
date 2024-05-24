package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Resource holds the schema definition for the Resource entity.
// NOTE: cannot set composite primary key due to lack of functionality of entgo.
// You should find workaround later.
type Resource struct {
	ent.Schema
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("image"),
		field.Uint16("port"),
		field.Float("cpu"),
		field.Uint64("memory"),
		field.Bool("isPrimary"),
		field.UUID("taskID", uuid.New()),
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).Field("taskID").
			Ref("resources").Unique().Required(),
	}
}
