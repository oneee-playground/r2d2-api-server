package repository

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/section"
)

type SectionRepository struct {
	client *model.SectionClient
}

var _ domain.SectionRepository = (*SectionRepository)(nil)

func NewSectionRepository(client *model.SectionClient) *SectionRepository {
	return &SectionRepository{client: client}
}

func (r *SectionRepository) CountByTaskID(ctx context.Context, taskID uuid.UUID) (uint8, error) {
	count, err := r.client.Query().Where(
		section.TaskID(taskID),
	).Count(ctx)
	if err != nil {
		return 0, err
	}

	return uint8(count), nil
}

func (r *SectionRepository) Create(ctx context.Context, section domain.Section) error {
	return r.client.Create().
		SetID(section.ID).
		SetDescription(section.Description).
		SetIndex(section.Index).
		SetTitle(section.Title).
		SetType(string(section.Type)).
		SetTaskID(section.TaskID).
		Exec(ctx)
}

func (r *SectionRepository) FetchByID(ctx context.Context, id uuid.UUID) (domain.Section, error) {
	entity, err := r.client.Get(ctx, id)
	if err != nil {
		if model.IsNotFound(err) {
			return domain.Section{}, domain.ErrSectionNotFound
		}
		return domain.Section{}, err
	}

	section := domain.Section{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		Index:       entity.Index,
		Type:        domain.SectionType(entity.Type),
		TaskID:      entity.TaskID,
	}

	return section, nil
}

func (r *SectionRepository) Update(ctx context.Context, section domain.Section) error {
	return r.client.UpdateOneID(section.ID).
		SetDescription(section.Description).
		SetIndex(section.Index).
		SetTitle(section.Title).
		SetType(string(section.Type)).
		SetTaskID(section.TaskID).
		Exec(ctx)
}

// TODO: Optimize this. Using temporary table for this should be suitable.
// Reference: https://gngsn.tistory.com/189
// I don't think entgo supports this at the moment.
func (r *SectionRepository) SaveIndexes(ctx context.Context, sections []domain.Section) error {
	for _, section := range sections {
		err := r.client.UpdateOneID(section.ID).SetIndex(section.Index).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SectionRepository) FetchAllByTaskID(ctx context.Context, taskID uuid.UUID, opt domain.FetchSectionsOption) ([]domain.Section, error) {
	builder := r.client.Query().
		Where(section.TaskID(taskID)).
		Order(section.ByIndex(sql.OrderAsc()))

	var models []*model.Section
	var err error
	if opt.IncludeContent {
		models, err = builder.All(ctx)
	} else {
		// Don't include title and description
		models, err = builder.Select(
			section.FieldID, section.FieldIndex,
			section.FieldType, section.FieldTaskID,
		).All(ctx)
	}

	if err != nil {
		return nil, err
	}

	sections := make([]domain.Section, len(models))
	for idx, model := range sections {
		sections[idx] = domain.Section{
			ID:          model.ID,
			Title:       model.Title,
			Description: model.Description,
			Index:       model.Index,
			Type:        model.Type,
			TaskID:      model.TaskID,
		}
	}

	return sections, nil
}
