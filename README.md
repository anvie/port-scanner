Simple port scanner library for Go

Example Usage
----------------

```go
package main

import (
    "fmt"
	"github.com/anvie/port-scanner"
)

func main(){

     ps := portscanner.NewPortScanner("localhost")

     // get opened port
     fmt.Printf("scanning port %d-%d...\n", 20, 30000)

     openedPorts := ps.GetOpenedPort(20, 30000)

     for i := 0; i < len(openedPorts); i++ {
     	port := openedPorts[i]
     	fmt.Print(" ", port, " [open]")
     	fmt.Println("  -->  ", ps.DescribePort(port))
     }
}


```

