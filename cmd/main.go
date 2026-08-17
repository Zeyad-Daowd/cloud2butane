package main

import (
	"encoding/json"
	"fmt"
	"os"

	cloud2butane "github.com/Zeyad-Daowd/cloud2butane"
	"gopkg.in/yaml.v3"
)

type Test struct {
	Mode int `yaml:"mode"`
}

func main() {
	// read the cloud-config file from arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <cloud-config-file> --debug(optional)")
		return
	}
	cloudConfigFile := os.Args[1]

	debug := false
	if len(os.Args) > 2 && os.Args[2] == "--debug" {
		debug = true
	}
	// read the cloud-config file
	data, err := os.ReadFile(cloudConfigFile)
	if err != nil {
		fmt.Printf("Error reading cloud-config file: %v\n", err)
		return
	}

	// parse the cloud-config file
	var cloudConfig cloud2butane.CloudConfig
	err = yaml.Unmarshal(data, &cloudConfig)
	if err != nil {
		fmt.Printf("Error parsing cloud-config file: %v\n", err)
		return
	}

	if debug {
		// print the parsed cloud-config as a json object
		output, err := json.MarshalIndent(cloudConfig, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("========DEBUG MODE========")
		fmt.Println("Parsed cloud-config:")
		fmt.Println(string(output))
		fmt.Println("==========================")
	}

	butane, err := cloud2butane.TranslateCloudConfig(cloudConfig)
	if err != nil {
		fmt.Printf("Error translating cloud-config to butane: %v\n", err)
		return
	}

	// print the butane struct as yaml
	butaneOutput, err := yaml.Marshal(butane)
	if err != nil {
		fmt.Printf("Error marshaling Butane config: %v\n", err)
		return
	}

	fmt.Println(string(butaneOutput))

}
