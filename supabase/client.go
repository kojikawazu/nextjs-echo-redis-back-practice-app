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
	query := `
        SELECT id, user_id, description, completed, created_at, updated_at
        FROM todos
        ORDER BY created_at DESC
    `

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
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
			return nil, err
		}
		todos = append(todos, todo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return todos, nil
}

func CreateTodo(todo *models.TodoData) error {
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
		return err
	}

	return nil
}

func UpdateTodo(todo *models.TodoData) error {
	query := `
        UPDATE todos
        SET user_id = $1, description = $2, completed = $3, updated_at = $4
        WHERE id = $5
    `

	todo.UpdatedAt = time.Now()

	cmdTag, err := pool.Exec(ctx, query, todo.UserID, todo.Description, todo.Completed, todo.UpdatedAt, todo.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no todo found with id %s", todo.ID)
	}

	return nil
}

func DeleteTodo(id string) error {
	query := `
        DELETE FROM todos
        WHERE id = $1
    `

	cmdTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no todo found with id %s", id)
	}

	return nil
}

func ClosePool() {
	if pool != nil {
		pool.Close()
		log.Println("Supabase connection pool closed")
	}
}

func TestQuery() error {
	query := `SELECT 1`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			return err
		}
		fmt.Println("Test Query Result:", num)
	}

	return rows.Err()
}
