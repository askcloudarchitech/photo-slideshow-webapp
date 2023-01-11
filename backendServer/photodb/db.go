package photodb

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Database *sql.DB
}

type PhotoDetails struct {
	Name       string
	TimeTaken  int
	LastViewed int
}

func (d *DB) InitDB() {
	db, err := sql.Open("sqlite3", "/data/database/photo.db")
	if err != nil {
		log.Fatal(err)
	}
	d.Database = db

	initStmt := `
	CREATE TABLE IF NOT EXISTS PHOTOS (FILENAME string unique not null primary key, LAST_VIEWED int, TIME_TAKEN int);
	CREATE TABLE IF NOT EXISTS PHOTO_SEMAPHORE (LAST_POLL_TIMESTAMP int);
	`
	_, err = d.Database.Exec(initStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, initStmt)
	}
}

func (d *DB) AddPhoto(name string, timeTaken int) error {
	tx, err := d.Database.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into PHOTOS(FILENAME, LAST_VIEWED, TIME_TAKEN) values(?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, 0, timeTaken)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (d *DB) UpdateLastViewed(name string) error {
	tx, err := d.Database.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE PHOTOS SET LAST_VIEWED=? WHERE FILENAME=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(time.Now().Unix(), name)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (d *DB) GetNextSlideshowPhoto() (PhotoDetails, error) {
	rows, err := d.Database.Query("select FILENAME, LAST_VIEWED, TIME_TAKEN from PHOTOS ORDER BY LAST_VIEWED LIMIT 1")
	if err != nil {
		return PhotoDetails{}, err
	}
	defer rows.Close()

	for rows.Next() {
		details := PhotoDetails{}
		err := rows.Scan(&details.Name, &details.LastViewed, &details.TimeTaken)
		if err != nil {
			return PhotoDetails{}, nil
		}
		return details, nil
	}

	err = rows.Err()
	return PhotoDetails{}, err
}

func (d *DB) GetPhotoList(count int, offset int) ([]PhotoDetails, error) {

	return []PhotoDetails{}, nil
}
