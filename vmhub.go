package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	subcommands := `Available subcommands:
get - SNMP get
set - SNMP set
walk - SNMP walk`

	if len(os.Args) < 2 {
		fmt.Println(subcommands)
		os.Exit(1)
	}

	client := Client{
		client:  http.DefaultClient,
		baseURL: "http://192.168.0.1",
	}

	err := client.Login(os.Getenv("VM_PASSWORD"))
	if err != nil {
		fmt.Println("Login failed: ", err)
	}
	var ret []byte

	switch os.Args[1] {
	case "get":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s %s OID ...\n", os.Args[0], os.Args[1])
			os.Exit(1)
		}
		ret, err = client.SNMPGet(os.Args[2:])
	case "set":
		if len(os.Args) != 5 {
			fmt.Printf("Usage: %s %s OID VALUE TYPE\n", os.Args[0], os.Args[1])
			os.Exit(1)
		}
		ret, err = client.SNMPSet(os.Args[2], os.Args[3], os.Args[4])
	case "walk":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s %s OID ...\n", os.Args[0], os.Args[1])
			os.Exit(1)
		}
		ret, err = client.SNMPWalk(os.Args[2:])
	default:
		fmt.Println(subcommands)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Request failed: ", err)
	}
	fmt.Println(string(ret))
}
