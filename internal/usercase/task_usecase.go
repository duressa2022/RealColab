package usecase

import (
	"context"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUseCase struct {
	TaskRepository *repository.TaskRepository
	Timeout        time.Duration
}

func NewTaskUseCase(timeout time.Duration, task *repository.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		TaskRepository: task,
		Timeout:        timeout,
	}
}

// method for adding new tasks
func (tu *TaskUseCase) AddTask(cxt context.Context, task *domain.TaskRequest, ID string) (*domain.Task, error) {
	userID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	taskID := primitive.NewObjectID()

	var createdTask domain.Task
	createdTask.Title = task.Title
	createdTask.TimePerDay = task.TimePerDay
	createdTask.Description = task.Description
	createdTask.StartDate = task.StartDate
	createdTask.DueDate = task.DueDate
	createdTask.Status = task.Status
	createdTask.Type = task.Type

	createdTask.UserID = userID
	createdTask.TaskID = taskID
	return tu.TaskRepository.PostTask(cxt, &createdTask)
}

// methof fot getting private tasks
func (tu *TaskUseCase) GetPrivateTasks(cxt context.Context, Id string, size int64, page int64) ([]*domain.PrivateTask, int64, error) {
	return tu.TaskRepository.GetPrivateTasks(cxt, Id, size, page)
}

// method for archiving the task

func (tu *TaskUseCase) ArchiveTask(cxt context.Context, ID string) error {
	return tu.TaskRepository.ArchiveTask(cxt, ID)
}

// method for editing the task
func (tu *TaskUseCase) EditTask(cxt context.Context, task *domain.EditTask, Id string) (*domain.Task, error) {
	return tu.TaskRepository.EditTask(cxt, task, Id)
}

// method for searching the task
func (tu *TaskUseCase) SearchTask(cxt context.Context, searchTerm string, size int64, page int64) ([]*domain.Task,int64, error) {
	return tu.TaskRepository.SearchTask(cxt, searchTerm, page, size)
}
