package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

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

		input := strings.SplitN(inputStr, " ", 2)
		var result []string
		if len(input) > 1 {
			result = processArgs(input[1])
		}
		// whatt?
		switch input[0] {
		case "echo":
			fmt.Println(strings.Join(result, " "))
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("error")
			}
			fmt.Println(dir)
		case "cd":
			pathChange := ""
			if strings.TrimSpace(input[1]) == "~" {
				pathChange = os.Getenv("HOME")
			} else {
				pathChange = path.Clean(input[1])
			}
			if !path.IsAbs(pathChange) {
				dir, _ := os.Getwd()
				pathChange = path.Join(dir, pathChange)
			}
			err = os.Chdir(pathChange)
			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", pathChange)
			}
		case "cat":
			cat(result)
		case "type":
			switch input[1] {
			case "echo":
				fmt.Println("echo is a shell builtin")
			case "exit":
				fmt.Println("exit is a shell builtin")
			case "type":
				fmt.Println("type is a shell builtin")
			case "pwd":
				fmt.Println("pwd is a shell builtin")
			case "cd":
				fmt.Println("cd is a shell builtin")
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
	paths := strings.Split(env, ":")
	for _, path := range paths {
		exec := filepath.Join(path, input)
		if _, err := os.Stat(exec); err == nil {
			fmt.Fprintf(os.Stdout, "%v is %v\n", input, exec)
			return
		}
	}
	fmt.Printf("%s: not found\n", input)
}

func cat(result []string) {
	for _, filename := range result {
		data, err := os.Open(filename)
		if err != nil {
			fmt.Println("error, invalid argument")
			return
		}
		defer data.Close()
		content, err := io.ReadAll(data)
		if err != nil {
			fmt.Println("rerrr reading file")
			return
		}
		fmt.Print(string(content))
	}
}

// magic
func processArgs(argstr string) []string {
	var singleQuote bool
	var doubleQuote bool
	var backslash bool
	var arg string
	var args []string
	for _, r := range argstr {
		switch r {
		case '\'':
			if backslash && doubleQuote {
				arg += "\\"
			}
			if backslash || doubleQuote {
				arg += string(r)
			} else {
				singleQuote = !singleQuote
			}
			backslash = false
		case '"':
			if backslash || singleQuote {
				arg += string(r)
			} else {
				doubleQuote = !doubleQuote
			}
			backslash = false
		case '\\':
			if backslash || singleQuote {
				arg += string(r)
				backslash = false
			} else {
				backslash = true
			}
		case ' ':
			if backslash && doubleQuote {
				arg += "\\"
			}
			if backslash || singleQuote || doubleQuote {
				arg += string(r)
			} else if arg != "" {
				args = append(args, arg)
				arg = ""
			}
			backslash = false
		default:
			if doubleQuote && backslash {
				arg += "\\"
			}
			arg += string(r)
			backslash = false
		}
	}
	if arg != "" {
		args = append(args, arg)
	}
	return args
}
