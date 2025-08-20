package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
		Action string `json:"action"`
	} `json:"payload"`
}

func main(){
	if len(os.Args)<2{
		fmt.Println("Usage: GitHub_user_Activity <username>")
		return
	}

	username:=os.Args[1]

	if username==""{
		fmt.Println("Please provide a GitHub username.")
		return
	}

	fmt.Printf("Fetching activity for user: %s\n", username)

	url:= fmt.Sprintf("https://api.github.com/users/%s/events", username)
	
	resp, err := http.Get(url)
	
	if err != nil {
		fmt.Println("Error fetching data: ",err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode!= http.StatusOK {
		fmt.Println("Error: ", resp.Status)
		return
	}

	body,err:=io.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("Error reading response body: ", err)
		return
	}

	var events [] Event
	if err:=json.Unmarshal(body,&events); err!=nil{
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	if len(events)==0{
		fmt.Println("No activity found for user:", username)
		return
	}

	for _, event := range events {

		switch event.Type{
			case "PushEvent":
				fmt.Printf("Pushed %d commits to %s:\n",len(event.Payload.Commits),event.Repo.Name)
		
			case "PullRequestEvent":
				fmt.Printf("Pull request in %s:\n",  event.Repo.Name)

			case "StarrerEvent":
				fmt.Printf("User %s starred %s:\n", username, event.Repo.Name)

			case "ForkEvent":
				fmt.Printf("User %s forked %s:\n", username, event.Repo.Name)
			
			case "opened":
				fmt.Printf("Opened a new issue in %s:\n", event.Repo.Name)

			default:
				fmt.Printf("%s in repo %s\n", event.Type, event.Repo.Name)
			
		}
	}
}