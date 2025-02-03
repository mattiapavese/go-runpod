package queries

var QueryGpuTypes = `
query GpuTypes {
  gpuTypes {
    id
    displayName
    memoryInGb
  }
}
`

// format with gpuId, gpuCount
//
//	// example
//	var query := fmt.Sprintf(queries.QueryGpuType, "NVIDIA GeForce RTX 4090", 2)
var QueryGpuType = `
query GpuTypes {
      gpuTypes(input: {id: "%s"}) {
        maxGpuCount
        id
        displayName
        manufacturer
        memoryInGb
        cudaCores
        secureCloud
        communityCloud
        securePrice
        communityPrice
        oneMonthPrice
        threeMonthPrice
        oneWeekPrice
        communitySpotPrice
        secureSpotPrice
        lowestPrice(input: {gpuCount: %d}) {
          minimumBidPrice
          uninterruptablePrice
        }
      }
    }
`
