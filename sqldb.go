// package findmyaadigoserver

// import (
// 	"log"
// 	"database/sql"
// )

// import (
// 	"database/sql"
// 	"log"
// 	"sync"

// 	_ "github.com/mattn/go-sqlite3"
// )

// var (
// 	instance *sql.DB
// 	once     sync.Once
// )

package findmyaadigoserver

import (
	"database/sql"
	"fmt"
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
		initDb()
	})
	return instance
}

// initDb initializes the in-memory SQLite DB and loads from file if needed
func initDb() {
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
func saveToDisk() error {
	// Use SQLite backup feature:
	_, err := instance.Exec("BACKUP TO './example.db';")
	return err
}

// addEntry inserts a new entry into the DB
func addEntry(timestamp int64, latitude float64, longitude float64) error {
    res, err := instance.Exec(
        "INSERT INTO location_history (timestamp, latitude, longitude) VALUES (?, ?, ?)",
        timestamp, latitude, longitude,
    )
    if err != nil {
        return err
    }
    fmt.Println(res)
    return nil
}

// getEntryRange retrieves all entries with timestamp between start and end (inclusive)
func getEntryRange(start int64, end int64) ([]LocationEntry, error) {
    rows, err := instance.Query(
        "SELECT timestamp, latitude, longitude FROM location_history WHERE timestamp >= ? AND timestamp <= ? ORDER BY timestamp ASC",
        start, end,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var entries []LocationEntry
    for rows.Next() {
        var e LocationEntry
        if err := rows.Scan(&e.timestamp, &e.latitude, &e.longitude); err != nil {
            return nil, err
        }
        entries = append(entries, e)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return entries, nil
}


// getLatestEntry retrieves the entry with the largest timestamp
func getLatestEntry() (*LocationEntry, error) {
	instance.Query("SELECT timestamp, latitude, longitude FROM location_history ORDER BY timestamp DESC LIMIT 1")
	return nil, nil
}

type LocationEntry struct {
	timestamp int64
	latitude  float64
	longitude float64
}
