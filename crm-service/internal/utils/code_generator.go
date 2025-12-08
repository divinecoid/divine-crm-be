package utils

import (
	"fmt"
	"gorm.io/gorm"
)

// GenerateCode generates a unique code for entities
func GenerateCode(db *gorm.DB, prefix string, model interface{}) string {
	// ✅ Handle nil database - use simple counter
	if db == nil {
		// Simple sequential code without database check
		return fmt.Sprintf("%s%06d", prefix, 1)
	}

	// Count existing records
	var count int64
	db.Model(model).Count(&count)

	// Generate code with prefix and padded number
	return fmt.Sprintf("%s%06d", prefix, count+1)
}
