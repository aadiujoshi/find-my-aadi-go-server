package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"server/db"
)

func main2() {
	fmt.Println("=== SQLite In-Memory Database Test ===")

	// Clean up any existing test database file
	if err := os.Remove("./location-history.db"); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Could not remove existing database file: %v\n", err)
	}

	// Test 1: Initialize database
	fmt.Println("1. Testing database initialization...")
	database := db.GetDB()
	if database == nil {
		log.Fatal("Failed to get database instance")
	}
	fmt.Println("✓ Database initialized successfully")

	// Test 2: Add some sample entries
	fmt.Println("2. Testing AddEntry functionality...")
	
	// Current time and some test timestamps
	now := time.Now().Unix()
	
	testEntries := []struct {
		timestamp int64
		lat       float64
		lng       float64
		desc      string
	}{
		{now - 3600, 37.7749, -122.4194, "San Francisco (1 hour ago)"},
		{now - 1800, 37.7849, -122.4094, "Near San Francisco (30 min ago)"},
		{now - 900, 37.7949, -122.3994, "Moving north (15 min ago)"},
		{now, 37.8049, -122.3894, "Current location"},
	}

	for i, entry := range testEntries {
		err := db.AddEntry(entry.timestamp, entry.lat, entry.lng)
		if err != nil {
			log.Printf("Failed to add entry %d: %v", i+1, err)
		} else {
			fmt.Printf("✓ Added entry %d: %s (%.4f, %.4f)\n", i+1, entry.desc, entry.lat, entry.lng)
		}
	}
	fmt.Println()

	// Test 3: Get all entries in range
	fmt.Println("3. Testing GetEntryRange functionality...")
	
	// Get all entries
	allEntries, err := db.GetEntryRange(now-4000, now+100)
	if err != nil {
		log.Printf("Failed to get entry range: %v", err)
	} else {
		fmt.Printf("✓ Retrieved %d entries in full range:\n", len(allEntries))
		for i, entry := range allEntries {
			timeStr := time.Unix(entry.Timestamp, 0).Format("15:04:05")
			fmt.Printf("  Entry %d: %s - (%.4f, %.4f)\n", i+1, timeStr, entry.Latitude, entry.Longitude)
		}
	}
	fmt.Println()

	// Test 4: Get partial range
	fmt.Println("4. Testing partial range query...")
	partialEntries, err := db.GetEntryRange(now-2000, now-500)
	if err != nil {
		log.Printf("Failed to get partial range: %v", err)
	} else {
		fmt.Printf("✓ Retrieved %d entries in partial range (last 33-8 minutes):\n", len(partialEntries))
		for i, entry := range partialEntries {
			timeStr := time.Unix(entry.Timestamp, 0).Format("15:04:05")
			fmt.Printf("  Entry %d: %s - (%.4f, %.4f)\n", i+1, timeStr, entry.Latitude, entry.Longitude)
		}
	}
	fmt.Println()

	// Test 5: Get latest entry
	fmt.Println("5. Testing GetLatestEntry functionality...")
	latest, err := db.GetLatestEntry()
	if err != nil {
		log.Printf("Failed to get latest entry: %v", err)
	} else if latest == nil {
		fmt.Println("✗ No latest entry found")
	} else {
		timeStr := time.Unix(latest.Timestamp, 0).Format("15:04:05")
		fmt.Printf("✓ Latest entry: %s - (%.4f, %.4f)\n", timeStr, latest.Latitude, latest.Longitude)
	}
	fmt.Println()

	// Test 6: Save to disk
	fmt.Println("6. Testing SaveToDisk functionality...")
	err = db.SaveToDisk()
	if err != nil {
		log.Printf("Failed to save to disk: %v", err)
	} else {
		fmt.Println("✓ Successfully saved database to disk")
		
		// Verify file was created
		if _, err := os.Stat("./location-history.db"); os.IsNotExist(err) {
			fmt.Println("✗ Database file was not created")
		} else {
			fmt.Println("✓ Database file exists on disk")
		}
	}
	fmt.Println()

	// Test 7: Test singleton pattern
	fmt.Println("7. Testing singleton pattern...")
	db2 := db.GetDB()
	if db2 == database {
		fmt.Println("✓ Singleton pattern working - same instance returned")
	} else {
		fmt.Println("✗ Singleton pattern failed - different instances returned")
	}
	fmt.Println()

	// Test 8: Add more entries and test again
	fmt.Println("8. Adding more test data...")
	moreEntries := []struct {
		timestamp int64
		lat       float64
		lng       float64
		desc      string
	}{
		{now + 300, 37.8149, -122.3794, "5 minutes later"},
		{now + 600, 37.8249, -122.3694, "10 minutes later"},
	}

	for i, entry := range moreEntries {
		err := db.AddEntry(entry.timestamp, entry.lat, entry.lng)
		if err != nil {
			log.Printf("Failed to add additional entry %d: %v", i+1, err)
		} else {
			fmt.Printf("✓ Added additional entry %d: %s (%.4f, %.4f)\n", i+1, entry.desc, entry.lat, entry.lng)
		}
	}
	fmt.Println()

	// Test 9: Final range query with all data
	fmt.Println("9. Final verification - all entries:")
	finalEntries, err := db.GetEntryRange(now-4000, now+1000)
	if err != nil {
		log.Printf("Failed to get final entries: %v", err)
	} else {
		fmt.Printf("✓ Total entries in database: %d\n", len(finalEntries))
		for i, entry := range finalEntries {
			timeStr := time.Unix(entry.Timestamp, 0).Format("15:04:05")
			fmt.Printf("  Entry %d: %s - (%.4f, %.4f)\n", i+1, timeStr, entry.Latitude, entry.Longitude)
		}
	}
	fmt.Println()

	// Test 10: Test empty range
	fmt.Println("10. Testing empty range query...")
	emptyEntries, err := db.GetEntryRange(now+2000, now+3000)
	if err != nil {
		log.Printf("Failed to query empty range: %v", err)
	} else {
		fmt.Printf("✓ Empty range query returned %d entries (expected 0)\n", len(emptyEntries))
	}
	fmt.Println()

	// Test 11: Final save to disk
	fmt.Println("11. Final save to disk...")
	err = db.SaveToDisk()
	if err != nil {
		log.Printf("Failed final save to disk: %v", err)
	} else {
		fmt.Println("✓ Final save to disk completed")
	}
	
	fmt.Println("\n=== All tests completed ===")
}