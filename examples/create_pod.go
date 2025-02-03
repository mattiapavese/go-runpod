package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/mattiapavese/go-runpod/runpod"
)

func ExampleCreatePod() {

	client := runpod.NewClient(os.Getenv("RUNPOD_API_KEY"))

	config := runpod.NewCreatePodConfig()
	config.PodName = "my_pod_test"
	config.ImageName = "runpod/pytorch:2.2.0-py3.10-cuda12.1.1-devel-ubuntu22.04"
	config.GPUTypeId = "NVIDIA RTX 2000 Ada Generation"
	config.Ports = config.Ports + ",8888/http,8000/http" //to add ports to the default ports
	config.NetworkVolumeId = "4h1a4hsiju"

	pod, err := client.CreatePod(config)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(pod)
}
