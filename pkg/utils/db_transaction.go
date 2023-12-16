package utils

import (
	"gorm.io/gorm"
)

// Transactional executes the given callback within a database transaction
func Transactional(db *gorm.DB, callback func(tx *gorm.DB) error) error {
	// Begin a transaction
	tx := db.Begin()

	// Check for errors in beginning the transaction
	if err := tx.Error; err != nil {
		return err
	}

	// Execute the callback
	if err := callback(tx); err != nil {
		// If an error occurs, rollback the transaction and return the error
		tx.Rollback()
		return err
	}

	// Commit the transaction if no errors occurred
	return tx.Commit().Error
}
