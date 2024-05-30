package repository

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/submission"
)

type SubmissionRepository struct {
	client *model.SubmissionClient
}

var _ domain.SubmissionRepository = (*SubmissionRepository)(nil)

func NewSubmissionRepository(client *model.SubmissionClient) *SubmissionRepository {
	return &SubmissionRepository{client: client}
}

func (r *SubmissionRepository) Create(ctx context.Context, submission domain.Submission) error {
	return r.client.Create().
		SetID(submission.ID).
		SetIsDone(submission.IsDone).
		SetRepository(submission.Repository).
		SetCommitHash(submission.CommitHash).
		SetTimestamp(submission.Timestamp).
		SetTaskID(submission.TaskID).
		SetUserID(submission.UserID).
		Exec(ctx)
}

func (r *SubmissionRepository) FetchByID(ctx context.Context, id uuid.UUID) (domain.Submission, error) {
	entity, err := r.client.Get(ctx, id)
	if err != nil {
		if model.IsNotFound(err) {
			return domain.Submission{}, domain.ErrSubmissionNotFound
		}
		return domain.Submission{}, err
	}

	submission := domain.Submission{
		ID:         entity.ID,
		Timestamp:  entity.Timestamp,
		IsDone:     entity.IsDone,
		Repository: entity.Repository,
		CommitHash: entity.CommitHash,
		TaskID:     entity.TaskID,
		UserID:     entity.UserID,
	}

	return submission, nil
}

func (r *SubmissionRepository) UndoneExists(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) (bool, error) {
	return r.client.Query().Where(
		submission.And(
			submission.TaskID(taskID),
			submission.UserID(userID),
			submission.IsDone(false),
		),
	).Exist(ctx)
}

func (r *SubmissionRepository) Update(ctx context.Context, submission domain.Submission) error {
	return r.client.UpdateOneID(submission.ID).
		SetIsDone(submission.IsDone).
		SetRepository(submission.Repository).
		SetCommitHash(submission.CommitHash).
		SetTimestamp(submission.Timestamp).
		Exec(ctx)
}

func (r *SubmissionRepository) FetchPaginated(ctx context.Context, taskID uuid.UUID, offset int, limit int) ([]domain.Submission, error) {
	models, err := r.client.Query().
		Where(submission.TaskID(taskID)).
		WithUser().
		Order(submission.ByTimestamp(sql.OrderDesc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	submissions := make([]domain.Submission, len(models))
	for idx, model := range models {
		submissions[idx] = domain.Submission{
			ID:         model.ID,
			Timestamp:  model.Timestamp,
			IsDone:     model.IsDone,
			Repository: model.Repository,
			CommitHash: model.CommitHash,
			TaskID:     model.TaskID,
			UserID:     model.UserID,
			User: &domain.User{
				ID:         model.Edges.User.ID,
				Username:   model.Edges.User.Username,
				Email:      model.Edges.User.Email,
				ProfileURL: model.Edges.User.ProfileURL,
				Role:       domain.UserRole(model.Edges.User.Role),
			},
		}
	}

	return submissions, nil
}
