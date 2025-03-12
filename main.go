package main

import (
	"fmt"
	"github.com/philippeckel/pair/cmd"
	"github.com/philippeckel/pair/internal/config"
)

func main() {

	if err := config.InitViper(); err != nil {
		fmt.Printf("Warning: %v\n", err)
		// Continue execution even if config initialization fails
	}
	commands.Execute()
}
