package mutations

// this is the mutation that follows
// the schema in the python sdk

// format with the input string
// containing all input parameters
// var MutationCreatePod = `
// mutation {
// 	podFindAndDeployOnDemand(input: {%s}){
// 		id
// 		desiredStatus
// 		imageName
// 		env
// 		machineId
// 		machine {
// 			podHostId
// 		}
// 	}
// }
// `

// this version of the mutation is different from python sdk and returns
// the whole Pod object information
// Difference is just in return values taht returns the whole Pod info

var podReturn = `
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
`

// format with the formatted input string
// containing all input parameters
//
//	// example
//	var query := fmt.Sprintf(mutations.MutationFindAndDeployOnDemand, inputStr)
var MutationFindAndDeployOnDemand = `
mutation{podFindAndDeployOnDemand(input:{%s}){` + podReturn + `}}`

// format with podId (string)
var MutationStopPod = `mutation{podStop(input:{podId:"%s" }){` + podReturn + `}}`

// format with podId (string) and gpuCount (integer)
var MutationResumePod = `mutation{podResume(input:{podId:"%s",gpuCount:%d}){` + podReturn + `}}`

// format with podId (string)
var MutationPodTerminate = `mutation{podTerminate(input:{podId:"%s"})}`
