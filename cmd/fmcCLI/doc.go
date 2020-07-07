//
/*
This is a command tool to interact with the cisco fmc.

It uses the fmc api to get data from the fmc.

You can find all availble options by running the application with option "-h".

Usage of fmcCLI:
  -cert string
    	adding x509 Certificate if client does not trust the fmc certificate
  -domain string
    	domain id (default "global")
  -function string
    	possible GetNetworks (default "GetNetworks")
  -pw string
    	Username Password (default "admin")
  -u string
    	url to FPM api (default "https://fmc/api")
  -url string
    	url to FPM api (default "https://fmc/api")
  -user string
    	API Username (default "admin")
*/
package main // import "github.com/chifu1234/fmcCLI"
