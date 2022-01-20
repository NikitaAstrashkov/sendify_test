package models

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

// InitGormConnection - func that initializing conn to MySQL DB
func InitGormConnection(connectionString string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Failed to open DB connection, error: ", err.Error())
	}
	db.DB().SetConnMaxLifetime(time.Second)
	db.DB().SetMaxIdleConns(4)
	db.DB().SetMaxOpenConns(10)
	return db
}

// LoadEnv - func for retrieving vars from .env
func LoadEnv(appConfig interface{}) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load config, error: ", err.Error())
	}
	err = env.Parse(appConfig)
	if err != nil {
		log.Fatal("Failed to parse config, error: ", err.Error())
	}
}

//PrintHTTPResult - func for printing HTTP results
func PrintHTTPResult(w http.ResponseWriter, resultHTTPCode int, jsonAnswer interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if jsonAnswer == nil && resultHTTPCode == http.StatusInternalServerError {
		body, _ := json.Marshal(map[string]string{"error": "Internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}

	body, err := json.Marshal(jsonAnswer)
	if err != nil {
		log.Print("failed to marshal json, error: ", err.Error())
		body, _ = json.Marshal(map[string]string{"error": "Internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}

	if resultHTTPCode >= 400 {
		body, _ = json.Marshal(map[string]interface{}{"error": jsonAnswer})
		w.WriteHeader(resultHTTPCode)
		w.Write(body)
		return
	}

	w.WriteHeader(resultHTTPCode)
	w.Write(body)
}
