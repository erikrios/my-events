package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Environment string
	Port        uint16
	DBHost      string
	DBPort      uint16
	DBUsername  string
	DBPassword  string
	DBName      string
)

var envError = errors.New("Failed to get the environment variables. Please check .env.example file for the example.")

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %s\n", err.Error())
	}

	if err := getEnvironments(); err != nil {
		log.Fatalln(err)
	}
}

func getEnvironments() (err error) {

	if v, ok := os.LookupEnv("ENV"); !ok {
		err = envError
		return
	} else {
		Environment = v
	}

	if v, ok := os.LookupEnv("PORT"); !ok {
		err = envError
		return
	} else {
		if convVal, convErr := strconv.Atoi(v); convErr != nil {
			err = envError
			return
		} else {
			Port = uint16(convVal)
		}
	}

	if v, ok := os.LookupEnv("DB_HOST"); !ok {
		err = envError
		return
	} else {
		DBHost = v
	}

	if v, ok := os.LookupEnv("DB_PORT"); !ok {
		err = envError
		return
	} else {
		if convVal, convErr := strconv.Atoi(v); convErr != nil {
			err = envError
			return
		} else {
			DBPort = uint16(convVal)
		}
	}

	if v, ok := os.LookupEnv("DB_USERNAME"); !ok {
		err = envError
		return
	} else {
		DBUsername = v
	}

	if v, ok := os.LookupEnv("DB_PASSWORD"); !ok {
		err = envError
		return
	} else {
		DBPassword = v
	}

	if v, ok := os.LookupEnv("DB_NAME"); !ok {
		err = envError
		return
	} else {
		DBName = v
	}

	return
}
