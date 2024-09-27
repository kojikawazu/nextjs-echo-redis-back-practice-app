package supabase

import (
	"backend/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ctx  = context.Background()
	pool *pgxpool.Pool
)

func InitSupabase() error {
	log.Println("Initializing Supabase client...")
	supabaseURL := os.Getenv("SUPABASE_URL") + "?sslmode=require"

	config, err := pgxpool.ParseConfig(supabaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
		return fmt.Errorf("unable to parse database URL: %v", err)
	}

	// コネクションプールの設定
	config.MaxConns = 10 // 必要に応じて調整
	config.MaxConnIdleTime = 30 * time.Second
	// Prepared Statementの競合を防ぐためにSimple Protocolを優先
	config.ConnConfig.PreferSimpleProtocol = true

	log.Println("Connecting supabase database...")
	pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to Supabase: %v", err)
		return fmt.Errorf("unable to connect to Supabase: %v", err)
	}

	// 接続の確認
	log.Println("Pinging supabase database...")
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping Supabase: %v", err)
		return fmt.Errorf("unable to ping Supabase: %v", err)
	}

	log.Println("Connected to Supabase successfully")
	return nil
}

func FetchTodos() ([]models.TodoData, error) {
	log.Println("Fetching todos from Supabase...")
	query := `
        SELECT id, user_id, description, completed, created_at, updated_at
        FROM todos
        ORDER BY created_at DESC
    `

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to fetch todos: %v", err)
		return nil, err
	}
	log.Println("Fetched todos successfully")
	defer rows.Close()

	var todos []models.TodoData

	for rows.Next() {
		var todo models.TodoData
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			log.Fatalf("Failed to scan todo: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}

	if rows.Err() != nil {
		log.Fatalf("Failed to fetch todos: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched %d todos", len(todos))
	return todos, nil
}

func CreateTodo(todo *models.TodoData) error {
	log.Println("Creating todo in Supabase...")
	query := `
        INSERT INTO todos (id, user_id, description, completed, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	// タイムスタンプの設定
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	_, err := pool.Exec(ctx, query, todo.ID, todo.UserID, todo.Description, todo.Completed, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		log.Fatalf("Failed to create todo: %v", err)
		return err
	}

	log.Println("Created todo successfully")
	return nil
}

func UpdateTodo(todo *models.TodoData) error {
	log.Println("Updating todo in Supabase...")
	query := `
        UPDATE todos
        SET user_id = $1, description = $2, completed = $3, updated_at = $4
        WHERE id = $5
    `

	todo.UpdatedAt = time.Now()

	cmdTag, err := pool.Exec(ctx, query, todo.UserID, todo.Description, todo.Completed, todo.UpdatedAt, todo.ID)
	if err != nil {
		log.Fatalf("Failed to update todo: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		log.Fatalf("No todo found with id %s", todo.ID)
		return fmt.Errorf("no todo found with id %s", todo.ID)
	}

	log.Println("Updated todo successfully")
	return nil
}

func DeleteTodo(id string) error {
	log.Println("Deleting todo in Supabase...")
	query := `
        DELETE FROM todos
        WHERE id = $1
    `

	cmdTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		log.Fatalf("Failed to delete todo: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		log.Fatalf("No todo found with id %s", id)
		return fmt.Errorf("no todo found with id %s", id)
	}

	log.Println("Deleted todo successfully")
	return nil
}

func ClosePool() {
	if pool != nil {
		pool.Close()
		log.Println("Supabase connection pool closed")
	}
}

func TestQuery() error {
	log.Println("Testing query...")
	query := `SELECT 1`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to test query: %v", err)
		return err
	}
	log.Println("Test query successful")
	defer rows.Close()

	for rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			log.Fatalf("Failed to scan test query result: %v", err)
			return err
		}
		fmt.Println("Test Query Result:", num)
	}

	log.Println("Test query completed")
	return rows.Err()
}
