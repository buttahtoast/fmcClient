package fmcClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client api needed.
type Client struct {
	user        string
	password    string
	baseURL     string
	HTTPClient  *http.Client
	accessToken string
	Domain      string
}

// NewClient creates new FMC client with given Username and Password it will login
func NewClient(user string, password string, baseurl string, domain string, cert []byte) (*Client, error) {
	t := &Client{
		user:     user,
		password: password,
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
		baseURL: baseurl,
		Domain:  domain,
	}
	// add cert if defined

	if cert != nil {
		err := t.trustcert(cert)
		if err != nil {
			return nil, err
		}
	}
	// Login and retrun new Client
	m, err := t.login()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// adding cert from Client to trusted certificate
func (c *Client) trustcert(cert []byte) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		return fmt.Errorf("failed to parse certificate")
	}
	tlsConf := &tls.Config{RootCAs: roots}
	tr := &http.Transport{TLSClientConfig: tlsConf}
	c.HTTPClient.Transport = tr

	// if no error return nill
	return nil
}

// Login to FMC Center to get accessToken
func (c *Client) login() (*Client, error) {

	var buf bytes.Buffer
	req, err := http.NewRequest("POST", c.baseURL+"/fmc_platform/v1/auth/generatetoken", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.user, c.password)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 401 {
		return nil, fmt.Errorf("wrong username or password %d %v", res.StatusCode, req.URL)
	} else if res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("cannot login unknown error, status code: %d %v", res.StatusCode, req.URL)
	}
	c.accessToken = res.Header.Get("X-Auth-Access-Token")
	return c, nil
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Auth-Access-Token", c.accessToken)
	//c.HTTPClient.Timeout
	//t.HTTPClient.Timeout = time.Second * 6000

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		//
	} else if res.StatusCode == http.StatusCreated {
		//
	} else if res.StatusCode == http.StatusAccepted {
		//
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unknown error, status code: %d  URL: %v %v", res.StatusCode, req.URL, ioutil.NopCloser(bytes.NewBuffer(body)))
	}
	// parse to json with pointer from calling function
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return fmt.Errorf("test %v", err)
	}

	return nil
}
