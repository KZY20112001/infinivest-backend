package main

import (
	"github.com/KZY20112001/infinivest-backend/internal/initializers"
	"github.com/KZY20112001/infinivest-backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.SetupDB()
}

func main() {
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run()
}
