package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/mattiapavese/go-runpod/runpod"
)

func ExampleGetUser() {

	client := runpod.NewClient(os.Getenv("RUNPOD_API_KEY"))

	user, err := client.GetUser(false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(user)
}
