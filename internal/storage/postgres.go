package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Record struct {
	ID        int       `json:"id"`
	Filename  string    `json:"filename"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type DB struct {
	conn *sql.DB
}

func New() (*DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "voicetotext"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Migrate() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS transcriptions (
			id         SERIAL PRIMARY KEY,
			filename   TEXT NOT NULL,
			text       TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	return err
}

func (db *DB) Save(filename, text string) (int, error) {
	var id int
	err := db.conn.QueryRow(
		`INSERT INTO transcriptions (filename, text) VALUES ($1, $2) RETURNING id`,
		filename, text,
	).Scan(&id)
	return id, err
}

func (db *DB) History() ([]Record, error) {
	rows, err := db.conn.Query(
		`SELECT id, filename, text, created_at FROM transcriptions ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Filename, &r.Text, &r.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}
