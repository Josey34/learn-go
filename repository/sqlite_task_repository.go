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
		return nil, &domain.DatabaseError{Operation: "connect", Err: err}
	}

	repo := &SQLiteTaskRepository{db: db}

	if !dbExists {
		err := repo.initSchema()
		if err != nil {
			return nil, &domain.DatabaseError{Operation: "init schema", Err: err}
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
		return &domain.DatabaseError{Operation: "init schema", Err: err}
	}

	return nil
}

func (r *SQLiteTaskRepository) GetAll() ([]domain.Task, error) {
	query := "SELECT id, title, description, status, priority FROM tasks"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, &domain.DatabaseError{Operation: "get all tasks", Err: err}
	}

	defer rows.Close()

	tasks := []domain.Task{}

	for rows.Next() {
		var task domain.Task

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
		if err != nil {
			return nil, &domain.DatabaseError{Operation: "scan task row", Err: err}
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, &domain.DatabaseError{Operation: "iterate task rows", Err: err}
	}

	return tasks, nil
}

func (r *SQLiteTaskRepository) Create(task domain.Task) (domain.Task, error) {
	insertSQL := `INSERT INTO tasks (title, description, status, priority) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(insertSQL, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return domain.Task{}, &domain.DatabaseError{Operation: "insert task", Err: err}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}, &domain.DatabaseError{Operation: "get last insert id", Err: err}
	}

	task.ID = int(id)

	return task, nil
}

func (r *SQLiteTaskRepository) GetByID(id int) (domain.Task, error) {
	query := "SELECT * FROM tasks WHERE id = ?"

	row := r.db.QueryRow(query, id)

	var task domain.Task

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, &domain.NotFoundError{Resource: "Task Get By ID", ID: id}
		}
		return domain.Task{}, &domain.DatabaseError{Operation: "get task by id", Err: err}
	}

	return task, nil
}

func (r *SQLiteTaskRepository) Update(updatedTask domain.Task) (domain.Task, error) {

	_, err := r.GetByID(updatedTask.ID)
	if err != nil {
		return domain.Task{}, &domain.DatabaseError{Operation: "get task by id before update", Err: err}
	}

	query := "UPDATE tasks SET title = ?, description = ?, status = ?, priority = ? WHERE id = ?"

	_, err = r.db.Exec(query, updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Priority, updatedTask.ID)
	if err != nil {
		return domain.Task{}, &domain.DatabaseError{Operation: "update task", Err: err}
	}

	return updatedTask, nil
}

func (r *SQLiteTaskRepository) Delete(id int) error {
	_, err := r.GetByID(id)

	if err != nil {
		return &domain.DatabaseError{Operation: "get task by id before delete", Err: err}
	}

	query := "DELETE FROM tasks WHERE id = ?"

	_, err = r.db.Exec(query, id)
	if err != nil {
		return &domain.DatabaseError{Operation: "delete task", Err: err}
	}

	return nil
}

func (r *SQLiteTaskRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
