package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abdivasiyev/gogrep/internal/app"
)

func main() {
	// 0 - app name
	// 1 - command
	// 2 - pattern
	if len(os.Args) < 3 {
		fmt.Println("Usage: gogrep -p \"gogrep\" <file>")
		os.Exit(1)
	}

	if os.Args[1] != "-p" {
		fmt.Println("Usage: gogrep -p \"gogrep\" <file>")
		os.Exit(1)
	}

	var (
		pattern = os.Args[2]
		srcList = os.Args[3:]
	)

	ctx := context.Background()

	app, err := app.New(pattern, srcList)
	if err != nil {
		fmt.Printf("gopgrep: %s\n", err)
		os.Exit(1)
	}

	err = app.Run(ctx)
	if err != nil {
		fmt.Printf("gopgrep: %s\n", err)
		os.Exit(1)
	}
}
