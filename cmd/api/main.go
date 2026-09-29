package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miriyaladhanwinn/on-duty/internal/api"
	"github.com/miriyaladhanwinn/on-duty/internal/db"
)

const Banner = `
   ____          ____        _         
  / __ \____    / __ \__  __/ /___  __ 
 / / / / __ \  / / / / / / / __/ / / / 
/ /_/ / / / / / /_/ / /_/ / /_/ /_/ /  
\____/_/ /_/ /_____/\__,_/\__/\__, /   
                             /____/    
 Amrita CSE On-Duty & Attendance Workflow Engine
`

func main() {
	fmt.Print(Banner)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize store
	store := db.NewMemoryStore()
	server := api.NewServer(store)

	log.Printf("[+] Amrita CSE On-Duty REST API server listening on :%s\n", port)
	log.Printf("[+] Swagger/API Docs: http://localhost:%s/api/v1/health\n", port)
	if err := server.Router().Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
