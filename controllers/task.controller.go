package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manlikehenryy/go-gin-project/helpers"
	"github.com/manlikehenryy/go-gin-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		log.Println("Unable to parse body:", err)
		helpers.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userId, ok := c.MustGet("userId").(primitive.ObjectID)
	if !ok {
		helpers.SendError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	task.UserId = userId

	insertResult, err := tasksCollection.InsertOne(context.Background(), task)
	if err != nil {
		log.Println("Database error:", err)
		helpers.SendError(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	task.ID = insertResult.InsertedID.(primitive.ObjectID)

	helpers.SendJSON(c, http.StatusCreated, gin.H{
		"data":    task,
		"message": "Task created successfully",
	})
}

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task
	filter := bson.M{}

	params, err := helpers.PaginateCollection(c, tasksCollection, filter, &tasks)
	if err != nil {
		helpers.SendError(c, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	helpers.SendPaginatedResponse(c, tasks, params)
}

func GetTask(c *gin.Context) {

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.SendError(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	userId, ok := c.MustGet("userId").(primitive.ObjectID)
	if !ok {
		helpers.SendError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	var task models.Task
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": id, "userId": userId}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helpers.SendError(c, http.StatusNotFound, "Task not found")
		} else {
			helpers.SendError(c, http.StatusInternalServerError, "Failed to retrieve task")
		}
		return
	}

	helpers.SendJSON(c, http.StatusOK, gin.H{
		"data": task,
	})
}

func UpdateTask(c *gin.Context) {

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.SendError(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	userId, ok := c.MustGet("userId").(primitive.ObjectID)
	if !ok {
		helpers.SendError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		helpers.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var existingTask models.Task
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helpers.SendError(c, http.StatusNotFound, "Task not found")
		} else {
			helpers.SendError(c, http.StatusInternalServerError, "Failed to retrieve task")
		}
		return
	}
	

	if existingTask.UserId != userId {
		helpers.SendError(c, http.StatusForbidden, "Unauthorized to update this task")
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":  task.Title,
			"desc":   task.Desc,
			"status": task.Status,
		},
	}

	result, err := tasksCollection.UpdateOne(context.Background(), bson.M{"_id": id, "userId": userId}, update)
	if err != nil {
		log.Println("Database error:", err)
		helpers.SendError(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	if result.MatchedCount == 0 {
		helpers.SendError(c, http.StatusNotFound, "Task not found or unauthorized")
		return
	}

	helpers.SendJSON(c, http.StatusOK, gin.H{
		"message": "Task updated successfully",
	})
}

func UsersTask(c *gin.Context) {
	userId, ok := c.MustGet("userId").(primitive.ObjectID)
	if !ok {
		helpers.SendError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	filter := bson.M{"userId": userId}

	// Define a slice to hold the tasks
	var tasks []models.Task
	params, err := helpers.PaginateCollection(c, tasksCollection, filter, &tasks)
	if err != nil {
		helpers.SendError(c, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	helpers.SendPaginatedResponse(c, tasks, params)
}

func DeleteTask(c *gin.Context) {

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.SendError(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	userId, ok := c.MustGet("userId").(primitive.ObjectID)
	if !ok {
		helpers.SendError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	var existingTask models.Task
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helpers.SendError(c, http.StatusNotFound, "Task not found")
		} else {
			helpers.SendError(c, http.StatusInternalServerError, "Failed to retrieve task")
		}
		return
	}

	if existingTask.UserId != userId {
		helpers.SendError(c, http.StatusForbidden, "Unauthorized to delete this task")
		return
	}

	filter := bson.M{"_id": id, "userId": userId}
	result, err := tasksCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		helpers.SendError(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	if result.DeletedCount == 0 {
		helpers.SendError(c, http.StatusNotFound, "Task not found or unauthorized")
		return
	}

	helpers.SendJSON(c, http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
