package runpod

import "fmt"

type NetworkVolume struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	DataCenterId string `json:"dataCenterId"`
}

type ErrNoNetworkVolumeFound struct {
	Id string
}

func (e *ErrNoNetworkVolumeFound) Error() string {
	return fmt.Sprintf("no network volume found with id '%s'", e.Id)
}

type ErrNetworkVolumeDataCenterMismatch struct {
	NetworkVolumeId string
	DataCenterId    string
}

func (e *ErrNetworkVolumeDataCenterMismatch) Error() string {
	return fmt.Sprintf(
		`network volume with id '%s' doesn't belong to datacenter with id '%s'; 
if you are trying to create a pod, consider setting config.DataCenterId to empty string 
to automatically detect the correct dataCenterId`,
		e.NetworkVolumeId, e.DataCenterId)
}
