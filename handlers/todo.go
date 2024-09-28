package handlers

import (
	"backend/models"
	"backend/myredis"
	"backend/supabase"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func GetTodos(c echo.Context) error {
	log.Println("Fetching todos...")
	ctx := context.Background()

	// Redisクライアントの取得
	rdb := myredis.GetRedisClient()
	log.Println("Connected to Redis successfully")

	// Redisからキャッシュを取得
	cachedTodos, err := rdb.Get(ctx, "todos").Result()
	if err == redis.Nil {
		// キャッシュがない場合はSupabaseからデータを取得
		log.Println("Cache miss: fetching data from Supabase")
		todos, err := supabase.FetchTodos()
		if err != nil {
			log.Printf("Error fetching todos from Supabase: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch todos",
			})
		}

		// キャッシュにデータを保存（10分間）
		todosJson, err := json.Marshal(todos)
		if err == nil {
			rdb.Set(ctx, "todos", todosJson, 10*time.Minute)
			log.Println("Data cached in Redis for 10 minutes")
		} else {
			log.Printf("Error marshalling todos for caching: %v", err)
		}

		log.Println("Data fetched from Supabase")
		return c.JSON(http.StatusOK, todos)
	} else if err != nil {
		// Redisに接続できないなどのエラー
		log.Printf("Error fetching cache from Redis: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch todos",
		})
	}

	// キャッシュにデータがあった場合
	log.Println("Cache hit: returning data from Redis")
	var todos []models.TodoData
	if err := json.Unmarshal([]byte(cachedTodos), &todos); err != nil {
		log.Printf("Error unmarshalling todos from cache: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch todos from cache",
		})
	}

	log.Println("Data fetched from Redis")
	return c.JSON(http.StatusOK, todos)
}

func CreateTodo(c echo.Context) error {
	log.Println("Creating todo...")
	var todo models.TodoData
	if err := c.Bind(&todo); err != nil {
		log.Printf("Error binding todo: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// 必要なフィールドのバリデーション
	if todo.UserID == "" || todo.Description == "" {
		log.Printf("UserID and Description are required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "UserID and Description are required",
		})
	}

	// UUIDの生成
	if todo.ID == "" {
		log.Println("Generating UUID for todo...")
		newUUID, err := uuid.NewRandom()
		if err != nil {
			log.Printf("Error generating UUID: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate UUID",
			})
		}
		log.Println("Generated UUID successfully")
		todo.ID = newUUID.String()
	}

	if err := supabase.CreateTodo(&todo); err != nil {
		log.Printf("Error creating todo: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create todo",
		})
	}

	log.Println("Created todo successfully")
	return c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c echo.Context) error {
	log.Println("Updating todo...")
	id := c.Param("id")
	var todo models.TodoData
	if err := c.Bind(&todo); err != nil {
		log.Printf("Error binding todo: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// IDの設定
	todo.ID = id

	log.Printf("Updating todo with ID: %s", id)
	if err := supabase.UpdateTodo(&todo); err != nil {
		log.Printf("Error updating todo: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update todo",
		})
	}

	log.Println("Updated todo successfully")
	return c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c echo.Context) error {
	log.Println("Deleting todo...")
	id := c.Param("id")

	if err := supabase.DeleteTodo(id); err != nil {
		log.Printf("Error deleting todo: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete todo",
		})
	}

	log.Println("Deleted todo successfully")
	return c.NoContent(http.StatusNoContent)
}
