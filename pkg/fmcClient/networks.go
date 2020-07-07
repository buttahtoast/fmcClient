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

// Network structure for FMC .
type Network struct {
	Items []struct {
		Links struct {
			Self   string `json:"self"`
			Parent string `json:"parent"`
		} `json:"links"`
		Type        string `json:"type"`
		Value       string `json:"value"`
		Overridable bool   `json:"overridable"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Name        string `json:"name"`
	} `json:"items"`
}

// CreateNetworkOutput fmc struct for every network
type CreateNetworkOutput []struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Overrides struct {
		Target struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"target"`
	} `json:"overrides"`
	Overridable bool   `json:"overridable"`
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

// CreateNetworkInput structur for POST .
type CreateNetworkInput struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Overridable bool   `json:"overridable"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ID          string `json:"id"`
}

type UpdateNetworkOutput struct {
	Links struct {
		Self   string `json:"self"`
		Parent string `json:"parent"`
	} `json:"links"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Overridable bool   `json:"overridable"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Metadata    struct {
		LastUser struct {
			Name string `json:"name"`
		} `json:"lastUser"`
		Domain struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"domain"`
		IPType     string `json:"ipType"`
		ParentType string `json:"parentType"`
	} `json:"metadata"`
}

// GetNetworks will retrun all Networks from FMC
func (c *Client) GetNetworks() (*Network, error) {
	// todo: implement limits
	// todo: implement filtering
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/fmc_config/v1/domain/%s/object/networks?expanded=true", c.baseURL, c.Domain), nil)
	if err != nil {
		return nil, err
	}

	res := Network{}
	// create Pointer for Network Struct
	err = c.sendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}
	return &res, nil
}

// CreateNetworks will create a Network from FMC
func (c *Client) CreateNetworks(i string) (*CreateNetworkOutput, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/fmc_config/v1/domain/%s/object/networks", c.baseURL, c.Domain), nil)
	if err != nil {
		return nil, err
	}
	var test CreateNetworkInput
	//test that the json structur is OK
	err = json.Unmarshal([]byte(i), &test)
	if err != nil {
		return nil, fmt.Errorf("error found in input '%v' failed with '%v'", i, err)
	}

	req.Body = ioutil.NopCloser(strings.NewReader(i))
	res := CreateNetworkOutput{}
	// create Pointer for Network Struct
	err = c.sendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}
	return &res, nil
}

// UpdateNetworksByObject will overwrite a Network Object
func (c *Client) UpdateNetworksByObject(i string) (*UpdateNetworkOutput, error) {
	var (
		network CreateNetworkInput
		err     error
		res     *UpdateNetworkOutput
	)
	json.Unmarshal([]byte(i), &network)
	networks, err := c.GetNetworks()
	if err != nil {
		return nil, fmt.Errorf("test %v", err)
	}
	idx := sort.Search(len(networks.Items), func(i int) bool {
		return string(networks.Items[i].Name) >= network.Name
	})
	network.ID = networks.Items[idx].ID
	res, err = c.UpdateNetworks(network)
	if err != nil {
		return nil, err
	}
	//
	return res, err
	//
}

// UpdateNetworks will overwrite a Network Object
func (c *Client) UpdateNetworks(i CreateNetworkInput) (*UpdateNetworkOutput, error) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/fmc_config/v1/domain/%s/object/networks/%s", c.baseURL, c.Domain, i.ID), nil)
	if err != nil {
		return nil, err
	}
	t, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(t))
	res := UpdateNetworkOutput{}
	// create Pointer for Network Struct
	err = c.sendRequest(req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}
	return &res, nil
}
