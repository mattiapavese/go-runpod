package runpod

import (
	"fmt"

	"github.com/mattiapavese/go-runpod/runpod/queries"
)

type Gpu struct {
	ID                 string      `json:"id"`
	DisplayName        string      `json:"displayName"`
	Manufacturer       string      `json:"manufacturer"`
	MemoryInGb         int         `json:"memoryInGb"`
	CudaCores          int         `json:"cudaCores"`
	SecureCloud        bool        `json:"secureCloud"`
	CommunityCloud     bool        `json:"communityCloud"`
	SecurePrice        float64     `json:"securePrice"`
	CommunityPrice     float64     `json:"communityPrice"`
	OneMonthPrice      float64     `json:"oneMonthPrice"`
	ThreeMonthPrice    float64     `json:"threeMonthPrice"`
	OneWeekPrice       float64     `json:"oneWeekPrice"`
	CommunitySpotPrice float64     `json:"communitySpotPrice"`
	SecureSpotPrice    float64     `json:"secureSpotPrice"`
	LowestPrice        LowestPrice `json:"lowestPrice"`
}

type LowestPrice struct {
	MinimumBidPrice      float64 `json:"minimumBidPrice"`
	UninterruptablePrice float64 `json:"uninterruptablePrice"`
}

type ErrNoGpuFound struct {
	GpuId string
}

func (e *ErrNoGpuFound) Error() string {
	msg :=
		`no GPU found with the specified Id: %s, run client.GetGpus() to get a list of all GPUs available`

	return fmt.Sprintf(msg, e.GpuId)
}

func (c *Client) GetGpu(gpuId string, gpuCount int) (*Gpu, error) {

	q := fmt.Sprintf(queries.QueryGpuType, gpuId, gpuCount)
	resp, err := c.query(q)

	if err != nil {
		fmt.Println("err1")
		return nil, err
	}

	var wrapper = struct {
		GpuTypes []Gpu `json:"gpuTypes"`
	}{}

	err = c.mapToStruct(resp.Data, &wrapper)
	if err != nil {
		fmt.Println("err2")
		return nil, err
	}

	if len(wrapper.GpuTypes) == 0 {
		return nil, &ErrNoGpuFound{GpuId: gpuId}
	}

	return &wrapper.GpuTypes[0], nil
}

type GpuCondensed struct {
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
	MemoryInGb  int    `json:"memoryInGb"`
}

func (c *Client) GetGpus() ([]GpuCondensed, error) {

	resp, err := c.query(queries.QueryGpuTypes)

	if err != nil {
		return nil, err
	}

	var wrapper = struct {
		GpuTypes []GpuCondensed `json:"gpuTypes"`
	}{}

	err = c.mapToStruct(resp.Data, &wrapper)
	if err != nil {
		return nil, err
	}

	var gpus []GpuCondensed
	// ensure the gpus with 'unknown' id are not included
	for _, gpu := range wrapper.GpuTypes {
		if gpu.Id != "unknown" {
			gpus = append(gpus, gpu)
		}
	}

	return gpus, nil
}
