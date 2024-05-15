package validator

import "github.com/oneee-playground/r2d2-api-server/internal/domain"

func TaskStageValid(s domain.TaskStage) bool {
	switch s {
	case domain.StageDraft, domain.StageFixing, domain.StageAvailable:
		return true
	}
	return false
}
