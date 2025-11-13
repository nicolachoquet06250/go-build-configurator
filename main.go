package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

var currentDir, _ = os.Getwd()

var configFile = flag.String("file", "build.go.json", "config file to read")
var configSystemBuild = flag.String("system", "", "on what system you will build your program")
var configArchBuild = flag.String("arch", "", "on what arch you will build your program")

var install = flag.Bool("install", false, "Install executable")

type ConfigItem struct {
	Output     string `yaml:"output,omitempty" json:"output,omitempty"`
	Executable string `yaml:"executableName,omitempty" json:"executableName,omitempty"`
}

type Config struct {
	Linux   ConfigItem `yaml:"linux,omitempty" json:"linux,omitempty"`
	Windows ConfigItem `yaml:"windows,omitempty" json:"windows,omitempty"`
	Darwin  ConfigItem `yaml:"darwin,omitempty" json:"darwin,omitempty"`
}

var config Config

var Systems = map[string]func() ConfigItem{
	"windows": func() ConfigItem { return config.Windows },
	"linux":   func() ConfigItem { return config.Linux },
	"darwin":  func() ConfigItem { return config.Darwin },
}

var System func() ConfigItem

func main() {
	flag.Parse()

	if *install {
		exePath, _ := os.Executable()
		homeDir, _ := os.UserHomeDir()
		exeName := strings.Split(exePath, string(os.PathSeparator))[len(strings.Split(exePath, string(os.PathSeparator)))-1]

		destPath := homeDir + string(os.PathSeparator) + "go-build-configurator"

		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			_ = os.MkdirAll(destPath, os.ModePerm)
		}
		destPath += string(os.PathSeparator) + exeName

		Install(exePath, destPath)

		return
	}

	if *configSystemBuild == "" {
		*configSystemBuild = runtime.GOOS
		if os.Getenv("GOOS") != "" {
			*configSystemBuild = os.Getenv("GOOS")
		}
	}
	if *configSystemBuild != "" {
		_ = os.Setenv("GOOS", *configSystemBuild)
	}

	if *configArchBuild == "" {
		*configArchBuild = runtime.GOARCH
		if os.Getenv("GOOS") != "" {
			*configSystemBuild = os.Getenv("GOOS")
		}
	}
	if *configArchBuild != "" {
		_ = os.Setenv("GOARCH", *configArchBuild)
	}

	configFilePath := currentDir + string(os.PathSeparator) + *configFile

	if _, err := os.Stat(configFilePath); err != nil {
		panic(fmt.Errorf("config file not found: %s", *configFile))
	}

	if strings.Contains(*configFile, ".yaml") {
		content, _ := os.ReadFile(configFilePath)

		if err := yaml.Unmarshal(content, &config); err != nil {
			panic(fmt.Errorf("error parsing config file: %s", err))
		}
	} else if strings.Contains(*configFile, ".json") {
		content, _ := os.ReadFile(*configFile)

		if err := json.Unmarshal(content, &config); err != nil {
			panic(fmt.Errorf("error parsing config file: %s", err))
		}
	}

	System = Systems[*configSystemBuild]

	Build()
}
