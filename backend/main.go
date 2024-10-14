package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Book struct (Model)
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title"`
	Author string             `json:"author"`
	Price  float64            `json:"price"`
}

var bookCollection *mongo.Collection

// Connect to MongoDB
func initMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}
	bookCollection = client.Database("library").Collection("books")
	log.Println("Connected to MongoDB!")
}

// Create a new book
func createBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Generate a new MongoDB ObjectID
	newBook.ID = primitive.NewObjectID()

	// Insert the book into the database
	_, err := bookCollection.InsertOne(context.Background(), newBook)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating book"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"id": newBook.ID.Hex()}) // Return the ID as hex string
}

// Get all books
func getBooks(c *gin.Context) {
	var books []Book

	cursor, err := bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error fetching books"})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book Book
		err := cursor.Decode(&book)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error decoding book"})
			return
		}
		books = append(books, book)
	}

	c.IndentedJSON(http.StatusOK, books)
}

// Get a single book by ID
func getBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid book ID"})
		return
	}

	var book Book
	err = bookCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&book)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// Update a book by ID
func updateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid book ID"})
		return
	}

	var updatedBook Book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Update the book in the database
	_, err = bookCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": updatedBook})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error updating book"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book updated"})
}

// Delete a book by ID
func deleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid book ID"})
		return
	}

	_, err = bookCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error deleting book"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func main() {
	// Initialize MongoDB connection
	initMongoDB()

	// Create a new Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Define routes
	r.POST("/books", createBook)
	r.GET("/books", getBooks)
	r.GET("/books/:id", getBookByID)
	r.PUT("/books/:id", updateBook)
	r.DELETE("/books/:id", deleteBook)

	// Start the Gin server on port 8080
	r.Run(":8080")
}
