package repository

import (
	"context"
	"errors"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// method for counting all tasks of given id
func (tr *TaskRepository) TaskInformation(cxt context.Context, ID string) (map[string]int64, error) {
	taskCollection := tr.database.Collection(tr.collection)
	taskInformation := map[string]int64{}

	userId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return taskInformation, err
	}

	allTasks, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userId}})
	if err != nil {
		return taskInformation, err
	}
	taskInformation["allTasks"] = allTasks
	completedTasks, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userId}, {Key: "status", Value: "compeleted"}})
	if err != nil {
		return taskInformation, err
	}
	taskInformation["compeletedTasks"] = completedTasks

	onProgressTasks, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userId}, {Key: "status", Value: "onProgress"}})
	if err != nil {
		return taskInformation, err
	}
	taskInformation["onProgress"] = onProgressTasks

	expiredTasks, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userId}, {Key: "status", Value: "expired"}})
	if err != nil {
		return taskInformation, err
	}
	taskInformation["expiredTasks"] = expiredTasks

	archivedTasks, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userId}, {Key: "status", Value: "archived"}})
	if err != nil {
		return taskInformation, err
	}
	taskInformation["archived"] = archivedTasks

	return taskInformation, nil

}

// method for getting recently compeleted tasks
func (tr *TaskRepository) GetRecentlyCompletedTasks(cxt context.Context, ID string) ([]*domain.Task, error) {
	var recentlyCompeltedTasks []*domain.Task
	taskCollection := tr.database.Collection(tr.collection)

	taskid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	cursor, err := taskCollection.Find(cxt, bson.D{{Key: "status", Value: "compeleted"}, {Key: "taskID", Value: taskid}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var task *domain.Task
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		recentlyCompeltedTasks = append(recentlyCompeltedTasks, task)
	}
	return recentlyCompeltedTasks, err
}

// method for getting upcoming tasks by using id
func (tr *TaskRepository) GetUpComingTasks(cxt context.Context, Id string) ([]*domain.PrivateTask, error) {
	var upComing []*domain.PrivateTask
	taskCollection := tr.database.Collection(tr.collection)
	opts := options.Find().SetSort(bson.D{{Key: "dueDate", Value: 1}})

	taskid, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return nil, err
	}
	cursor, err := taskCollection.Find(cxt, bson.D{{Key: "status", Value: "ongoing"}, {Key: "taskID", Value: taskid}}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var task *domain.PrivateTask
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		upComing = append(upComing, task)
	}
	return upComing, nil

}

// method for searching tasks by using title
func (tr *TaskRepository) SearchArchivedTasks(cxt context.Context, title string) (*domain.Task, error) {
	taskCollection := tr.database.Collection(tr.collection)
	var task *domain.Task
	err := taskCollection.FindOne(cxt, bson.D{{Key: "title", Value: title}, {Key: "status", Value: "archived"}}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil

}

// method for deleting archived tasks by using id
func (tr *TaskRepository) DeleteArchived(cxt context.Context, ID string) error {
	taskCollection := tr.database.Collection(tr.collection)

	taskid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	_, err = taskCollection.DeleteOne(cxt, bson.D{{Key: "taskID", Value: taskid}})
	if err != nil {
		return err
	}
	return nil
}

// method for restoring archived tasks by using id
func (tr *TaskRepository) RestoreArchived(cxt context.Context, ID string) error {
	taskCollection := tr.database.Collection(tr.collection)

	taskid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	updatingValue := bson.M{
		"status": "ongoing",
	}
	updatedResult, err := taskCollection.UpdateOne(cxt, bson.D{{Key: "_taskID", Value: taskid}}, bson.D{{Key: "$set", Value: updatingValue}})
	if err != nil {
		return err
	}
	if updatedResult.MatchedCount == 0 {
		return errors.New("no matched docs found")
	}
	if updatedResult.ModifiedCount == 0 {
		return errors.New("no modified docs found")
	}
	return nil
}

// method for getting archived tasks by using sie and page
func (tr *TaskRepository) GetArchivedTasks(cxt context.Context, ID string, size int64, page int64) ([]*domain.Task, int64, error) {
	var archivedTasks []*domain.Task

	taskCollection := tr.database.Collection(tr.collection)
	skip := (page - 1) / size
	opts := options.Find().SetSkip(skip).SetLimit(size)

	userID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := taskCollection.Find(cxt, bson.D{{Key: "status", Value: "archived"}, {Key: "_userID", Value: userID}}, opts)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(cxt) {
		var archivedTask *domain.Task
		err := cursor.Decode(&archivedTask)
		if err != nil {
			return nil, 0, err
		}
		archivedTasks = append(archivedTasks, archivedTask)
	}

	numberDocs, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userID}, {Key: "status", Value: "archived"}})
	if err != nil {
		return nil, 0, err
	}
	return archivedTasks, numberDocs, nil

}

// method for getting all shared tasks based size and page
func (tr *TaskRepository) GetSharedTasks(cxt context.Context, ID string, size int64, page int64) ([]*domain.SharedTask, int64, error) {
	var sharedTasks []*domain.SharedTask

	taskCollection := tr.database.Collection(tr.collection)
	skip := (page - 1) / size
	opts := options.Find().SetSkip(skip).SetLimit(size)

	userId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{"members": userId}
	cursor, err := taskCollection.Find(cxt, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(cxt) {
		var sharedTask *domain.SharedTask
		err := cursor.Decode(&sharedTask)
		if err != nil {
			return nil, 0, err
		}
		sharedTasks = append(sharedTasks, sharedTask)
	}

	numberDocs, err := taskCollection.CountDocuments(cxt, filter)
	if err != nil {
		return nil, 0, err
	}
	return sharedTasks, numberDocs, nil

}

// method for getting all private tasks  based size and page
func (tr *TaskRepository) GetPrivateTasks(cxt context.Context, ID string, size int64, page int64) ([]*domain.PrivateTask, int64, error) {
	var privateTask []*domain.PrivateTask

	taskCollection := tr.database.Collection(tr.collection)
	skip := (page - 1) / size
	opts := options.Find().SetSkip(skip).SetLimit(size)

	userId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, 0, err
	}

	curser, err := taskCollection.Find(cxt, bson.D{{Key: "_userID", Value: userId}}, opts)
	if err != nil {
		return nil, 0, err
	}

	for curser.Next(cxt) {
		var task *domain.PrivateTask
		err := curser.Decode(&task)
		if err != nil {
			return nil, 0, err
		}
		privateTask = append(privateTask, task)
	}

	numberOfDocs, err := taskCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userId}})
	if err != nil {
		return nil, 0, err
	}
	return privateTask, numberOfDocs, nil
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

	updatedResult, err := taskcollection.UpdateOne(cxt, bson.D{{Key: "_taskID", Value: task_id}}, bson.D{{Key: "$set", Value: updatingValue}})
	if err != nil {
		return nil, err
	}
	if updatedResult.MatchedCount == 0 {
		return nil, errors.New("no matched docs")
	}
	if updatedResult.ModifiedCount == 0 {
		return nil, errors.New("no modified docs")
	}

	var updatedTask *domain.Task
	err = taskcollection.FindOne(cxt, bson.D{{Key: "_taskID", Value: task_id}}).Decode(&updatedTask)
	if err != nil {
		return nil, err
	}
	return updatedTask, err
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
