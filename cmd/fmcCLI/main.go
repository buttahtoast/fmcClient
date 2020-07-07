package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chifu1234/fmcClient"
)

func main() {
	// parse args
	arg := args()
	var cert []byte

	// read certificate
	cert, err := ioutil.ReadFile(arg.Cert)
	if err != nil {
		if err.Error() == "open : no such file or directory" && arg.Cert == "" {
			// do nothing
		} else {
			log.Fatal(err)
		}
	}

	// create basic client and login
	t, err := fmcClient.NewClient(arg.user, arg.password, arg.baseURL, arg.Domain, cert)
	if err != nil {
		log.Fatal(err)
	}

	// run function
	switch arg.Function {
	case "GetNetworks":
		res, err := t.GetNetworks()
		if err != nil {
			log.Fatal(err)
		}
		//pprint
		empJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(empJSON))
	case "CreateNetworks":
		res, err := t.CreateNetworks(arg.Input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", res)
	// case "UpdateNetworks":
	// 	res, err := t.UpdateNetworks(arg.Input)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%+v\n", res)
	case "UpdateNetworksByObject":
		res, err := t.UpdateNetworksByObject(arg.Input)
		if err != nil {
			log.Fatal(err)
		}
		empJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(empJSON))
	default:
		log.Fatalf("Function not known to me: %v", arg.Function)
	}
}
