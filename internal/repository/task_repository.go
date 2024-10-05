package repository

import (
	"context"
	"errors"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository struct {
	database   mongo.Database
	collection string
}

func NewTaskRepository(collection string, database mongo.Database) *TaskRepository {
	return &TaskRepository{
		database:   database,
		collection: collection,
	}
}

// method for updating tasks
func (tr *TaskRepository) EditTask(cxt context.Context, task domain.EditTask, ID string) (*domain.Task, error) {
	taskcollection := tr.database.Collection(tr.collection)
	task_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	updatingValue := bson.M{
		"title":       task.Title,
		"timePerDay":  task.TimePerDay,
		"type":        task.Type,
		"description": task.Description,
	}

	updatedResult,err:=taskcollection.UpdateOne(cxt,bson.D{{Key: "_taskID",Value: task_id}},bson.D{{Key: "$set",Value: updatingValue}})
	if err!=nil{
		return nil,err
	}
	if updatedResult.MatchedCount==0{
		return nil,errors.New("no matched docs")
	}
	if updatedResult.ModifiedCount==0{
		return nil,errors.New("no modified docs")
	}
	
	var updatedTask *domain.Task
	err=taskcollection.FindOne(cxt,bson.D{{Key: "_taskID",Value: task_id}}).Decode(&updatedTask)
	if err!=nil{
		return nil,err
	}
	return updatedTask,err
}

// method for archiving tasks
func (tr *TaskRepository) ArchiveTask(cxt context.Context, taskID string) error {
	taskCollection := tr.database.Collection(tr.collection)
	task_id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	updatingValue := bson.M{
		"status": "Archived",
	}

	updatedResult, err := taskCollection.UpdateOne(cxt, bson.D{{Key: "_taskID", Value: task_id}}, bson.D{{Key: "$set", Value: updatingValue}})
	if err != nil {
		return err
	}
	if updatedResult.MatchedCount == 0 {
		return errors.New("no matching docs found")
	}
	if updatedResult.ModifiedCount == 0 {
		return errors.New("no modified docs found")
	}
	return nil
}

// method for creating new task
func (tr *TaskRepository) PostTask(cxt context.Context, task *domain.TaskInformation) (*domain.TaskInformation, error) {
	taskCollection := tr.database.Collection(tr.collection)
	_, err := taskCollection.InsertOne(cxt, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
