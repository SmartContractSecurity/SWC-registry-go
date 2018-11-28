package main

import (
    swc "github.com/SmartContractSecurity/SWC-registry-go/pkg"
    "fmt"
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