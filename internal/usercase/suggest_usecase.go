package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"working/super_task/internal/repository"
)

type SuggestUseCase struct {
	TaskRepository *repository.TaskRepository
	Timeout        time.Duration
}

const (
	limit = 5
)

func NewSuggestUseCase(task *repository.TaskRepository, timeout time.Duration) *SuggestUseCase {
	return &SuggestUseCase{
		TaskRepository: task,
		Timeout:        timeout,
	}
}

// method for creating prompt for suggestions
func (sc *SuggestUseCase) CreatePrompt(cxt context.Context, Id string) (string, error) {
	compeleted, err := sc.TaskRepository.GetTaskByCriteria(cxt, "compeleted", Id, limit)
	if err != nil {
		return "", err
	}

	ongoing, err := sc.TaskRepository.GetTaskByCriteria(cxt, "ongoing", Id, limit)
	if err != nil {
		return "", err
	}

	var prevTasksDescription string = ""
	taskCounter := 1
	for _, data := range compeleted {
		prevTasksDescription += strconv.Itoa(taskCounter) + " ." + data + " \n"
		taskCounter = taskCounter + 1
	}
	for _, data := range ongoing {
		prevTasksDescription += strconv.Itoa(taskCounter) + " ." + data + " \n"
		taskCounter = taskCounter + 1
	}
	prompt := fmt.Sprintf(`
		Based on the following previous task descriptions:
		%s
		Please suggest two new task descriptions in this format
		and two suggestions has to be each of two and half sentences:
			Suggest1: [description]
			Suggest2: [description]
		`, prevTasksDescription)
	return prompt, nil
}
