package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	treblle "github.com/treblle/treblle-go"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users  = make(map[int]User)
	nextID = 1
	mu     sync.Mutex
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || users[id].ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users[id])
}

func createUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user.ID = nextID
	nextID++
	users[user.ID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || users[id].ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedUser.ID = id
	users[id] = updatedUser

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || users[id].ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Configure Treblle with additional fields to mask
	treblle.Configure(treblle.Configuration{
		APIKey:                 "***REMOVED***",
		ProjectID:              "***REMOVED***",
		AdditionalFieldsToMask: []string{"bank_account", "routing_number", "tax_id", "auth_token", "ssn", "api_key", "password", "credit_card"},
	})

	r := mux.NewRouter()

	// Create a subrouter with the /api/v1 prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	// Apply Treblle middleware to the main router
	r.Use(treblle.Middleware)

	// Define routes with the API subrouter, wrapping handlers with treblle.WithRoutePath
	api.Handle("/users", treblle.WithRoutePath("GET /api/v1/users", http.HandlerFunc(getUsers))).Methods("GET")
	api.Handle("/users/{id:[0-9]+}", treblle.WithRoutePath("GET /api/v1/users/{id}", http.HandlerFunc(getUser))).Methods("GET")
	api.Handle("/users", treblle.WithRoutePath("POST /api/v1/users", http.HandlerFunc(createUser))).Methods("POST")
	api.Handle("/users/{id:[0-9]+}", treblle.WithRoutePath("PUT /api/v1/users/{id}", http.HandlerFunc(updateUser))).Methods("PUT")
	api.Handle("/users/{id:[0-9]+}", treblle.WithRoutePath("DELETE /api/v1/users/{id}", http.HandlerFunc(deleteUser))).Methods("DELETE")

	log.Println("Server running on port 8085")
	log.Fatal(http.ListenAndServe(":8085", r))
}

//to run the code
//go run main.go
//ngrok http 8085
// Change the ngrok URL to the one provided by ngrok
// Post command ( change ngrok URL in the curl command as well)
// curl -X POST https://d4aa-62-163-212-117.ngrok-free.app/api/v1/users \
// -H "Content-Type: application/json" \
// -d '{"name":"Vredran Cindric","email":"vcindric@example.com","bank_account":"123456789","routing_number":"021000021","tax_id":"12-3456789"}'
// Get command ( change ngrok URL in the curl command as well)
// curl -X GET https://d4aa-62-163-212-117.ngrok-free.app/api/v1/users/1 \
// -H "Content-Type: application/json"