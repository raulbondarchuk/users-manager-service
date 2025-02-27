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

// GetInstance — ленивая инициализация (singleton)
// При первом обращении вызывает HTTP() и возвращает *gin.Engine
func GetInstance() *gin.Engine {
	once.Do(func() {
		HTTP()
	})
	return instance
}

// MustLoad — запускает HTTP-сервер в отдельной горутине
func MustLoad() {
	router := GetInstance() // вызываем GetInstance() => создаёт/получает instance
	port := viper.GetInt("server.http.port")

	go func() {
		addr := ":" + strconv.Itoa(port)
		log.Printf("✅ Successfully started HTTP server: http://localhost%s", addr)
		if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
}

// HTTP — инициализация объекта *gin.Engine (CORS, middleware, маршруты)
// Вызывается один раз (через once.Do).
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
