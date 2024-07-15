package submission_module

import (
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

func toSubmissionListOutput(submissions []domain.Submission) *dto.SubmissionListOutput {
	out := make(dto.SubmissionListOutput, len(submissions))

	for i, submission := range submissions {
		url := "https://github.com/" + submission.Repository

		out[i] = dto.SubmissionListElem{
			ID:        submission.ID.String(),
			Timestamp: submission.Timestamp,
			IsDone:    submission.IsDone,
			SourceURL: url,
			User: dto.UserInfo{
				ID:         submission.User.ID.String(),
				Username:   submission.User.Username,
				ProfileURL: submission.User.ProfileURL,
				Role:       submission.User.Role.String(),
			},
		}
	}

	return &out
}

func toIDOutput(submission domain.Submission) *dto.IDOutput {
	return &dto.IDOutput{
		ID: submission.ID.String(),
	}
}
