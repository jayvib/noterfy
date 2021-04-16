package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
)

var payload = `
{
	"id": "ca727945-430b-4bf9-8401-e60a2179a6af"
}
`

type request struct {
	ID uuid.UUID `json:"id"`
}

func main() {

	var req request
	err := json.Unmarshal([]byte(payload), &req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(req.ID)

}
