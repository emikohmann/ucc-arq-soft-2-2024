package users

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`                    // Auto-increment primary key
	Username string `gorm:"size:100;not null;unique" binding:"required"` // Unique username, required
	Password string `gorm:"size:255;not null" binding:"required"`        // Password field, required
}
