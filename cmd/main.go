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
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Command  []string `yaml:"command"`
	Stdout   io.Writer
}

func (s *Services) getConf() {
	yamlFile, err := os.ReadFile("./services.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func execCommand(dir string, command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
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
			for range flags {
				for range services.Services {
					wg.Add(1)
				}
			}
			for _, flag := range flags {
				for _, service := range services.Services {
					if service.Name == flag {
						go func(service Service) {
							defer wg.Done()
							fmt.Printf("Starting %s\n", service.Name)
							cmd, _ := execCommand(service.Path, service.Command[0], service.Command[1:]...)
							cmd.Wait()
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
