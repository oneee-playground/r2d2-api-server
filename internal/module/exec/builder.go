package exec_module

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=builder.go -destination=../../../test/mocks/builder.go -package=mocks

type BuildOpts struct {
	ID         uuid.UUID `json:"id"`
	TaskID     uuid.UUID `json:"taskID"`
	Repository string    `json:"repository"`
	CommitHash string    `json:"commitHash"`
	Platform   string    `json:"platform"`
}

type ImageBuilder interface {
	// RequestBuild requests actual builder to build.
	// Error will be nil if build has started without error.
	// Result of build should be informed with eventbus.
	RequestBuild(ctx context.Context, opts BuildOpts) error
}
