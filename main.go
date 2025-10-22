package main

import (
	"backend/config"
	"backend/routes"
	"fmt"
	"log"
	"os"

	_ "backend/docs" // generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func main() {
	// === Setup logging ke file backend.log ===
	logFile, err := os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Tidak dapat membuka atau membuat file log:", err)
	}
	gin.DefaultWriter = logFile

	// Gunakan gin.Default() (seperti semula)
	r := gin.Default()

	// Tambahkan custom format log agar lebih rapi
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Format waktu: YYYY-MM-DD HH:MM:SS
		return fmt.Sprintf("[%s] %s %s %d %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	// === Middleware CORS (asli) ===
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// === Database dan routes (tidak diubah) ===
	config.ConnectDatabase()
	routes.SetupRoutes(r)

	// === Swagger route (asli) ===
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// === Jalankan server ===
	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
