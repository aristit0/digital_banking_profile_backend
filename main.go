package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/gin-gonic/gin"
)

// Configuration
const (
	CouchbaseHost     = "db"
	CouchbaseUser     = "admin"
	CouchbasePassword = "T1ku$H1t4m"
	BucketName        = "customer_360"
)

// Request/Response Models
type CustomerRequest struct {
	CustomerIDs []string `json:"customer_id" binding:"required"`
}

type Customer360Response struct {
	CustomerID  string                   `json:"customer_id"`
	Customer    map[string]interface{}   `json:"customer"`
	Address     map[string]interface{}   `json:"address"`
	Contact     map[string]interface{}   `json:"contact"`
	Accounts    []map[string]interface{} `json:"accounts"`
	Deposits    []map[string]interface{} `json:"deposits"`
	Loans       []map[string]interface{} `json:"loans"`
	Cards       []map[string]interface{} `json:"cards"`
	Investments []map[string]interface{} `json:"investments"`
	Segment     map[string]interface{}   `json:"segment"`
	Behavior    map[string]interface{}   `json:"behavior"`
	Preference  map[string]interface{}   `json:"preference"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Couchbase connection
var cluster *gocb.Cluster
var bucket *gocb.Bucket

func initCouchbase() error {
	var err error

	// Connect to cluster
	cluster, err = gocb.Connect(
		fmt.Sprintf("couchbase://%s", CouchbaseHost),
		gocb.ClusterOptions{
			Authenticator: gocb.PasswordAuthenticator{
				Username: CouchbaseUser,
				Password: CouchbasePassword,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	// Get bucket
	bucket = cluster.Bucket(BucketName)

	// Wait until bucket is ready
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return fmt.Errorf("bucket not ready: %w", err)
	}

	log.Println("✓ Connected to Couchbase successfully")
	return nil
}

func getCustomer360(customerID string) (*Customer360Response, error) {
	response := &Customer360Response{
		CustomerID:  customerID,
		Accounts:    []map[string]interface{}{},
		Deposits:    []map[string]interface{}{},
		Loans:       []map[string]interface{}{},
		Cards:       []map[string]interface{}{},
		Investments: []map[string]interface{}{},
	}

	// Get collections
	customersCol := bucket.Scope("demographics").Collection("customers")
	addressesCol := bucket.Scope("demographics").Collection("addresses")
	contactsCol := bucket.Scope("demographics").Collection("contacts")
	segmentsCol := bucket.Scope("analytics").Collection("segments")
	behaviorsCol := bucket.Scope("analytics").Collection("behaviors")
	preferencesCol := bucket.Scope("analytics").Collection("preferences")

	// 1. Get customer profile
	customerKey := fmt.Sprintf("customer::%s", customerID)
	getResult, err := customersCol.Get(customerKey, nil)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	var customer map[string]interface{}
	if err := getResult.Content(&customer); err != nil {
		return nil, err
	}
	response.Customer = customer

	// 2. Get address
	addressKey := fmt.Sprintf("address::%s::residential", customerID)
	if getResult, err := addressesCol.Get(addressKey, nil); err == nil {
		var address map[string]interface{}
		if err := getResult.Content(&address); err == nil {
			response.Address = address
		}
	}

	// 3. Get contact
	contactKey := fmt.Sprintf("contact::%s", customerID)
	if getResult, err := contactsCol.Get(contactKey, nil); err == nil {
		var contact map[string]interface{}
		if err := getResult.Content(&contact); err == nil {
			response.Contact = contact
		}
	}

	// 4. Get accounts
	query := fmt.Sprintf("SELECT * FROM `%s`.`products`.`accounts` WHERE customer_id = $1", BucketName)
	rows, err := cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID},
	})
	if err == nil {
		for rows.Next() {
			var account map[string]interface{}
			if err := rows.Row(&account); err == nil {
				if accountData, ok := account["accounts"].(map[string]interface{}); ok {
					response.Accounts = append(response.Accounts, accountData)
				}
			}
		}
		rows.Close()
	}

	// 5. Get deposits
	query = fmt.Sprintf("SELECT * FROM `%s`.`products`.`deposits` WHERE customer_id = $1", BucketName)
	rows, err = cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID},
	})
	if err == nil {
		for rows.Next() {
			var deposit map[string]interface{}
			if err := rows.Row(&deposit); err == nil {
				if depositData, ok := deposit["deposits"].(map[string]interface{}); ok {
					response.Deposits = append(response.Deposits, depositData)
				}
			}
		}
		rows.Close()
	}

	// 6. Get loans
	query = fmt.Sprintf("SELECT * FROM `%s`.`products`.`loans` WHERE customer_id = $1", BucketName)
	rows, err = cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID},
	})
	if err == nil {
		for rows.Next() {
			var loan map[string]interface{}
			if err := rows.Row(&loan); err == nil {
				if loanData, ok := loan["loans"].(map[string]interface{}); ok {
					response.Loans = append(response.Loans, loanData)
				}
			}
		}
		rows.Close()
	}

	// 7. Get cards
	query = fmt.Sprintf("SELECT * FROM `%s`.`products`.`cards` WHERE customer_id = $1", BucketName)
	rows, err = cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID},
	})
	if err == nil {
		for rows.Next() {
			var card map[string]interface{}
			if err := rows.Row(&card); err == nil {
				if cardData, ok := card["cards"].(map[string]interface{}); ok {
					response.Cards = append(response.Cards, cardData)
				}
			}
		}
		rows.Close()
	}

	// 8. Get investments
	query = fmt.Sprintf("SELECT * FROM `%s`.`products`.`investments` WHERE customer_id = $1", BucketName)
	rows, err = cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID},
	})
	if err == nil {
		for rows.Next() {
			var investment map[string]interface{}
			if err := rows.Row(&investment); err == nil {
				if investmentData, ok := investment["investments"].(map[string]interface{}); ok {
					response.Investments = append(response.Investments, investmentData)
				}
			}
		}
		rows.Close()
	}

	// 9. Get segment
	segmentKey := fmt.Sprintf("segment::%s", customerID)
	if getResult, err := segmentsCol.Get(segmentKey, nil); err == nil {
		var segment map[string]interface{}
		if err := getResult.Content(&segment); err == nil {
			response.Segment = segment
		}
	}

	// 10. Get behavior
	behaviorKey := fmt.Sprintf("behavior::%s", customerID)
	if getResult, err := behaviorsCol.Get(behaviorKey, nil); err == nil {
		var behavior map[string]interface{}
		if err := getResult.Content(&behavior); err == nil {
			response.Behavior = behavior
		}
	}

	// 11. Get preference
	preferenceKey := fmt.Sprintf("preference::%s", customerID)
	if getResult, err := preferencesCol.Get(preferenceKey, nil); err == nil {
		var preference map[string]interface{}
		if err := getResult.Content(&preference); err == nil {
			response.Preference = preference
		}
	}

	return response, nil
}

// API Handlers
func handleGetCustomer360(c *gin.Context) {
	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "customer_id array is required",
		})
		return
	}

	if len(req.CustomerIDs) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "at least one customer_id is required",
		})
		return
	}

	// Get data for all requested customers
	results := make([]interface{}, 0)
	errors := make([]map[string]string, 0)

	for _, customerID := range req.CustomerIDs {
		data, err := getCustomer360(customerID)
		if err != nil {
			errors = append(errors, map[string]string{
				"customer_id": customerID,
				"error":       err.Error(),
			})
		} else {
			results = append(results, data)
		}
	}

	response := gin.H{
		"success": len(errors) == 0,
		"count":   len(results),
		"data":    results,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "customer-360-api",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func handleGetStats(c *gin.Context) {
	// Get basic statistics
	query := fmt.Sprintf("SELECT COUNT(*) as total_customers FROM `%s`.`demographics`.`customers`", BucketName)
	rows, err := cluster.Query(query, nil)

	stats := gin.H{
		"service": "customer-360-api",
		"time":    time.Now().Format(time.RFC3339),
	}

	if err == nil && rows.Next() {
		var result map[string]interface{}
		if err := rows.Row(&result); err == nil {
			stats["total_customers"] = result["total_customers"]
		}
		rows.Close()
	}

	c.JSON(http.StatusOK, stats)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	// CORS configuration - Allow all origins for development
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", handleHealthCheck)
		api.GET("/stats", handleGetStats)
		api.POST("/customers", handleGetCustomer360)
	}

	return router
}

func main() {
	// Initialize Couchbase
	if err := initCouchbase(); err != nil {
		log.Fatalf("Failed to initialize Couchbase: %v", err)
	}
	defer cluster.Close(nil)

	// Setup router
	router := setupRouter()

	// Start server
	port := "2113"
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📍 API endpoint: http://localhost:%s/api/v1/customers", port)
	log.Printf("💊 Health check: http://localhost:%s/api/v1/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
