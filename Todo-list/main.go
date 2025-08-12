package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"context"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const(
	hostname string ="localhost:27017"
	dabName string="demo_todo"
	collectionName string="todo"
	port string="9000"
)

type(
	todoModel struct{
		ID bson.ObjectId `bson:"_id,omitempty"`
		Title string `bson:"title"`
		Completed bool `bson:"completed"`
		createdAt time.Time `bson:"created_at"`
	}

	todo struct{
		ID string `json:"id"`
		Title string `json:"title"`
		Completed bool `json:"completed"`
		CreatedAt string `json:"created_at"`
	}

)


func init(){
	rnd=renderer.New()
	sess,err:= mgo.Dial(hostname)
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
}

// checkErr logs the error and exits if err is not nil
func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main(){
	stopChan:=make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r:=chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/",homeHandler)
	r.Mount("/todo", todoHandlers)

	srv := http.Server{
		Addr:port,
		Handler:r,
		ReadTimeout:60*time.Second,
		WriteTimeout:60*time.Second,
		IdleTimeout:60*time.Second,
	}
	go func(){
		log.Println("Listening on port",port)
		if err:=srv.ListenAndServe(); err!=nil{
			log.Printf("listen:%s\n",err)
		}
	
	}()

	<- stopChan
	log.Println("Shutting down server...")
	ctx,cancel:=context.WithTimeout(context,Background(),5*time.Second)
	srv.Shutdown(ctx)
	defer cancel(
		log.Println("Sever gracefully stopped")
	)
}

func todoHandlers()http.Handler{
	rg:=chi.NewRouter()
	rg.Group(func(r chi.Router){
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Get("/{id}", getTodoByID)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
}

func homeHandler(w http.ResponseWriter,r *http.Request){
	err:= rnd.Template(w,http.StatusOK,[]string{"static/home.tpl"},nil)
	checkErr(err)
}


// Fetches the todo items from the DB
func fetchTodos(w http.ResponseWriter, r *http.Request){
	todos:=[]todoModel{}

	if err:=db.C(collectionName).Find(bson.M{}).All(&todos); err!=nil{
		rnd.JSON(w,http.StatusProcessing,renderer.M{
			"message":"Failed to fetch todos",
			"error":err,
		})
		return
	}

	todoList :=[]todo{}

	for _,t:=range todos{
		todoList=append(todoList,todo{
			ID: t.ID.Hex(),
			Title:t.Title,
			Completed: t.Completed,
			CreatedAt: t.createdAt,
		})
	}
	rnd.JSON(w,http.StatusOK,renderer.M{
		"data": todoList,
	})
}

func createTodo(w http.ReqestResponse, r* http.Request){
	var t todo
	if err:=json.NewDecoder(r.Body).Decode(&t); err!=nil{
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title ==""{
		rnd.JSON(w,http.StatusBadRequest,renderer.M{
		})
		return
	}

	tm:=todoModel{
		ID:bson.NewObjectId(),
		Title:t.Title,
		Completed:false,
		CreatedAt:time.Now(),
	}

	if err:=db.C(collectionName).Insert(tm); err!=nil{
		rnd.JSON(w,http.StatusProcessing,renderer.M{
			"message":"Failed to save todo",
			"error":err,
		})
			return

	}
	rnd.JSON(w,http.StatusCreated,renderer.M{
		"message":"Todo created successfully",
		"todo_id":tm.ID.Hex(),
	})
}