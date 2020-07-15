package fmcClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// NetworkGroups .
type NetworkGroups struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Items []struct {
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Literals []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"literals"`
		Type        string `json:"type"`
		Overridable bool   `json:"overridable"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		Metadata    struct {
			ReadOnly struct {
				State  bool   `json:"state"`
				Reason string `json:"reason"`
			} `json:"UpdateNetworksreadOnly"`
			Timestamp int64 `json:"timestamp"`
			LastUser  struct {
				Name string `json:"name"`
			} `json:"lastUser"`
			Domain struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"domain"`
		} `json:"metadata"`
	} `json:"items"`
	Paging struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Count  int `json:"count"`
		Pages  int `json:"pages"`
	} `json:"paging"`
}

// NetworkGroupsInput .
type NetworkGroupsInput struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Objects []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"objects"`
	Literals []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"literals"`
}

// GetNetworkGroups Get networkgroups from fmc
func (c *Client) GetNetworkGroups() (*NetworkGroups, error) {
	// todo: implement limits
	// todo: implement filtering
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/fmc_config/v1/domain/%s/object/networkgroups?limit=10000&expanded=true", c.baseURL, c.Domain), nil)
	if err != nil {
		return nil, err
	}

	res := NetworkGroups{}
	// create Pointer for Network Struct
	err = c.sendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}
	return &res, nil
}

// CreateNetworkGroups will create a Network from FMC
func (c *Client) CreateNetworkGroups(i string) (*NetworkGroups, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/fmc_config/v1/domain/%s/object/networkgroups", c.baseURL, c.Domain), nil)
	if err != nil {
		return nil, err
	}
	var test NetworkGroupsInput
	//test that the json structur is OK
	err = json.Unmarshal([]byte(i), &test)
	if err != nil {
		return nil, fmt.Errorf("error found in input '%v' failed with '%v'", i, err)
	}

	req.Body = ioutil.NopCloser(strings.NewReader(i))
	res := NetworkGroups{}
	// create Pointer for Network Struct
	err = c.sendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}
	return &res, nil
}

// SearchNetworkGroups .
func (c *Client) SearchNetworkGroups(network *NetworkGroupsInput, networks *NetworkGroups) (int, error) {
	var (
		err error
	)
	if err != nil {
		return 0, err
	}
	idx := sort.Search(len(networks.Items), func(i int) bool {
		return string(networks.Items[i].Name) >= network.Name
	})
	return idx, err
}

// UpdateNetworkGroupsByObject will overwrite a Network Object
func (c *Client) UpdateNetworkGroupsByObject(i string) (*NetworkGroups, error) {
	var (
		network NetworkGroupsInput
		err     error
		res     *NetworkGroups
	)
	json.Unmarshal([]byte(i), &network)
	networks, err := c.GetNetworkGroups()
	if err != nil {
		return nil, err
	}
	idx, err := c.SearchNetworkGroups(&network, networks)
	if err != nil {
		return nil, err
	}
	// Create Network if not exists
	if networks.Items[idx].Name != network.Name {
		// create new object
		c.CreateNetworkGroups(i)
		if err != nil {
			return nil, err
		}
		networks, err := c.GetNetworkGroups()
		if err != nil {
			return nil, err
		}

		// Get Network list
		idx, err := c.SearchNetworkGroups(&network, networks)
		if err != nil {
			return nil, err
		}
		network.ID = networks.Items[idx].ID
	} else {
		network.ID = networks.Items[idx].ID
	}

	res, err = c.UpdateNetworkGroups(network)
	if err != nil {
		return nil, err
	}
	//
	return res, err
	//
}

// UpdateNetworkGroups will overwrite a Network Object
func (c *Client) UpdateNetworkGroups(i NetworkGroupsInput) (*NetworkGroups, error) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/fmc_config/v1/domain/%s/object/networkgroups/%s", c.baseURL, c.Domain, i.ID), nil)
	if err != nil {
		return nil, err
	}
	t, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(t))
	res := NetworkGroups{}

	// create Pointer for Network Struct
	err = c.sendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}
	return &res, nil
}
