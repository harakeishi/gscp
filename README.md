# gscp
`gscp` is the ssh-config parser available in go.

`gscp` is named as an abbreviation of `go ssh config parser`.


## Importing
```go
import (
    "github.com/harakeishi/gscp"
)
```

## Documentation
Visit the docs on [GoDoc](https://pkg.go.dev/github.com/harakeishi/gscp)

## usage

If no arguments are passed to `LoadConfig()`, `~/.ssh/config` will be read.

```go
package main

import (
	"fmt"
	"github.com/harakeishi/gscp"
)

func main() {
	// Reads a config file and gets it as a string
	s, _ := gscp.LoadConfig()
	// parse
	r, _ := gscp.Parse(s)
	fmt.Printf("%+v", r)
}
```

```sh
$ go run ./cmd/main.go
[{Name:testhost Options:[{Name:HostName Value:192.0.2.1} {Name:User Value:myuser} {Name:IdentityFile Value:~/.ssh/id_rsa} {Name:ServerAliveInterval Value:60}]}]
```

If you want to parse a config in a specific directory, pass the path as follows.

```go
package main

import (
	"fmt"
	"github.com/harakeishi/gscp"
)

func main() {
	// Reads a config file and gets it as a string
	path := gscp.Path("./testData/test1_config")
	s, _ := gscp.LoadConfig(path)
	// parse
	fmt.Println(gscp.Parse(s))
}
```

Obtains information on the specified host and the specified options.

```go
package main

import (
	"fmt"
	"github.com/harakeishi/gscp"
)

func main() {
	// Specify path
	path := gscp.Path("./testData/test1_config")
	// Load config.
	s, _ := gscp.LoadConfig(path)
	// Parse config.
	hosts, _ := gscp.Parse(s)
	// Obtain information on the specified host
	host := hosts.FindHost("testhost")
	fmt.Printf("%+v\n", host)
	// Obtain information on specified options
	hostname := host.FindOption("HostName").Value
	fmt.Printf("%+v\n", hostname)
}
```

```sh
PS C:\src\smp> go run .\cmd\main.go
{Name:testhost Options:[{Name:HostName Value:192.0.2.1} {Name:User Value:myuser} {Name:IdentityFile Value:~/.ssh/id_rsa} {Name:ServerAliveInterval Value:60}]}
192.0.2.1
```

## License
Copyright (c) 2023 harakeishi

Licensed under [MIT](LICENSE)
