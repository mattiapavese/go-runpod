package examples

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/mattiapavese/go-runpod/runpod"
)

func ExampleGetGPus() {

	client := runpod.NewClient(os.Getenv("RUNPOD_API_KEY"))

	gpus, err := client.GetGpus()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Gpu types available")
	for _, gpu := range gpus {
		fmt.Printf("Id: %s | Display Name: %s | VRAM: %d\n", gpu.Id, gpu.DisplayName, gpu.MemoryInGb)
	}

	idx := rand.Intn(len(gpus))
	gpu, err := client.GetGpu(gpus[idx].Id, 1)

	if runpod.ErrorIs(err, &runpod.ErrNoGpuFound{}) { // util to check what type of error
		log.Fatalln(err)
	}

	fmt.Println(gpu)
}
