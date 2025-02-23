package router

import (
	"bytes"
	"encoding/json"
	"go-drive-project/entities"
	"go-drive-project/logger"
	"go-drive-project/utility"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) *gin.Engine {

	// Routes to exclude from JWT verification
	excludedRoutes := []string{"/g-drive/signup", "/g-drive/login"}

	// Define routes
	for _, route := range routes {
		switch route.Method {
		case "POST":
			r.POST("/g-drive"+route.Pattern, PostPutMiddleware(), TokenMiddleware(excludedRoutes), route.HandlerFunc)
		case "PUT":
			r.PUT("/g-drive"+route.Pattern, PostPutMiddleware(), TokenMiddleware(excludedRoutes), route.HandlerFunc)
		case "GET":
			r.GET("/g-drive"+route.Pattern, GetDeleteMiddleware(), TokenMiddleware(excludedRoutes), route.HandlerFunc)
		case "DELETE":
			r.DELETE("/g-drive"+route.Pattern, GetDeleteMiddleware(), TokenMiddleware(excludedRoutes), route.HandlerFunc)
		case "OPTIONS":
			r.OPTIONS("/g-drive"+route.Pattern, OptionsMiddleware(), route.HandlerFunc)
		}
	}

	return r
}

// PostPutMiddleware allows only POST and PUT requests
func PostPutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type, Authorization,Auth")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Autho")
		c.Writer.Header().Set("Cache-Control", "no-cache,no-store")

		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only POST and PUT methods are allowed"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetDeleteMiddleware allows only GET and DELETE requests
func GetDeleteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type, Authorization,Auth")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Autho")
		c.Writer.Header().Set("Cache-Control", "no-cache,no-store")

		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodDelete {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only GET and DELETE methods are allowed"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionsMiddleware handles preflight OPTIONS requests
func OptionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type, Authorization,Auth")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Autho")
		c.Writer.Header().Set("Cache-Control", "no-cache,no-store")

		// Send a 200 OK response for preflight requests
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, gin.H{"message": "Preflight request successful"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func TokenMiddleware(excludedRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip JWT verification for excluded routes
		for _, route := range excludedRoutes {
			if strings.HasPrefix(path, route) {
				c.Next()
				return
			}
		}

		// Get token from Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Token"})
			c.Abort()
			return
		}

		// Validate Token
		claims, success, err, errType := utility.ValidateToken(token)
		if !success {
			log.Printf("Token error: %v, Type: %v", err, errType)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// Inject User ID into request body for POST & PUT requests
		logger.Log.Println(c.Request.URL.Path)
		if (c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut) && !strings.HasPrefix(c.Request.URL.Path, "/g-drive/uploadFile") {
			var reqData map[string]interface{}

			// Read the request body
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				c.Abort()
				return
			}

			// Decode JSON
			if err := json.Unmarshal(body, &reqData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
				c.Abort()
				return
			}

			// Inject User ID into request data
			reqData["userid"] = claims.UserId

			// Encode the modified body back to JSON
			updatedBody, err := json.Marshal(reqData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode request body"})
				c.Abort()
				return
			}

			// Replace request body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(updatedBody))
			c.Request.ContentLength = int64(len(updatedBody))
		}

		// Check User ID in Query Params for GET & DELETE requests
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodDelete {
			queryUserId := c.Query("userid")
			logger.Log.Println("Query User ID: ", queryUserId)
			if queryUserId != "" {
				userId, _ := strconv.ParseInt(queryUserId, 10, 64)
				logger.Log.Println("User ID: ", userId)
				logger.Log.Println("Claims User ID: ", claims.UserId)
				if userId != claims.UserId {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
					c.Abort()
					return
				}
			}
		}

		// Store user ID in context for handlers to access
		c.Set("userid", claims.UserId)

		// Continue processing request
		c.Next()
	}
}

func ThrowBlankResponse(c *gin.Context) {
	ThrowJSONResponse(c, BlankPathCheckResponse())
}

// ThrowJSONResponse sends a structured JSON response
func ThrowJSONResponse(c *gin.Context, response entities.APIResponse) {
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internal Server Error") // Fix spelling from "Internel" to "Internal"
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(jsonResponse)
}

// BlankPathCheckResponse returns a standard 404 error response
func BlankPathCheckResponse() entities.APIResponse {
	logger.Log.Println("Blank request called")
	return entities.APIResponse{
		Status:  false,
		Message: "404 Not Found!",
	}
}

func NotPostMethodResponse() entities.APIResponse {
	var response = entities.APIResponse{}
	response.Status = false
	response.Message = "405 method not allowed."
	return response
}
