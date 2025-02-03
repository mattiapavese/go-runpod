package runpod

import (
	"fmt"

	"github.com/mattiapavese/go-runpod/runpod/mutations"
)

func (c *Client) TerminatePodFromId(podId string) (*Pod, error) {
	q := fmt.Sprintf(mutations.MutationPodTerminate, podId)

	_, err := c.query(q)
	if err != nil {
		return nil, err
	}

	return &Pod{DesiredStatus: "TERMINATED", Id: podId}, nil
}

func (c *Client) TerminatePod(pod *Pod) (*Pod, error) {
	return c.TerminatePodFromId(pod.Id)
}

func (c *Client) ResumePodFromId(podId string, gpuCount int) (*Pod, error) {

	q := fmt.Sprintf(mutations.MutationResumePod, podId, gpuCount)

	resp, err := c.query(q)
	if err != nil {
		return nil, err
	}

	var wrapper = struct {
		Pod Pod `json:"podResume"`
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

func (c *Client) ResumePod(pod *Pod, gpuCount int) (*Pod, error) {
	if gpuCount == 0 {
		return c.ResumePodFromId(pod.Id, pod.GPUCount)
	}
	return c.ResumePodFromId(pod.Id, gpuCount)
}

func (c *Client) StopPodFromId(podId string) (*Pod, error) {
	q := fmt.Sprintf(mutations.MutationStopPod, podId)

	resp, err := c.query(q)
	if err != nil {
		return nil, err
	}

	var wrapper = struct {
		Pod Pod `json:"podStop"`
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

func (c *Client) StopPod(pod *Pod) (*Pod, error) {
	return c.StopPodFromId(pod.Id)
}
