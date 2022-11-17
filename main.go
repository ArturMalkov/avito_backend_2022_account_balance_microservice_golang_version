package main

import (
	apiV1 "account-balance-microservice/api/v1"
	"account-balance-microservice/storage"
)

func main() {
	db := storage.GetDatabaseConnection()
	// defer db.Close()

	storage.SetupDatabase(db)

	router := apiV1.GetRouter()
	router.Run() // listen and serve on 0.0.0.0:8080
}
