package http

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	instance *gin.Engine
	once     sync.Once
)

// GetInstance — lazy initialization (singleton)
// On first call, it calls HTTP() and returns *gin.Engine
func GetInstance() *gin.Engine {
	once.Do(func() {
		HTTP()
	})
	return instance
}

// MustLoad — starts HTTP-server in a separate goroutine
func MustLoad() {
	router := GetInstance() // call GetInstance() => creates/gets instance
	port := viper.GetInt("server.http.port")

	go func() {
		addr := ":" + strconv.Itoa(port)
		log.Printf("✅ Successfully started HTTP server: http://localhost%s", addr)
		if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
}

// HTTP — initialization of the *gin.Engine object (CORS, middleware, routes)
// Called once (through once.Do).
func HTTP() {
	setMode()
	instance = gin.Default()
	setCors(instance)
	instance.Use(TimeoutMiddleware(viper.GetString("server.http.timeout")))
	instance.Use(RouteLogger())
	InitRoutes(instance)
}

func setMode() {
	if mode := viper.GetString("server.http.mode"); mode == "local" || mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
