package main

import (
	"fmt"
	swc "github.com/SmartContractSecurity/SWC-registry-go/pkg"
)

func main() {
	s, err := swc.GetSWC("SWC-101")
	if err != nil {
		fmt.Println(err)
		fmt.Println(swc.GetRegistry())
	} else {
		fmt.Println(s.GetTitle())
	}
}
