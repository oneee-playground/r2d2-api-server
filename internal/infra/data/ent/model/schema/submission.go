package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Submission holds the schema definition for the Submission entity.
type Submission struct {
	ent.Schema
}

// Fields of the Submission.
func (Submission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Unique(),
		field.Time("timestamp"),
		field.UUID("userID", uuid.New()),
		field.UUID("taskID", uuid.New()),
	}
}

// Edges of the Submission.
func (Submission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).Field("taskID").
			Ref("submissions").Unique().Required(),
		edge.From("user", User.Type).Field("userID").
			Ref("submissions").Unique().Required(),
		edge.To("events", Event.Type),
	}
}
