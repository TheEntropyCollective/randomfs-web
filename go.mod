module github.com/TheEntropyCollective/randomfs-web

go 1.21

require (
	github.com/TheEntropyCollective/randomfs-core v0.1.5
	github.com/gorilla/mux v1.8.1
)

replace github.com/TheEntropyCollective/randomfs-core => ../randomfs-core

// This module contains static web assets and does not require Go dependencies
