package main

import (
	"fmt"

	"github.com/harakeishi/gscp"
)

func main() {
	s, _ := gscp.LoadConfig()
	fmt.Println(gscp.Parse(s))
}
