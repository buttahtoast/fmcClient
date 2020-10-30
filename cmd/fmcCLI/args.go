package main

import (
	"flag"
)

// Args .
type Args struct {
	user        string
	password    string
	baseURL     string
	accessToken string
	Domain      string
	Function    string
	Cert        string
	Input       string
	Timeout     string
}

// reads/parses user input .
func args() *Args {
	urlPtr := flag.String("url", "https://fmc/api", "url to FPM api")
	flag.StringVar(urlPtr, "u", "https://fmc/api", "url to FPM api")
	domainPtr := flag.String("domain", "global", "domain id")
	functionPtr := flag.String("function", "GetNetworks", "possible GetNetworks|GetNetworkGroups|CreateNetworks|CreateNetworkGroups|UpdateNetworks|UpdateNetworkGroups")
	inputPtr := flag.String("input", "", "function Input in json")
	userPtr := flag.String("user", "admin", "API Username")
	pwPtr := flag.String("pw", "admin", "Username Password")
	fmcCertPtr := flag.String("cert", "", "adding x509 Certificate if client does not trust the fmc certificate")
	timeoutPtr := flag.Int("timeout", "60", "timeout")

	flag.Parse()
	return &Args{
		user:     *userPtr,
		password: *pwPtr,
		baseURL:  *urlPtr,
		Domain:   *domainPtr,
		Function: *functionPtr,
		Cert:     *fmcCertPtr,
		Input:    *inputPtr,
		Timeout:  *timeoutPtr,
	}
}
