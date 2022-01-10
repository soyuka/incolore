package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DB string
	ShortenerHostname string
	IdAlphabet        string
	IdLength          int
	Port              string
	Directory         string
	MaxSize           int64
}

func GetConfig() Config {
	dbPath := os.Getenv("INCOLORE_DB")

	if dbPath == "" {
		dbPath = "bolt://data.bolt"
	}

	shortenerHostname := os.Getenv("INCOLORE_HOSTNAME")

	if shortenerHostname == "" {
		shortenerHostname = "http://localhost:5377"
	}

	port := os.Getenv("INCOLORE_PORT")
	if port == "" {
		port = "5377"
	}

	idLength, err := strconv.ParseInt(os.Getenv("INCOLORE_ID_LENGTH"), 10, 32)

	if idLength == 0 || err != nil {
		idLength = 12
	}

	idAlphabet := os.Getenv("INCOLORE_ID_ALPHABET")

	if idAlphabet == "" {
		idAlphabet = "0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
	}

	uploadDirectory := os.Getenv("INCOLORE_DIRECTORY") 

	if uploadDirectory == "" {
		uploadDirectory = "./upload"
	}

	maxSize, err := strconv.ParseInt(os.Getenv("INCOLORE_MAX_SIZE"), 10, 32)

	if maxSize == 0 || err != nil {
		maxSize = 10000000
	}

	// todo: log config
	log.Println("DB Path", dbPath)
	log.Println("Hostname", shortenerHostname)
	log.Println("Directory", uploadDirectory)

	return Config{
		ShortenerHostname: shortenerHostname,
		IdLength:          int(idLength),
		IdAlphabet:        idAlphabet,
		Port:              port,
		DB:                dbPath,
		Directory:         uploadDirectory,
		MaxSize:           maxSize,
	}
}
