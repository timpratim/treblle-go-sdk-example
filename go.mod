module ngroktest

go 1.23.2

require (
	github.com/gorilla/mux v1.8.1
	github.com/treblle/treblle-go v0.7.2
)

require golang.org/x/sync v0.11.0 // indirect

replace github.com/treblle/treblle-go => github.com/timpratim/treblle-go v0.0.0-20250305182204-c99dda44ed8a
