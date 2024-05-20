package validator

import "github.com/oneee-playground/r2d2-api-server/internal/domain"

func TaskStageValid(s domain.TaskStage) bool {
	switch s {
	case domain.StageDraft, domain.StageFixing, domain.StageAvailable:
		return true
	}
	return false
}

func SubmissionAcitonValid(a domain.SubmissionAction) bool {
	switch a {
	case domain.ActionApprove, domain.ActionReject:
		return true
	}
	return false
}

func SectionTypeValid(t domain.SectionType) bool {
	switch t {
	case domain.TypeScenario, domain.TypeLoad:
		return true
	}
	return false
}
