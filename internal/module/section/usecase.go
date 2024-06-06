package section_module

import (
	"context"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/pkg/errors"
)

type sectionUsecase struct {
	lock tx.Locker

	taskRepository    domain.TaskRepository
	sectionRepository domain.SectionRepository
}

var _ domain.SectionUsecase = (*sectionUsecase)(nil)

func NewSectionUsecase(sr domain.SectionRepository, tr domain.TaskRepository, l tx.Locker) *sectionUsecase {
	return &sectionUsecase{
		sectionRepository: sr,
		taskRepository:    tr,
		lock:              l,
	}
}

func (s *sectionUsecase) GetList(ctx context.Context, in dto.IDInput) (out *dto.SectionListOutput, err error) {
	if err := s.assureTaskExists(ctx, in.ID); err != nil {
		return nil, err
	}

	fetchOpts := domain.FetchSectionsOption{
		IncludeContent: true,
	}

	sections, err := s.sectionRepository.FetchAllByTaskID(ctx, in.ID, fetchOpts)
	if err != nil {
		return nil, errors.Wrap(err, "fetching sections")
	}

	return toSectionListOutput(sections), nil
}

func (s *sectionUsecase) CreateSection(ctx context.Context, in dto.CreateSectionInput) (err error) {
	ctx = tx.NewAtomic(ctx)
	defer tx.Evaluate(ctx, &err)

	ctx, release, err := s.lock.Acquire(ctx, "task", in.ID.String())
	if err != nil {
		return errors.Wrap(err, "acquiring lock")
	}
	defer release()

	taskID := in.ID

	if err := s.assureTaskExists(ctx, taskID); err != nil {
		return err
	}

	count, err := s.sectionRepository.CountByTaskID(ctx, taskID)
	if err != nil {
		return errors.Wrap(err, "counting sections")
	}

	section := domain.Section{
		ID:          uuid.New(),
		Title:       in.Title,
		Description: in.Description,
		Type:        domain.SectionType(in.Type),
		TaskID:      taskID,
		Index:       count,
	}

	if err := s.sectionRepository.Create(ctx, section); err != nil {
		return errors.Wrap(err, "creating section")
	}

	return nil
}

func (s *sectionUsecase) UpdateSection(ctx context.Context, in dto.UpdateSectionInput) (err error) {
	ctx = tx.NewAtomic(ctx)
	defer tx.Evaluate(ctx, &err)

	if err := s.assureTaskExists(ctx, in.TaskID); err != nil {
		return err
	}

	section, err := s.sectionRepository.FetchByID(ctx, in.SectionID)
	if err != nil {
		if errors.Is(err, domain.ErrSectionNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return errors.Wrap(err, "fetching section")
	}

	section.Title = in.Title
	section.Description = in.Description

	// Ignore section type for now.
	// section.Type = domain.SectionType(in.Type)

	if err = s.sectionRepository.Update(ctx, section); err != nil {
		return errors.Wrap(err, "updating section")
	}

	return nil
}

func (s *sectionUsecase) ChangeIndex(ctx context.Context, in dto.SectionIndexInput) (err error) {
	ctx = tx.NewAtomic(ctx)
	defer tx.Evaluate(ctx, &err)

	ctx, release, err := s.lock.Acquire(ctx, "task", in.TaskID.String())
	if err != nil {
		return errors.Wrap(err, "acquiring lock")
	}
	defer release()

	if err := s.assureTaskExists(ctx, in.TaskID); err != nil {
		return err
	}

	fetchOpts := domain.FetchSectionsOption{
		IncludeContent: false,
	}

	sections, err := s.sectionRepository.FetchAllByTaskID(ctx, in.TaskID, fetchOpts)
	if err != nil {
		return errors.Wrap(err, "fetching sections")
	}

	idx := in.Index

	isOutOfRange := idx < 0 || idx >= len(sections)
	if isOutOfRange {
		return status.NewErr(http.StatusForbidden, "index out of range")
	}

	idxFunc := func(section domain.Section) bool {
		return section.ID == in.SectionID
	}

	curIdx := slices.IndexFunc(sections, idxFunc)
	if curIdx == -1 {
		return status.NewErr(http.StatusNotFound, "section not found")
	}

	changeIndex(sections, curIdx, idx)

	if err := s.sectionRepository.SaveIndexes(ctx, sections); err != nil {
		return errors.Wrap(err, "saving indexes")
	}

	return nil
}

// assureTaskExists checks if task exists. if not exists, it will return an error.
func (u *sectionUsecase) assureTaskExists(ctx context.Context, taskID uuid.UUID) error {
	exists, err := u.taskRepository.ExistsByID(ctx, taskID)
	if err != nil {
		return errors.Wrap(err, "checking task exists")
	}

	if !exists {
		return status.NewErr(http.StatusNotFound, "task not found")
	}

	return nil
}
