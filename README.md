# The Golang SWC Registry Library

This package aims to make the [Smart Contract Weakness Classification](https://smartcontractsecurity.github.io/SWC-registry/) (SWC) Registry accessible in Golang. To get an overview over the available methods, please check the Go Doc for this repo.

## Example
```go
package main

import (
    swc "github.com/SmartContractSecurity/SWC-registry-go/pkg"
    "fmt"
)


func main() {
    s, err := swc.GetSWC("SWC-101")
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(s.GetMarkdown())
    }
}
```

## Behaviour
On first use of the `GetSWC` function, the SWC registry is initialized from the SWC classification repository if it's not filled already. Further queries are then relayed to the already downloaded dataset - unless an update is triggered explicitly by the user.

There are two ways to update the internal SWC registry: from a local file or a remote URL. The user will have to get the package's SWC registry instance (a singleton and thread-safe) and call one of the update methods:

### Update from URL
```go
r := GetRegistry()
err := r.UpdateRegistryFromURL()
if err != nil {
    t.Fatalf("Loading from URL failed. Got error %s", err)
}
```
Alternatively, a custom URL can be passed. Otherwise we will default to the JSON file in the official SWC Registry repo: Also contained as `DefaultGithubURL` in the package.

### Update from File
```go
r := GetRegistry()
err := r.UpdateRegistryFromFile()
if err != nil {
    t.Fatalf("Loading from file failed. Got error %s", err)
}
```
This works analogously to the above example and defaults to `DefaultFilePath` unless given another valid file path to a JSON file.