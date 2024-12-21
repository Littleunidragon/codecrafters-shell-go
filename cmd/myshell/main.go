package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		inputStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
		inputStr = strings.TrimRight(inputStr, "\n")
		inputStr = strings.TrimSpace(inputStr)
		if err != nil {
			log.Fatal(err)
		}

		switch inputStr {
		case "exit 0":
			os.Exit(0)
		}

		input := strings.Split(inputStr, " ")
		switch input[0] {
		case "echo":
			fmt.Println(input[1])
		case "type":
			switch input[1] {
			case "echo":
				fmt.Println("echo is a shell builtin")
			case "exit":
				fmt.Println("exit is a shell builtin")
			case "type":
				fmt.Println("type is a shell builtin")
			default:
				builtin(input[1])
			}
		default:
			command := exec.Command(input[0], input[1:]...)
			command.Stdout = os.Stdout
			err = command.Run()
			if err != nil {
				fmt.Printf("%s: command not found\n", input[0])
			}
		}
	}

}

func builtin(input string) {
	env := os.Getenv("PATH")
	//fmt.Println(env)
	paths := strings.Split(env, ":") // for windows its ";", but tests done with ":" anywayyy os.PathListSeparator
	for _, path := range paths {
		//fmt.Println(path)
		exec := filepath.Join(path, input) // literary checks every possible path, to find the needed one?
		if _, err := os.Stat(exec); err == nil {
			fmt.Fprintf(os.Stdout, "%v is %v\n", input, exec)
			return
		}
	}
	fmt.Printf("%s: not found\n", input)
}
