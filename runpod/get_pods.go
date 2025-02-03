package runpod

import (
	"fmt"

	"github.com/mattiapavese/go-runpod/runpod/queries"
)

type ErrNoPodFound struct {
	Id string
}

func (e *ErrNoPodFound) Error() string {
	return fmt.Sprintf("no pod with id '%s' exists", e.Id)
}

func (c *Client) GetPod(podId string) (*Pod, error) {

	q := fmt.Sprintf(queries.QueryPod, podId)

	resp, err := c.query(q)
	if err != nil {
		return nil, err
	}

	var wrapper = struct {
		Pod Pod `json:"pod"`
	}{}

	err = c.mapToStruct(resp.Data, &wrapper)
	if err != nil {
		return nil, err
	}

	if wrapper.Pod.Id == "" {
		return &wrapper.Pod, &ErrNoPodFound{Id: podId}
	}

	return &wrapper.Pod, nil
}
