package main

import (
	"context"

	"github.com/Untanky/modern-auth/internal/gin"
)

func main() {
	app := gin.NewGinApplication()

	app.Start(context.Background())
}
