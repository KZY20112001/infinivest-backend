package initializers

import "github.com/KZY20112001/infinivest-backend/internal/db"

func SetupDB() {
	db.ConnectToPostgres()
	db.ConnectToRedis()
}
