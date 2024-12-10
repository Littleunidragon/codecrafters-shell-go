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
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimRight(input, "\n")
	// fmt.Fprintf(os.Stdout, "%s: command not found\n", input)
	fmt.Fprintf(os.Stdout, "%s: command not found\n", input)

}
