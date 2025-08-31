package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var DB_FILE_NAME = "./location-history.db"
var (
	instance *sql.DB
	once sync.Once
)

// GetDB returns the singleton SQLite in-memory DB instance
func GetDB() *sql.DB {
	once.Do(func() {
		InitDb()
	})
	return instance
}

func ensureFileDB() {
    db, err := sql.Open("sqlite3", DB_FILE_NAME)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS location_history (
        Timestamp INTEGER,
        Latitude REAL,
        Longitude REAL
    );`)
    if err != nil {
        log.Fatal(err)
    }
}

// initDb initializes the in-memory SQLite DB and loads from file if needed
func InitDb() {
	ensureFileDB()
	
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Failed to open in-memory DB:", err)
	}

	// Load existing file-based DB into memory
	fileDB, err := sql.Open("sqlite3", DB_FILE_NAME)
	if err != nil {
		log.Fatal("Failed to open file DB:", err)
	}
	defer fileDB.Close()

	_, err = db.Exec("ATTACH DATABASE '" + DB_FILE_NAME + "' AS filedb;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE location_history AS SELECT * FROM filedb.location_history;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DETACH DATABASE filedb;")
	if err != nil {
		log.Fatal(err)
	}

	instance = db
}

// saveToDisk writes the in-memory DB back to the file for persistence
func SaveToDisk() error {
    // Attach the original file and copy data back
    _, err := instance.Exec("ATTACH DATABASE ? AS diskdb", DB_FILE_NAME)
    if err != nil {
        return err
    }
    defer instance.Exec("DETACH DATABASE diskdb")
    
    _, err = instance.Exec("DELETE FROM diskdb.location_history")
    if err != nil {
        return err
    }
    
    _, err = instance.Exec("INSERT INTO diskdb.location_history SELECT * FROM location_history")
    return err
}

// addEntry inserts a new entry into the DB
func AddEntry(Timestamp int64, Latitude float64, Longitude float64) error {
    _, err := instance.Exec(
        "INSERT INTO location_history (Timestamp, Latitude, Longitude) VALUES (?, ?, ?)",
        Timestamp, Latitude, Longitude,
    )
    if err != nil {
        return err
    }
    return nil
}

// getEntryRange retrieves all entries with Timestamp between start and end (inclusive)
func GetEntryRange(start int64, end int64) ([]LocationEntry, error) {
    rows, err := instance.Query(
        "SELECT Timestamp, Latitude, Longitude FROM location_history WHERE Timestamp >= ? AND Timestamp <= ? ORDER BY Timestamp ASC",
        start, end,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var entries []LocationEntry
    for rows.Next() {
        var e LocationEntry
        if err := rows.Scan(&e.Timestamp, &e.Latitude, &e.Longitude); err != nil {
            return nil, err
        }
        entries = append(entries, e)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return entries, nil
}


// getLatestEntry retrieves the entry with the largest Timestamp
func GetLatestEntry() (*LocationEntry, error) {
    row := instance.QueryRow("SELECT Timestamp, Latitude, Longitude FROM location_history ORDER BY Timestamp DESC LIMIT 1")
    
    var e LocationEntry
    err := row.Scan(&e.Timestamp, &e.Latitude, &e.Longitude)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // No entries found
        }
        return nil, err
    }
    return &e, nil
}

type LocationEntry struct {
	Timestamp int64
	Latitude  float64
	Longitude float64
}