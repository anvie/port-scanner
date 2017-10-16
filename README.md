Simple port scanner library for Go

Installation
--------------

```bash
$ go get github.com/anvie/port-scanner
```

Example Usage
----------------

```go
package main

import (
    "fmt"
	"time"
	"github.com/anvie/port-scanner"
)

func main(){
     // scan localhost with a 2 second timeout per port in 5 concurrent threads
     ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)

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

output:

```
scanning port 20-30000...
 22 [open]  -->   SSH
 25 [open]  -->   SMTP
 53 [open]  -->   <unknown>
 80 [open]  -->   web service nginx
 139 [open]  -->   netbios
 445 [open]  -->   Samba
 548 [open]  -->   <unknown>
 587 [open]  -->   <unknown>
 631 [open]  -->   cups
 2181 [open]  -->   <unknown>
 5800 [open]  -->   VNC remote desktop
 5900 [open]  -->   <unknown>
 6379 [open]  -->   <unknown>
 6942 [open]  -->   <unknown>
 9009 [open]  -->   <unknown>
 17500 [open]  -->   <unknown>
 27017 [open]  -->   mongodb [ http://www.mongodb.org/ ]
```
