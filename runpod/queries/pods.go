package queries

var QueryPods = `
query myPods {
    myself {
        pods {
            id
            containerDiskInGb
            costPerHr
            desiredStatus
            dockerArgs
            dockerId
            env
            gpuCount
            imageName
            lastStatusChange
            machineId
            memoryInGb
            name
            podType
            port
            ports
            uptimeSeconds
            vcpuCount
            volumeInGb
            volumeMountPath
            runtime {
                ports{
                    ip
                    isIpPublic
                    privatePort
                    publicPort
                    type
                }
            }
            machine {
                gpuDisplayName
            }
        }
    }
}`

// format with podId
//
//	// example
//	var query := fmt.Sprintf(queries.QueryPod, "Abd89XpqRsj3")
var QueryPod = `
query pod {
	pod(input: {podId: "%s"}) {
		id
		containerDiskInGb
		costPerHr
		desiredStatus
		dockerArgs
		dockerId
		env
		gpuCount
		imageName
		lastStatusChange
		machineId
		memoryInGb
		name
		podType
		port
		ports
		uptimeSeconds
		vcpuCount
		volumeInGb
		volumeMountPath
		runtime {
			ports {
				ip
				isIpPublic
				privatePort
				publicPort
				type
			}
		}
		machine {
			gpuDisplayName
		}
	}
}`
