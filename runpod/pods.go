package runpod

type Pod struct {
	Id                string   `json:"id"`
	ContainerDiskInGb int      `json:"containerDiskInGb"`
	CostPerHr         float64  `json:"costPerHr"`
	DesiredStatus     string   `json:"desiredStatus"`
	DockerArgs        string   `json:"dockerArgs"`
	DockerID          string   `json:"dockerId"`
	Env               []string `json:"env"`
	GPUCount          int      `json:"gpuCount"`
	ImageName         string   `json:"imageName"`
	LastStatusChange  string   `json:"lastStatusChange"`
	MachineId         string   `json:"machineId"`
	MemoryInGb        int      `json:"memoryInGb"`
	Name              string   `json:"name"`
	PodType           string   `json:"podType"`
	Port              int      `json:"port"`
	Ports             string   `json:"ports"`
	UptimeSeconds     int64    `json:"uptimeSeconds"`
	VCPUCount         int      `json:"vcpuCount"`
	VolumeInGb        int      `json:"volumeInGb"`
	VolumeMountPath   string   `json:"volumeMountPath"`
	Runtime           Runtime  `json:"runtime"`
	Machine           Machine  `json:"machine"`
}

type Runtime struct {
	Ports []Port `json:"ports"`
}

type Port struct {
	IP          string `json:"ip"`
	IsIPPublic  bool   `json:"isIpPublic"`
	PrivatePort int    `json:"privatePort"`
	PublicPort  int    `json:"publicPort"`
	Type        string `json:"type"`
}

type Machine struct {
	GPUDisplayName string `json:"gpuDisplayName"`
}
