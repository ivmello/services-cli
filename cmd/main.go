package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Services struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name string `yaml:"name"`
	Path  string `yaml:"path"`
	Command  []string `yaml:"command"`
}

func (s *Services) getConf(){
	yamlFile, err := os.ReadFile("./services.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func execCommand(dir string, command string, args ...string) (io.ReadCloser, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return stdout, nil
}

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI tool for starting services",
	Run: func(cmd *cobra.Command, args []string) {
		flagsString := cmd.Flags().Lookup("services").Value.String()
		flags := strings.Split(flagsString, ",")
		var services Services
		services.getConf()
		if len(flags) > 0 {
			var wg sync.WaitGroup
			wg.Add(len(flags))
			for _, flag := range flags {
				for _, service := range services.Services {
					if service.Name == flag {
						go func(service Service) {
							defer wg.Done()
							stdout, _ := execCommand(service.Path, service.Command[0], service.Command[1:]...)
							fmt.Printf("Starting %s\n", service.Name)
							if _, err := io.Copy(os.Stdout, stdout); err != nil {
								fmt.Println(err)
								return
							}
						}(service)
					}
				}
			}
			wg.Wait()
		}
	},
}

func init() {
	rootCmd.Flags().StringP("services", "s", "", "Select services to start")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
