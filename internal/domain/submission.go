package domain

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID        uuid.UUID
	Timestamp time.Time

	UserID uuid.UUID
	User   *User

	TaskID uuid.UUID
	Task   *Task
}
