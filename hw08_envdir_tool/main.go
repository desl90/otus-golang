package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatal("Usage: go-envdir /path/to/evndir command arg1 arg2")
	}

	dir, cmd := args[1], args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("Environment directory not found: %s", dir)
	}

	code := RunCmd(cmd, env)

	os.Exit(code)
}
