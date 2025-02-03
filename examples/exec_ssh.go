package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/mattiapavese/go-runpod/runpod"
)

func ExecSshOnPod() {
	client := runpod.NewClient(os.Getenv("RUNPOD_API_KEY"))

	podId := "k9t70kydlm5q1n" // your pod's id

	pod, err := client.GetPod(podId)
	if err != nil {
		log.Fatalln(err)
	}

	out, err := pod.Exec("ls /", os.Getenv("PATH_TO_PRIVATE_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(out)
}
