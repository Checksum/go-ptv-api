Go API for Public Transport Victoria
====================================

Installation
------------

`go get github.com/Checksum/go-ptv-api`

Usage
-----

```go
import "github.com/Checksum/go-ptv-api"

func main() {
    client := ptv.NewClient("developerID", "secretKey")
    healthCheck, resp, err := client.HealthCheck.Get()
}
```
