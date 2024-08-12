package helpers

import (
	"context"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/manlikehenryy/go-gin-project/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationParams struct {
	Page        int
	Limit       int
	Total       int64
	PageCount   int
	HasNextPage bool
	HasPrevPage bool
}

// func PaginateCollection(
// 	c *gin.Context,
// 	collection *mongo.Collection,
// 	filter interface{},
// 	result interface{}, // Expects a slice of target type, e.g., *[]models.Task
// ) (*PaginationParams, error) {
// 	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	if err != nil || page < 1 {
// 		page = 1
// 		log.Printf("Invalid page number. Setting to default: %d", page)
// 	}

// 	limit, err := strconv.Atoi(c.DefaultQuery("perPage", "10"))
// 	if err != nil || limit < 1 {
// 		limit = 10
// 		log.Printf("Invalid perPage value. Setting to default: %d", limit)
// 	}

// 	offset := (page - 1) * limit

// 	total, err := collection.CountDocuments(context.Background(), filter)
// 	if err != nil {
// 		log.Printf("Error counting documents: %v", err)
// 		return nil, err
// 	}

// 	cursor, err := collection.Find(
// 		context.Background(),
// 		filter,
// 		options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)),
// 	)
// 	if err != nil {
// 		log.Printf("Error finding documents: %v", err)
// 		return nil, err
// 	}
// 	defer func() {
// 		if err := cursor.Close(context.Background()); err != nil {
// 			log.Printf("Error closing cursor: %v", err)
// 		}
// 	}()

// 	// Use reflection to decode the cursor into the target type slice
// 	if resultSlice, ok := result.(*[]interface{}); ok {
// 		for cursor.Next(context.Background()) {
// 			var item interface{}
// 			if err := cursor.Decode(&item); err != nil {
// 				log.Printf("Error decoding document: %v", err)
// 				return nil, err
// 			}
// 			*resultSlice = append(*resultSlice, item)
// 		}
// 	} else {
// 		// Decode directly into the result slice type
// 		for cursor.Next(context.Background()) {
// 			if err := cursor.Decode(result); err != nil {
// 				log.Printf("Error decoding document: %v", err)
// 				return nil, err
// 			}
// 		}
// 	}

// 	if err := cursor.Err(); err != nil {
// 		log.Printf("Cursor error: %v", err)
// 		return nil, err
// 	}

// 	pageCount := int(math.Ceil(float64(total) / float64(limit)))
// 	hasNextPage := page < pageCount
// 	hasPrevPage := page > 1

// 	return &PaginationParams{
// 		Page:        page,
// 		Limit:       limit,
// 		Total:       total,
// 		PageCount:   pageCount,
// 		HasNextPage: hasNextPage,
// 		HasPrevPage: hasPrevPage,
// 	}, nil
// }

func PaginateCollection(
	c *gin.Context,
	collection *mongo.Collection,
	filter interface{},
	result interface{}, // Must be a pointer to a slice of the target type
) (*PaginationParams, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Println("Error counting documents:", err)
		return nil, err
	}

	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)),
	)
	if err != nil {
		log.Println("Error finding documents:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Ensure result is a pointer to a slice of the target type
	slicePtr, ok := result.(*[]models.Task)
	if !ok {
		return nil, errors.New("result parameter must be a pointer to a slice of the target type")
	}

	// Decode documents into the target slice
	for cursor.Next(context.Background()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Error decoding document:", err)
			return nil, err
		}
		*slicePtr = append(*slicePtr, task)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	pageCount := int(math.Ceil(float64(total) / float64(limit)))
	hasNextPage := page < pageCount
	hasPrevPage := page > 1

	return &PaginationParams{
		Page:        page,
		Limit:       limit,
		Total:       total,
		PageCount:   pageCount,
		HasNextPage: hasNextPage,
		HasPrevPage: hasPrevPage,
	}, nil
}

func SendPaginatedResponse(c *gin.Context, result interface{}, params *PaginationParams) {
	SendJSON(c, http.StatusOK, gin.H{
		"data":    result,
		"message": "Data fetched successfully",
		"meta": gin.H{
			"page":      params.Page,
			"perPage":   params.Limit,
			"total":     params.Total,
			"pageCount": params.PageCount,
			"nextPage": func() int {
				if params.HasNextPage {
					return params.Page + 1
				}
				return 0
			}(),
			"hasNextPage": params.HasNextPage,
			"hasPrevPage": params.HasPrevPage,
		},
	})
}
