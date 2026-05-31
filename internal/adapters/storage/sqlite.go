package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ekosachev/loom/internal/domain/models"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteAdapter struct {
	db *sql.DB
}

func NewSQLiteAdapter(dsn string) (*SQLiteAdapter, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS workspaces (
		name TEXT PRIMARY KEY
	);
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		parent_id INTEGER,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (parent_id) REFERENCES messages (id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS branches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		workspace_id TEXT NOT NULL,
		current_message_id INTEGER,
		FOREIGN KEY (workspace_id) REFERENCES workspaces (name) ON DELETE CASCADE,
		FOREIGN KEY (current_message_id) REFERENCES messages (id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS state (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`
	if _, err = db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteAdapter{db}, err
}

func (s *SQLiteAdapter) SaveWorkspace(ctx context.Context, ws *models.Workspace) error {
	query := `
	INSERT INTO workspaces (name) VALUES (?) ON CONFLICT(name) DO NOTHING;
	`
	_, err := s.db.ExecContext(ctx, query, ws.Name)
	return err
}

func (s *SQLiteAdapter) GetWorkspace(ctx context.Context, name string) (*models.Workspace, error) {
	query := `SELECT name FROM workspaces WHERE name = ?;`
	row := s.db.QueryRowContext(ctx, query, name)

	ws := models.Workspace{}

	if err := row.Scan(&ws.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &ws, nil
}

func (s *SQLiteAdapter) GetAllWorkspaces(ctx context.Context) ([]models.Workspace, error) {
	query := `SELECT name FROM workspaces`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []models.Workspace{}, err
	}
	workspaces := make([]models.Workspace, 0)

	for rows.Next() {
		ws := models.Workspace{}
		rows.Scan(&ws.Name)
		workspaces = append(workspaces, ws)
	}

	return workspaces, nil
}

func (s *SQLiteAdapter) SaveBranch(ctx context.Context, branch *models.Branch) error {
	query := `INSERT INTO branches (name, current_message_id, workspace_id)
	VALUES (?, ?, ?);`

	reslut, err := s.db.ExecContext(ctx, query, branch.Name, branch.CurrentMessageID, branch.WorkspaceID)
	branchID, err := reslut.LastInsertId()
	if err != nil {
		return err
	}
	branch.ID = branchID

	return err
}

func (s *SQLiteAdapter) GetBranch(ctx context.Context, branchID int64) (*models.Branch, error) {
	query := `SELECT id, name, current_message_id, workspace_id FROM branches WHERE id = ?`

	row := s.db.QueryRowContext(ctx, query, branchID)

	var br models.Branch

	if err := row.Scan(&br.ID, &br.Name, &br.CurrentMessageID, &br.WorkspaceID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &br, nil
}

func (s *SQLiteAdapter) UpdateBranchHead(ctx context.Context, branchID int64, messageID int64) error {
	query := `UPDATE branches SET current_message_id = ? WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query, messageID, branchID)

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("branch not found")
	}
	return nil
}

func (s *SQLiteAdapter) GetBranchByName(ctx context.Context, workspaceName string, branchName string) (*models.Branch, error) {
	query := `SELECT id, name, workspace_id, current_message_id FROM branches WHERE workspace_id = ? AND name = ?`
	row := s.db.QueryRowContext(ctx, query, workspaceName, branchName)
	var br models.Branch

	if err := row.Scan(&br.ID, &br.Name, &br.WorkspaceID, &br.CurrentMessageID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, err
	}

	return &br, nil
}

func (s *SQLiteAdapter) GetAllBranches(ctx context.Context, workspaceName string) ([]models.Branch, error) {
	query := `SELECT id, name, workspace_id, current_message_id FROM branches WHERE workspace_id = ?`
	rows, err := s.db.QueryContext(ctx, query, workspaceName)
	if err != nil {
		return []models.Branch{}, err
	}

	branches := make([]models.Branch, 0)

	for rows.Next() {
		var br models.Branch
		rows.Scan(&br.ID, &br.Name, &br.WorkspaceID, &br.CurrentMessageID)
		branches = append(branches, br)
	}

	return branches, nil
}

func (s *SQLiteAdapter) SaveMessage(ctx context.Context, msg *models.Message) error {
	query := `INSERT INTO messages (parent_id, role, content) VALUES (?, ?, ?);`

	res, err := s.db.ExecContext(ctx, query, msg.ParentID, msg.Role, msg.Content)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	msg.ID = id
	return nil
}

func (s *SQLiteAdapter) GetMessage(ctx context.Context, id int64) (*models.Message, error) {
	query := `SELECT id, parent_id, role, content, timestamp FROM messages WHERE id = ?;`
	row := s.db.QueryRowContext(ctx, query, id)

	var m models.Message
	if err := row.Scan(&m.ID, &m.ParentID, &m.Role, &m.Content, &m.Timestamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (s *SQLiteAdapter) GetThreadContext(ctx context.Context, headMessageID int64) ([]models.Message, error) {
	query := `
	WITH RECURSIVE thread AS (
		SELECT id, parent_id, role, content, timestamp FROM messages WHERE id = ?
		UNION ALL
		SELECT m.id, m.parent_id, m.role, m.content, m.timestamp 
		FROM messages m
		JOIN thread t ON m.id = t.parent_id
	)
	SELECT id, parent_id, role, content, timestamp FROM thread;`

	rows, err := s.db.QueryContext(ctx, query, headMessageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ParentID, &m.Role, &m.Content, &m.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}

func (s *SQLiteAdapter) GetState(ctx context.Context, key string) (*string, error) {
	query := `SELECT value FROM state WHERE key = ?`
	row := s.db.QueryRowContext(ctx, query, key)

	v := ""
	if err := row.Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (s *SQLiteAdapter) SetState(ctx context.Context, key string, value string) error {
	query := `
	INSERT INTO state (key, value)
	VALUES (?, ?)
	ON CONFLICT(key)
	DO UPDATE SET value = ?;
	`

	_, err := s.db.ExecContext(ctx, query, key, value, value)

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteAdapter) Close() error {
	return s.db.Close()
}
