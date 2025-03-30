package runpod

import (
	"fmt"

	"strings"

	"github.com/mattiapavese/go-runpod/runpod/mutations"
)

type CreatePodConfig struct {
	PodName             string
	ImageName           string
	GPUTypeId           string
	CloudType           string
	SupportPublicIP     bool
	StartSSH            bool
	GPUCount            int
	VolumeInGB          int
	MinVCPUCount        int
	MinMemoryInGB       int
	DockerArgs          string
	VolumeMountPath     string
	NetworkVolumeId     string            //nullable (set this to attach existing volume)
	ContainerDiskInGB   int               //nullable
	Ports               string            //nullable
	DataCenterId        string            //nullable
	CountryCode         string            //nullable
	Env                 map[string]string //nullable
	TemplateId          string            //nullable
	AllowedCUDAVersions []string          //nullable
	MinDownload         int               //nullable
	MinUpload           int               //nullable
}

func NewCreatePodConfig() *CreatePodConfig {
	return &CreatePodConfig{
		CloudType:       DefaultCloudType,
		SupportPublicIP: true,
		StartSSH:        true,
		GPUCount:        1,
		MinVCPUCount:    1,
		MinMemoryInGB:   1,
		VolumeMountPath: "/runpod-volume",
		Ports:           "80/http,22/tcp",
	}
}

func (c *Client) CreatePod(config *CreatePodConfig) (*Pod, error) {

	if config.ImageName == "" && config.TemplateId == "" {
		return nil, fmt.Errorf("one between 'config.ImageName' and 'config.TemplateId' id required")
	}

	if config.PodName == "" {
		return nil, fmt.Errorf("field 'config.PodName' is required")
	}

	if config.GPUTypeId == "" {
		return nil, fmt.Errorf("field 'config.GPUTypeId' is required")
	}

	_, err := c.GetGpu(config.GPUTypeId, 1)
	if err != nil {
		return nil, err
	}

	if !(config.CloudType == CloudTypeSecure ||
		config.CloudType == CloudTypeCommunity ||
		config.CloudType == CloudTypeAll) {
		return nil, fmt.Errorf("config.CloudType must be one of 'ALL', 'SECURE' or 'COMMUNITY'")
	}

	if config.NetworkVolumeId != "" {

		user, err := c.GetUser(false)
		if err != nil {
			return nil, err
		}

		foundVolume := false

		for _, volume := range user.NetworkVolumes {
			if config.NetworkVolumeId == volume.Id {
				if config.DataCenterId == "" {
					config.DataCenterId = volume.DataCenterId
				} else {
					if config.DataCenterId != volume.DataCenterId {
						return nil, &ErrNetworkVolumeDataCenterMismatch{
							NetworkVolumeId: config.NetworkVolumeId,
							DataCenterId:    config.DataCenterId,
						}
					}
				}
				foundVolume = true
				break
			}
		}

		if !foundVolume {
			return nil, &ErrNoNetworkVolumeFound{Id: config.NetworkVolumeId}
		}
	}

	if config.ContainerDiskInGB == 0 && config.TemplateId == "" {
		config.ContainerDiskInGB = 10
	}

	mutation := BuildFindAndDeployOnDemandMutation(config)

	resp, err := c.query(mutation)

	if err != nil {
		return nil, err
	}

	var wrapper = struct {
		Pod Pod `json:"podFindAndDeployOnDemand"`
	}{}

	err = c.mapToStruct(resp.Data, &wrapper)
	if err != nil {
		return nil, err
	}

	return &wrapper.Pod, nil

}

func BuildFindAndDeployOnDemandMutation(config *CreatePodConfig) string {

	inputs := `name: "%s", gpuTypeId: "%s", cloudType: %s`
	inputs = fmt.Sprintf(
		inputs, config.PodName, config.GPUTypeId, config.CloudType, // cloudType validate in callee method CreatePod
	)

	// copied from python sdk,
	// by default there was no backup on startSsh: false
	// likely the gql layer automatically handles this
	if config.StartSSH {
		inputs += ", startSsh: true"
	}

	if config.SupportPublicIP {
		inputs += ", supportPublicIp: true"
	} else {
		inputs += ", supportPublicIp: false"
	}

	// optional fields

	if config.DataCenterId != "" {
		inputs += fmt.Sprintf(`, dataCenterId: "%s"`, config.DataCenterId)
	}

	if config.CountryCode != "" {
		inputs += fmt.Sprintf(`, countryCode: "%s"`, config.CountryCode)
	}

	if config.GPUCount != 0 {
		inputs += fmt.Sprintf(`, gpuCount: %d`, config.GPUCount)
	}

	// can be 0 (?)
	inputs += fmt.Sprintf(`, volumeInGb: %d`, config.VolumeInGB)

	// should be never 0 because validated in the callee
	inputs += fmt.Sprintf(`, containerDiskInGb: %d`, config.ContainerDiskInGB)

	if config.MinVCPUCount != 0 {
		inputs += fmt.Sprintf(`, minVcpuCount: %d`, config.MinVCPUCount)
	}

	if config.MinMemoryInGB != 0 {
		inputs += fmt.Sprintf(`, minMemoryInGb: %d`, config.MinMemoryInGB)
	}

	inputs += fmt.Sprintf(`, dockerArgs: "%s"`, config.DockerArgs)

	if config.Ports != "" {
		ports := strings.ReplaceAll(config.Ports, " ", "")
		inputs += fmt.Sprintf(`, ports: "%s"`, ports)
	}

	if config.VolumeMountPath != "" {
		inputs += fmt.Sprintf(`, volumeMountPath: "%s"`, config.VolumeMountPath)
	}

	if config.NetworkVolumeId != "" {
		inputs += fmt.Sprintf(`, networkVolumeId: "%s"`, config.NetworkVolumeId)
	}

	if config.ImageName != "" {
		inputs += fmt.Sprintf(`, imageName: "%s"`, config.ImageName)
	}

	if config.TemplateId != "" {
		inputs += fmt.Sprintf(`, templateId: "%s"`, config.TemplateId)
	}

	if len(config.Env) > 0 {
		var envPairs []string
		for key, value := range config.Env {
			envPairs = append(envPairs, fmt.Sprintf(`{ key: "%s", value: "%s" }`, key, value))
		}
		inputs += fmt.Sprintf(`, env: [%s]`, strings.Join(envPairs, ", "))
	}

	if len(config.AllowedCUDAVersions) > 0 {
		var cudaVersions []string
		for _, version := range config.AllowedCUDAVersions {
			cudaVersions = append(cudaVersions, fmt.Sprintf(`"%s"`, version))
		}
		inputs += fmt.Sprintf(`, allowedCudaVersions: [%s]`, strings.Join(cudaVersions, ", "))
	}

	if config.MinDownload != 0 {
		inputs += fmt.Sprintf(`, minDownload: %v`, config.MinDownload)
	}

	if config.MinUpload != 0 {
		inputs += fmt.Sprintf(`, minUpload: %v`, config.MinUpload)
	}

	inputs = strings.Replace(inputs, "\n", "", -1)

	return fmt.Sprintf(
		mutations.MutationFindAndDeployOnDemand,
		inputs,
	)

}
