package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// wraps sql db connection pool
type SnippetService struct {
	DB *sql.DB
	//InsertStatement *sql.Stmt
	//LatestStatement *sql.Stmt
	//GetStatement    *sql.Stmt
}

//func NewSnippetService(db *sql.DB) (*SnippetService, error) {
//	insertStatement, err := db.Prepare(`INSERT INTO snippets (title, content, created, expires)
//	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`)
//	if err != nil {
//		return nil, err
//	}
//
//	latestStatement, err := db.Prepare(`INSERT INTO snippets (title, content, created, expires)
//	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`)
//	if err != nil {
//		return nil, err
//	}
//
//	getStatement, err := db.Prepare(`INSERT INTO snippets (title, content, created, expires)
//	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`)
//	if err != nil {
//		return nil, err
//	}
//
//	return &SnippetService{DB: db, LatestStatement: latestStatement, InsertStatement: insertStatement, GetStatement: getStatement}, nil
//}

func (s *SnippetService) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//result, err := s.InsertStatement.Exec(title, content, expires)
	result, err := s.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SnippetService) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	//row := s.GetStatement.QueryRow(id)
	row := s.DB.QueryRow(stmt, id)
	var snippet Snippet

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return snippet, nil
}

func (s *SnippetService) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	//rows, err := s.LatestStatement.Query()
	rows, err := s.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// important: we should close rows after error check
	defer rows.Close()
	var snippets []Snippet
	for rows.Next() {
		var snippet Snippet
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
