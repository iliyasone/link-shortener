package main

import (
	"flag"
	"fmt"
	"link-shortener/internal/storage"
	"os"
)

func main() {
	storageType := flag.String("storage", "ram", "Storage backend to use: 'ram' or 'postgres'")
	flag.Parse()

	var store storage.Storage
	switch *storageType {
	case "postgres":
		{
			fmt.Println("Postgres storage selected, but not implemented yet")
			os.Exit(0)
		}
	case "ram":
		{
			store = storage.NewRAMStorage()
		}
	default:
		{
			fmt.Println("Wrong storage type")
		}
	}

	fmt.Printf("Storage %T created successfully\n", store)
}
