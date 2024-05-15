package task_module

import (
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

func toIDOutput(task domain.Task) *dto.IDOutput {
	return &dto.IDOutput{ID: task.ID}
}

func toTaskListOutput(tasks []domain.Task) *dto.TaskListOutput {
	out := make(dto.TaskListOutput, len(tasks))

	for i, task := range tasks {
		out[i] = dto.TaskListElem{
			ID:    task.ID,
			Title: task.Title,
		}
	}

	return &out
}

func toTaskOutput(task domain.Task) *dto.TaskOutput {
	return &dto.TaskOutput{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Stage:       string(task.Stage),
	}
}
