package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var db *mongo.Database
var client *mongo.Client

const (
	hostname       string = "localhost:27017"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createdAt"`
	}

	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func init() {
	rnd = renderer.New()

	//create context with timeout
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	//connect to MongoDB
	clientOptions:= options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	client, err = mongo.Connect(ctx, clientOptions)

	if(err !=nil){
		log.Fatalf("Failed to connect to MongoDB: %v",err)
	}

	//Test the connection
	err=client.Ping(ctx,nil)
	if(err!=nil){
		log.Fatalf("Failed to ping MongoDB: %v",err)
	}

	db=client.Database(dbName)
	log.Println("Connected to MongoDB database")
}

// checkErr logs the error and exits if err is not nil
func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen:%s\n", err)
		}

	}()

	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Sever gracefully stopped")

}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		//r.Get("/{id}", getTodoByID)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})

	return rg
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

// Fetches the todo items from the DB
func fetchTodos(w http.ResponseWriter, r *http.Request) {
	collection:= db.Collection(collectionName)
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	cursor,err:=collection.Find(ctx,bson.M{})
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)
	
	var todos [] todoModel
	if err := cursor.All(ctx, &todos); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to decode todos",
			"error":   err.Error(),
		})
	}

	todoList := []todo{}

	for _, t := range todos {
		todoList = append(todoList, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
	})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Invalid Json",
			"error":   err.Error(),
		})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Title cannot be empty",
		})
		return
	}

	tm := todoModel{
		ID:        primitive.NewObjectID(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	collection :=db.Collection(collectionName)
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	result,err:=collection.InsertOne(ctx,tm)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return

	}
	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": result.InsertedID.(primitive.ObjectID).Hex(),
	})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	objectID, err:= primitive.ObjectIDFromHex(id)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	collection:=db.Collection(collectionName)
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	result,err:=collection.DeleteOne(ctx,primitive.M{"_id":objectID})

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":err.Error(),
		})
		return
	}

	if result.DeletedCount==0{
		rnd.JSON(w,http.StatusNotFound,renderer.M{
			"message": "Todo not found",
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo deleted successfully",
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	objectID,err:=primitive.ObjectIDFromHex(id)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to decode request body",
		})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Title cannot be empty",
		})
		return
	}

	collection:=db.Collection(collectionName)
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	update:=bson.M{
		"$set": bson.M{
			"title": t.Title,
			"completed": t.Completed,
		},
	}

	result,err:=collection.UpdateOne(ctx,bson.M{"_id":objectID},update)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err.Error(),
		})
		return
	}


	if result.MatchedCount==0{
		rnd.JSON(w,http.StatusNotFound,renderer.M{
				"message": "Todo not found",
			})
		return
	}
	rnd.JSON(w,http.StatusOK,renderer.M{
		"message":"Todo updated successfully",
	})
}