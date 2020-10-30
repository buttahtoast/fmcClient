package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kubernetli/fmcClient/pkg/fmcClient"
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

	// read input file
	input, err := ioutil.ReadFile(arg.Input)
	if err != nil {
		if err.Error() == "open : no such file or directory" && arg.Input == "" {
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

	//increase http timeout
	t.HTTPClient.Timeout = time.Minute * arg.Timeout

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
		res, err := t.CreateNetworks(string(input))
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
	case "UpdateNetworkGroupsByObject":
		res, err := t.UpdateNetworkGroupsByObject(string(input))
		if err != nil {
			log.Fatal(err)
		}
		empJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(empJSON))
	case "GetNetworkGroups":
		res, err := t.GetNetworkGroups()
		if err != nil {
			log.Fatal(err)
		}
		//pprint
		empJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(empJSON))
	default:
		log.Fatalf("Function not known to me: %v", arg.Function)
	}
}
