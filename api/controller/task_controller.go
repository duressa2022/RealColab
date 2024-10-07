package controller

import (
	"net/http"
	"strconv"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUseCase *usecase.TaskUseCase
	Env         *config.Env
}

func NewTaskController(env *config.Env, task *usecase.TaskUseCase) *TaskController {
	return &TaskController{
		TaskUseCase: task,
		Env:         env,
	}
}

// handler for working with new tasks
func (tc *TaskController) AddTaskHandler(c *gin.Context) {
	var task *domain.TaskRequest
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	ID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	usrID, ok := ID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data type", "success": false, "data": nil})
		return
	}

	createdTask, err := tc.TaskUseCase.AddTask(c, task, usrID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "new task is ceated",
		"success": true,
		"data":    createdTask,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for working with private tasks
func (tc *TaskController) GetTaskByTypeHandler(c *gin.Context) {
	size := c.Query("size")
	page := c.Query("page")

	sizeNumber, err := strconv.Atoi(size)
	if err != nil || sizeNumber < 1 {
		sizeNumber = 4
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	ID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	userID, ok := ID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data type", "success": false, "data": nil})
		return
	}

	privateTasks, total, err := tc.TaskUseCase.GetPrivateTasks(c, userID, int64(sizeNumber), int64(pageNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	totalPage := (total + int64(sizeNumber) - 1) / int64(sizeNumber)

	response := map[string]interface{}{
		"message": "successful",
		"success": true,
		"content": map[string]interface{}{
			"data":            privateTasks,
			"totalPageNumber": totalPage,
		},
	}

	c.IndentedJSON(http.StatusOK, response)

}

// handler for archiving the tasks by using id
func (tc *TaskController) ArchiveTaskHandler(c *gin.Context) {
	taskID := c.Param("taskID")

	err := tc.TaskUseCase.ArchiveTask(c, taskID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "archived",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for working with the task edit
func (tc *TaskController) EditTaskHandler(c *gin.Context) {
	var newTask *domain.EditTask
	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	taskID := c.Param("taskID")

	updatedTask, err := tc.TaskUseCase.EditTask(c, newTask, taskID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "updated",
		"success": true,
		"data":    updatedTask,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for working with the search tasks
func (tc *TaskController) SearchTaskHandler(c *gin.Context) {
	var searchTerm *domain.SearchTerm
	if err := c.BindJSON(&searchTerm); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	page := c.Query("page")
	size := c.Query("size")

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}
	sizeNumber, err := strconv.Atoi(size)
	if err != nil {
		sizeNumber = 4
	}

	searchResult, totalResult, err := tc.TaskUseCase.SearchTask(c, searchTerm.SearchTerm, int64(sizeNumber), int64(pageNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	totalPages := (totalResult + int64(sizeNumber) - 1) / int64(sizeNumber)

	response := map[string]interface{}{
		"success": true,
		"message": "result",
		"content": map[string]interface{}{
			"data":       searchResult,
			"totalPages": totalPages,
		},
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for working with the archived tasks
func (tc *TaskController) GetArchivedTasksHandler(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("Size")
	ID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	userID, ok := ID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data type", "success": false, "data": nil})
		return
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}
	sizeNumber, err := strconv.Atoi(size)
	if err != nil || sizeNumber < 1 {
		sizeNumber = 1
	}

	archived, numberTask, err := tc.TaskUseCase.GetArchivedTasks(c, userID, int64(pageNumber), int64(sizeNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": true, "gin": nil})
		return
	}
	totalPage := (numberTask + int64(sizeNumber) - 1) / int64(sizeNumber)

	response := map[string]interface{}{
		"message": "archived tasks",
		"success": true,
		"content": map[string]interface{}{
			"data":       archived,
			"totalPages": totalPage,
		},
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for working with restore the archived tasks
func (tc *TaskController) RestoreArchived(c *gin.Context) {
	taskID := c.Param("taskID")

	err := tc.TaskUseCase.RestoreArchived(c, taskID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	response := map[string]interface{}{
		"message": "restored",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for working with the deleting the tasks
func (tc *TaskController) DeleteArchived(c *gin.Context) {
	taskID := c.Param("taskID")

	err := tc.TaskUseCase.DeleteArchived(c, taskID)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "deleted",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)

}
