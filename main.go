package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"bytes"
	"encoding/json"
	
	
)

var taskCollection *mongo.Collection
var jwtSecret []byte

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string             `bson:"userId,omitempty" json:"userId,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Completed   bool               `bson:"completed,omitempty" json:"completed,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	Priority    string `json:"priority,omitempty"`
}

func extractJSONFromResponse(response string) (map[string]interface{}, error) {
	// Regex to extract JSON inside the response
	re := regexp.MustCompile("```json\\n(.*?)\\n```")
	match := re.FindStringSubmatch(response)

	if len(match) < 2 {
		return nil, fmt.Errorf("failed to extract JSON")
	}

	// Parse JSON from extracted text
	var parsedJSON map[string]interface{}
	err := json.Unmarshal([]byte(match[1]), &parsedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return parsedJSON, nil
}


// Middleware: JWT Authentication
func AuthMiddleware(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID") // Get userID from the request header
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing user ID"})
	}

	// Attach userID to request context
	c.Locals("userID", userID)
	log.Println("✅ Authenticated User ID:", userID)

	return c.Next()
}

// 📌 Get All Tasks (User-Specific)
func GetTasks(c *fiber.Ctx) error {
	userId :=  c.Get("X-User-ID")// Clerk provides userId in frontend
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID required"})
	}

	var tasks []Task
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}

// 📌 Create Task
func CreateTask(c *fiber.Ctx) error {
	var task Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request", "message": err.Error()})
	}

	// 📌 Get User ID
	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID required"})
	}

	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UserID = userId
	task.Completed = false

	// 🔥 Fetch previous tasks for AI analysis
	var tasks []Task
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		log.Println("Error fetching tasks:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch previous tasks"})
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var existingTask Task
		cursor.Decode(&existingTask)
		tasks = append(tasks, existingTask)
	}

	// 📌 Format past tasks into a list
	taskList := ""
	for _, t := range tasks {
		taskList += fmt.Sprintf("- %s: %s (Priority: %s)\n", t.Title, t.Description, t.Priority)
	}

	// 🔥 Strict JSON Prompt for AI
	prompt := fmt.Sprintf(`
		You are an AI assistant for task management.
		Here are the user's past tasks:

		%s

		The user has created a new task titled: "%s".

		Your task is to Analyze the previous tasks and tell the user whether this new task should be of **High, Medium, or Low Priority**.
		If the new task is **similar to past high-priority tasks**, it should also be **High Priority**.
		If it has **low effort but high outcome**, give it a **higher priority**.:
		- Assign a **priority** (High, Medium, Low).
		- Generate a **short description including reason**.

		🚨 **IMPORTANT:** Respond with ONLY the following JSON format:
		{
			"priority": "<High/Medium/Low>",
			"description": "<Generated description>"
		}
	`, taskList, task.Title)

	// 📌 Call Gemini AI API
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	})
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request body"})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Missing API key"})
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create API request"})
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI suggestion"})
	}
	defer resp.Body.Close()

	// ✅ Parse AI response directly as JSON
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Println("Error decoding response:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse AI response"})
	}

	// 🔥 Extract AI-generated priority and description
	// 🔥 Extract AI-generated priority and description
	if candidates, ok := responseData["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if text, ok := parts[0].(map[string]interface{})["text"].(string); ok {
					// ✅ Clean response: Trim spaces and remove Markdown backticks
					text = strings.TrimSpace(text)
	
					// ✅ Remove ```json ... ``` using regex
					re := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
					matches := re.FindStringSubmatch(text)
					if len(matches) > 1 {
						text = matches[1] // Extract the JSON part
					}
	
					// ✅ Ensure it's valid JSON before parsing
					if strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}") {
						var aiData map[string]string
						if err := json.Unmarshal([]byte(text), &aiData); err == nil {
							task.Priority = aiData["priority"]
							task.Description = aiData["description"]
						} else {
							log.Println("Failed to parse AI-generated JSON:", err)
						}
					} else {
						log.Println("Invalid AI response format after cleaning:", text)
					}
				}
			}
		}
	}


	// ✅ Save Task in MongoDB
	_, err = taskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create task"})
	}

	// ✅ Return the created task
	return c.JSON(task)
}


// 📌 Update Task
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := taskCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update task"})
	}

	return c.JSON(fiber.Map{"message": "Task updated"})
}

// 📌 Delete Task
func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := taskCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete task"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted"})
}


func AI_TaskSuggestions(c *fiber.Ctx) error {
	// 📌 Get user ID from headers (Clerk provides this)
	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID required"})
	}

	// 🔥 Fetch previous tasks from MongoDB
	var tasks []Task
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		log.Println("Error fetching tasks:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch previous tasks"})
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	// 📌 Get user query (New task to prioritize)
	userQuery := c.Query("query")
	if userQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing query parameter"})
	}

	// 📝 Create dynamic prompt with previous tasks
	taskList := ""
	for _, task := range tasks {
		taskList += fmt.Sprintf("- %s: %s (Priority: %s)\n", task.Title, task.Description, task.Priority)
	}

	prompt := fmt.Sprintf(`
		You are an AI assistant that helps in task prioritization based on previous tasks.
		The user has completed the following tasks:

		%s

		Now, the user wants to do the following task: "%s".

		Analyze the previous tasks and tell the user whether this task should be of **High, Medium, or Low Priority**.
		If the new task is **similar to past high-priority tasks**, it should also be **High Priority**.
		If it has **low effort but high outcome**, give it a **higher priority**.
		Respond in this JSON format:
		{
			"task": "<user task>",
			"priority": "<High/Medium/Low>",
			"reason": "<Why this priority is assigned>"
		}
	`, taskList, userQuery)

	// 📌 Prepare request body for Gemini API
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	})
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request body"})
	}

	// 🔥 Gemini API Call
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Missing API key"})
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create API request"})
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI suggestion"})
	}
	defer resp.Body.Close()

	// ✅ Parse response
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Println("Error decoding response:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse AI response"})
	}

	// Extract AI-generated suggestion
	generatedText := ""
	if candidates, ok := responseData["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if text, ok := parts[0].(map[string]interface{})["text"].(string); ok {
					generatedText = text
				}
			}
		}
	}

	if generatedText == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to extract AI-generated text"})
	}

	// ✅ Return AI suggestion
	return c.JSON(fiber.Map{"suggestion": generatedText})
}

func main() {
	_ = godotenv.Load()
	mongoURI := os.Getenv("MONGO_URI")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ MongoDB Connection Failed:", err)
	}

	// Check if the connection is successful
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB is not reachable:", err)
	}

	// Initialize taskCollection
	taskCollection = client.Database("taskmanager").Collection("tasks")
	if taskCollection == nil {
		log.Fatal("❌ taskCollection is nil, MongoDB connection failed.")
	}

	fmt.Println("✅ Connected to MongoDB!")

	// Initialize Fiber app
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("FRONTEND_URL"), // Change this to your Next.js frontend URL
		AllowHeaders: "Origin, Content-Type, Accept, X-User-ID",
	}))

	// Task CRUD Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Backend is running!",
		})
	})

	app.Get("/tasks", AuthMiddleware, GetTasks)
	app.Post("/tasks", CreateTask)
	app.Put("/tasks/:id", AuthMiddleware, UpdateTask)
	app.Delete("/tasks/:id", AuthMiddleware, DeleteTask)

	// AI-Powered Task Suggestions
	app.Get("/ai-suggestions", AuthMiddleware, AI_TaskSuggestions)

	// Start Server
	log.Fatal(app.Listen(":8080"))
}
