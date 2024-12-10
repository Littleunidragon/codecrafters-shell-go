package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		input = strings.TrimRight(input, "\n")
		input = strings.TrimSpace(input)
		if err != nil {
			log.Fatal(err)
		}

		switch input {
		case "exit 0":
			os.Exit(0)
		default:
			fmt.Printf("%s: command not found\n", input)
		}
	}

}
