package main

import (
	"fmt"
	"net/http"
	"webcrawler/prioritizer"
)


func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if len(header.Filename) < 4 ||
		header.Filename[len(header.Filename)-4:] != ".csv" {

		http.Error(w, "Only CSV files are allowed", http.StatusBadRequest)
		return
	}

	fmt.Println("recieved request!")

	prioritizer.Run(file)
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}