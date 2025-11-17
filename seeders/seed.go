package seeders

import (
	"database/sql"

	"fmt"
)

func RunAllSeeders(db *sql.DB) error {
	fmt.Println("Running database seeders...")

	if err := SeedEcoQuestions(db); err != nil {
		return fmt.Errorf("eco questions seeder failed: %w", err)
	}

	fmt.Println("All seeders completed")
	return nil
}
