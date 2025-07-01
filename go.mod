module github.com/TheEntropyCollective/randomfs-web

go 1.21

require (
	github.com/TheEntropyCollective/randomfs-core v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
)

replace github.com/TheEntropyCollective/randomfs-core => ../randomfs-core

// This module contains static web assets and does not require Go dependencies
