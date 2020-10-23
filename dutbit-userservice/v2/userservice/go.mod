module service

go 1.13

replace models => ./models

replace handlers => ./handlers

replace db => ./db

replace utils => ./utils

require (
	db v0.0.0-00010101000000-000000000000 // indirect
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/wonderivan/logger v1.0.0
	go.mongodb.org/mongo-driver v1.4.2 // indirect
	handlers v0.0.0-00010101000000-000000000000
	models v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)
