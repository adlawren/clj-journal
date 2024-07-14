package main

import (
	"bujo/lib"
	"flag"
	"log"
)

func main() {
	var dailyMigration bool
	var monthlyMigration bool

	flag.BoolVar(&dailyMigration, "m", false, "Run daily migration")
	flag.BoolVar(&monthlyMigration, "M", false, "Run monthly migration")

	flag.Parse()

	if dailyMigration {
		if err := lib.RunDailyMigration(); err != nil {
			log.Fatalf("Failed to run daily migration: %s", err)
		}
	} else if monthlyMigration {
		if err := lib.RunMonthlyMigration(); err != nil {
			log.Fatalf("Failed to run monthly migration: %s", err)
		}
	}
}
