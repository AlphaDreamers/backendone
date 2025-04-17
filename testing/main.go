package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	storage_go "github.com/supabase-community/storage-go"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}
	secret := os.Getenv("SECRET_ACCESS_KEY")
	storageClient := storage_go.NewClient(UrlConstructor(), secret, nil)
	result, err := storageClient.CreateBucket(uuid.New().String(), storage_go.BucketOptions{
		Public: true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created bucket: %s", result.Name)
}

func UrlConstructor() string {
	project_reference_id := os.Getenv("PROJECT_ID")
	return fmt.Sprintf("https://%s.supabase.co/storage/v1", project_reference_id)
}
