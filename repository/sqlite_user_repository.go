package repository

import (
	"database/sql"
	"os"
	"task-manager-api/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(dbPath string) (*SQLiteUserRepository, error) {
	dbExists := false

	if _, err := os.Stat(dbPath); err == nil {
		dbExists = true
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, &domain.DatabaseError{Operation: "open", Err: err}
	}

	err = db.Ping()
	if err != nil {
		return nil, &domain.DatabaseError{Operation: "ping", Err: err}
	}

	repo := &SQLiteUserRepository{db: db}

	if !dbExists {
		err = repo.initSchema()
		if err != nil {
			return nil, err
		}
	}

	return repo, nil
}

func (r *SQLiteUserRepository) initSchema() error {
	sqlBytes, err := os.ReadFile("migrations/add_users_table.sql")
	if err != nil {
		return &domain.DatabaseError{Operation: "read migration", Err: err}
	}

	sqlString := string(sqlBytes)
	_, err = r.db.Exec(sqlString)
	if err != nil {
		return &domain.DatabaseError{Operation: "init schema", Err: err}
	}

	return nil
}

func (r *SQLiteUserRepository) Create(user domain.User) (domain.User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?)`

	result, err := r.db.Exec(query, user.Email, user.PasswordHash)
	if err != nil {
		return domain.User{}, &domain.DatabaseError{Operation: "insert", Err: err}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.User{}, &domain.DatabaseError{Operation: "retrieve last insert id", Err: err}
	}

	user.ID = int(id)

	return user, nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (domain.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	var user domain.User

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, &domain.NotFoundError{Resource: "User", ID: 0}
		}
		return domain.User{}, &domain.DatabaseError{Operation: "get user by email", Err: err}
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetByID(id int) (domain.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var user domain.User

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, &domain.NotFoundError{Resource: "User", ID: id}
		}
		return domain.User{}, &domain.DatabaseError{Operation: "get user by id", Err: err}
	}

	return user, nil
}

func (r *SQLiteUserRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
