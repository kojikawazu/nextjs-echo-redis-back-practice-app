package handlers

import (
	"backend/models"
	"backend/supabase"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetTodos(c echo.Context) error {
	todos, err := supabase.FetchTodos()
	if err != nil {
		c.Logger().Errorf("Error fetching todos: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch todos",
		})
	}

	return c.JSON(http.StatusOK, todos)
}

func CreateTodo(c echo.Context) error {
	var todo models.TodoData
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// 必要なフィールドのバリデーション
	if todo.UserID == "" || todo.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "UserID and Description are required",
		})
	}

	// UUIDの生成
	if todo.ID == "" {
		newUUID, err := uuid.NewRandom()
		if err != nil {
			c.Logger().Errorf("Error generating UUID: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate UUID",
			})
		}
		todo.ID = newUUID.String()
	}

	if err := supabase.CreateTodo(&todo); err != nil {
		c.Logger().Errorf("Error creating todo: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create todo",
		})
	}

	return c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var todo models.TodoData
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// IDの設定
	todo.ID = id

	if err := supabase.UpdateTodo(&todo); err != nil {
		c.Logger().Errorf("Error updating todo: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update todo",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c echo.Context) error {
	id := c.Param("id")

	if err := supabase.DeleteTodo(id); err != nil {
		c.Logger().Errorf("Error deleting todo: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete todo",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
