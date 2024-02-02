package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/melbahja/goph"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"os"
)

type HostsSsh struct {
	Hosts []Config `yaml:"hosts"`
}

type Config struct {
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config"`
}

type SshInfo struct {
	Hostname string `mapstructure:"hostname" json:"hostname"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Port     string `mapstructure:"port" json:"port"`
}

func getYamlInformation() HostsSsh {
	var hosts HostsSsh

	yamlFile, _ := os.ReadFile("./hosts.yaml")
	err := yaml.Unmarshal(yamlFile, &hosts)
	if err != nil {
		fmt.Println(err)
	}
	return hosts
}

func collectInformationSSH() {

	hosts := getYamlInformation()

	for _, host := range hosts.Hosts {

		fmt.Println(host.Name)
		var sshInfo SshInfo
		err := mapstructure.Decode(host.Config, &sshInfo)
		if err != nil {
			fmt.Println(err)
		}
		result := getHostInformationWithSSH(sshInfo)
		fmt.Println(result)
	}

}

func getHostInformationWithSSH(sshInfo SshInfo) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	client, err := goph.New(sshInfo.Username, sshInfo.Hostname, goph.Password(sshInfo.Password))
	if err != nil {
		fmt.Println(err)
	}
	out, err2 := client.Run("ip addr | grep \\\\/24 | awk '{ print $2 }'")

	if err2 != nil {
		fmt.Println(err2)
	}

	result["ip"] = string(out)

	return result

}

func runServers() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Run as Server mode")
	})

	app.Listen(":4000")
}

func main() {
	runServer := flag.Bool("server", false, "Run as Server mode")
	runCollectApi := flag.Bool("api", false, "Collect API")

	flag.Parse()

	if *runServer {
		runServers()
	} else {
		if *runCollectApi {
			fmt.Println("Collect API")
		} else {
			collectInformationSSH()
		}

	}
}
