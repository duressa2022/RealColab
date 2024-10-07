package usecase

import (
	"context"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"
)

type HomeUseCase struct {
	UserRepository *repository.UserRepository
	TaskRepository *repository.TaskRepository
	Timeout        time.Duration
}

func NewHomeUseCase(user *repository.UserRepository, task *repository.TaskRepository, time time.Duration) *HomeUseCase {
	return &HomeUseCase{
		UserRepository: user,
		TaskRepository: task,
		Timeout:        time,
	}
}

func (hc *HomeUseCase) HomeInformation(cxt context.Context, Id string) (*domain.HomeInformation, error) {
	userInformation, err := hc.UserRepository.GetUserByID(cxt, Id)
	if err != nil {
		return nil, err
	}

	upcomingTask, err := hc.TaskRepository.GetUpComingTasks(cxt, Id)
	if err != nil {
		return nil, err
	}

	recentlyCompleted, err := hc.TaskRepository.GetRecentlyCompletedTasks(cxt, Id)
	if err != nil {
		return nil, err
	}

	var userHomeInformation domain.HomeInformation

	userHomeInformation.FirstName = userInformation.FirstName
	userHomeInformation.LastName = userInformation.LastName
	userHomeInformation.ProfileUrl = userInformation.ProfileUrl
	userHomeInformation.Rating = userInformation.Rating
	userHomeInformation.TaskData = userInformation.TaskInformation

	userHomeInformation.UpcomingTasks = append(userHomeInformation.UpcomingTasks, upcomingTask...)
	userHomeInformation.RecentlyCompletedTasks = append(userHomeInformation.RecentlyCompletedTasks, recentlyCompleted...)
	return &userHomeInformation, nil

}
