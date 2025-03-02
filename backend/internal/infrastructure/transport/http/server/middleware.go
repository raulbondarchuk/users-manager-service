package http

import (
	"app/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware establece un tiempo de espera para el procesamiento de la solicitud
func TimeoutMiddleware(timeoutString string) gin.HandlerFunc {
	timeout, err := time.ParseDuration(timeoutString) // Convert string to time.Duration
	if err != nil {
		log.Fatalf("Invalid timeout value: %v", err) // Handle error if value is incorrect
	}

	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the context in the request
		c.Request = c.Request.WithContext(ctx)

		// Process the request
		c.Next()

		// Check if the timeout has expired
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
			c.Abort() // Interrumpe el procesamiento futuro
		}
	}
}

// -- Routes logger --

func printRoutes(router *gin.Engine) {
	routes := router.Routes()
	routesPrint := ""
	// Create a map to group routes by prefix
	groupedRoutes := make(map[string][]gin.RouteInfo)

	// Group routes by prefix
	for _, route := range routes {
		prefix := strings.Split(route.Path, "/")[1] // Get the first segment of the path
		groupedRoutes[prefix] = append(groupedRoutes[prefix], route)
	}

	routesPrint += "Registered routes:\n"
	routesPrint += "------------------\n"

	// Get all prefixes and sort them
	var prefixes []string
	for prefix := range groupedRoutes {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)

	// Print routes by groups
	for _, prefix := range prefixes {
		routesPrint += fmt.Sprintf("\n[%s]\n", strings.ToUpper(prefix))
		for _, route := range groupedRoutes[prefix] {
			parts := strings.Split(route.Handler, "/")
			handler := parts[len(parts)-1]
			routesPrint += fmt.Sprintf("%-6s  %-25s | Handler: %v\n", route.Method, route.Path, handler)
		}
	}
	routesPrint += "\n------------------\n"
	fmt.Println(routesPrint)
}

func RouteLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Print("\n")
		logger.GetLogger().RouteInfo("Route accessed", map[string]interface{}{
			"path":   c.Request.RequestURI,
			"method": c.Request.Method,
		})
		c.Next()
	}
}
