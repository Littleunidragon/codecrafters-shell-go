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
	"regexp"
	"strings"
)

func main() {
	for {
		//what about now?
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
		input := strings.SplitN(inputStr, " ", 2)
		reg := regexp.MustCompile(`"((?:\\.|[^"\\])*)"|'([^']*)'|(\S+)`) // WARNING! double quotes doesnt work as wxpwcted yet!
		var args [][]string
		var result []string
		if len(input) > 1 {
			args = reg.FindAllStringSubmatch(input[1], -1)
			for _, arg := range args {
				if arg[1] != "" { // group 1 = ""
					result = append(result, treatSpecialChar(arg[1]))
				} else if arg[2] != "" { // group 2 =''
					result = append(result, arg[2])
				} else if arg[3] != "" { //group 3 = standalone words
					result = append(result, treatSpecialChar(arg[3]))
				}
			}
		}

		switch input[0] {
		case "echo":
			for i := 0; i < len(result); i++ {
				fmt.Print(result[i] + "")
			}
			fmt.Println()
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
	//fmt.Println(env)
	paths := strings.Split(env, ":") // for windows its ";", but tests done with ":" anywayyy os.PathListSeparator
	for _, path := range paths {
		//fmt.Println(path)
		exec := filepath.Join(path, input) // literary checks every possible path, to find THE one
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

func treatSpecialChar(escapeSpace string) string {
	double := regexp.MustCompile(`\\{2}`)
	placeholder := double.ReplaceAllString(escapeSpace, "\u0000")
	placeholder = strings.ReplaceAll(placeholder, "\\", "")
	return strings.ReplaceAll(placeholder, "\u0000", `\`)
}
