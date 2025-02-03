package runpod

import "github.com/mattiapavese/go-runpod/runpod/queries"

type User struct {
	Id             string          `json:"id"`
	PubKey         string          `json:"pubKey"`
	NetworkVolumes []NetworkVolume `json:"networkVolumes"`
}

func (c *Client) GetUser(withPubKey bool) (u *User, err error) {

	var resp *GQLResponse
	if withPubKey {
		resp, err = c.query(queries.QueryUser)
	} else {
		resp, err = c.query(queries.QueryUserNoPubKey)
	}

	if err != nil {
		return nil, err
	}

	var wrapper = struct {
		MySelf User `json:"mySelf"`
	}{}

	err = c.mapToStruct(resp.Data, &wrapper)
	if err != nil {
		return nil, err
	}

	return &wrapper.MySelf, nil
}
