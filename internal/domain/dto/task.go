package dto

import "github.com/google/uuid"

type TaskListElem struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TaskListOutput []TaskListElem

type TaskOutput struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stage       string    `json:"stage"`
}

type TaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateTaskInput struct {
	IDInput
	TaskInput
}

type TaskStageInput struct {
	IDInput
	Stage string `json:"stage" binding:"required" validate:"task_stage"`
}
