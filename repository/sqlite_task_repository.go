package repository

import (
	"database/sql"
	"os"
	"task-manager-api/domain"

	_ "modernc.org/sqlite"
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(dbPath string) (*SQLiteTaskRepository, error) {
	dbExists := false

	if _, err := os.Stat(dbPath); err == nil {
		dbExists = true
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &SQLiteTaskRepository{db: db}

	if !dbExists {
		err := repo.initSchema()
		if err != nil {
			return nil, err
		}
	}
	return repo, nil
}

func (r *SQLiteTaskRepository) initSchema() error {
	sqlBytes, err := os.ReadFile("migrations/init.sql")

	if err != nil {
		return err
	}

	sqlString := string(sqlBytes)

	_, err = r.db.Exec(sqlString)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteTaskRepository) GetAll() ([]domain.Task, error) {
	query := "SELECT id, title, description, status, priority FROM tasks"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []domain.Task{}

	for rows.Next() {
		var task domain.Task

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *SQLiteTaskRepository) Create(task domain.Task) (domain.Task, error) {
	insertSQL := `INSERT INTO tasks (title, description, status, priority) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(insertSQL, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return domain.Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}, err
	}

	task.ID = int(id)

	return task, nil
}

func (r *SQLiteTaskRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
