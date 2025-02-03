package queries

var QueryUser = `
query myself {
    myself {
        id
        pubKey
        networkVolumes {
            id
            name
            size
            dataCenterId
        }
    }
}`

var QueryUserNoPubKey = `
query myself {
    myself {
        id
        networkVolumes {
            id
            name
            size
            dataCenterId
        }
    }
}`
