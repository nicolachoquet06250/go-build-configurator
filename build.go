package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Build() {
	pwd, _ := os.Getwd()

	sys := System()

	completeDestPath := pwd
	if sys.Output != "" {
		completeDestPath = pwd + string(os.PathSeparator) + sys.Output
		if _, err := os.Stat(completeDestPath); os.IsNotExist(err) {
			if err = os.MkdirAll(completeDestPath, 0755); err != nil {
				panic(fmt.Errorf("can't create output directory:" + completeDestPath))
			}
		} else if err != nil {
			panic(err)
		}
	}

	if sys.Executable != "" {
		completeDestPath += string(os.PathSeparator) + sys.Executable
	} else {
		goModFile, _ := os.ReadFile(pwd + string(os.PathSeparator) + "go.mod")
		goModFileLines := strings.Split(string(goModFile), "\n")
		moduleParts := strings.Split(strings.Replace(goModFileLines[0], "module ", "", 1), "/")

		completeDestPath += string(os.PathSeparator) + moduleParts[len(moduleParts)-1] + ".exe"
	}

	out, err := exec.Command("go", "build", "-o", completeDestPath).CombinedOutput()
	println(string(out))
	if err != nil {
		panic(err)
	}
}
