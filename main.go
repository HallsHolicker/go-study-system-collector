package main

import (
	"Go-Collector/utils"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/k0kubun/pp/v3"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

type Hosts struct {
	Hosts *collectType `yaml:"hosts"`
}

type collectType struct {
	Ssh []*Config `yaml:"ssh"`
	Api []string  `yaml:"api"`
}

type Config struct {
	Name   string   `yaml:"name"`
	Config *SshInfo `yaml:"config"`
}

type SshInfo struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     uint   `yaml:"port"`
	Sudo     bool   `yaml:"sudo"`
}

var wg sync.WaitGroup

func getYamlInformation() Hosts {
	var hosts Hosts

	yamlFile, _ := os.ReadFile("./hosts.yaml")
	err := yaml.Unmarshal(yamlFile, &hosts)
	if err != nil {
		fmt.Println(err)
	}
	return hosts
}

func collectInformation() {
	hosts := getYamlInformation().Hosts

	var sshRequestChannel chan SshInfo = make(chan SshInfo)
	//var ResponseChannel chan SshInfo = make(chan SshInfo)
	var apiRequestChannel chan string = make(chan string)

	wg.Add(len(hosts.Ssh) + len(hosts.Api))
	go collectSSH(sshRequestChannel)
	go collectAPI(apiRequestChannel)

	for _, host := range hosts.Ssh {
		fmt.Println(host.Name)
		sshRequestChannel <- *host.Config
		//result := getHostInformationWithSSH(*host.Config)
		//fmt.Println(result)
	}
	for _, host := range hosts.Api {
		apiRequestChannel <- host
	}

	wg.Wait()
	close(sshRequestChannel)
	close(apiRequestChannel)
}

func collectSSH(c chan SshInfo) {
	for {
		select {
		case sshInfo := <-c:
			var hw utils.Hardware
			//fmt.Println(sshInfo)
			//result := make(map[string]map[string]string)
			sshConfig := goph.Config{
				User:     sshInfo.Username,
				Addr:     sshInfo.Hostname,
				Port:     sshInfo.Port,
				Auth:     goph.Password(sshInfo.Password),
				Timeout:  goph.DefaultTimeout,
				Callback: ssh.InsecureIgnoreHostKey(),
			}
			//client, err := goph.New(sshInfo.Username, sshInfo.Hostname, goph.Password(sshInfo.Password))
			client, err := goph.NewConn(&sshConfig)
			if err != nil {
				fmt.Println(err)
			}
			//out, err2 := client.Run("ip addr | grep \\\\/24 | awk '{ print $2 }'")
			//cpuName, cpuNameErr := client.Run("grep 'model name' /proc/cpuinfo | cut -f2 -d ':' | uniq")
			//utils.HandleErr(cpuNameErr)
			//cpuCore, cpuCoreErr := client.Run("grep 'cpu cores' /proc/cpuinfo | tail -1")
			//utils.HandleErr(cpuCoreErr)
			//cpuCount, cpuCountErr := client.Run("grep 'physical id' /proc/cpuinfo | sort -u")
			//utils.HandleErr(cpuCountErr)

			//var testmap map[string]interface{}

			var command string

			if sshInfo.Sudo {
				command = "sudo "
			} else {
				command = ""
			}

			lshwCommand := command + "lshw -quiet -xml"
			test, err := client.Run(lshwCommand)

			//test, err := client.Run("cat /root/test.log")
			//test, err := os.ReadFile("./test2.log")

			utils.HandleErr(err)
			//fmt.Println(string(test))
			hw = utils.LshwParser(test)

			//Collector Nvme Disk
			nvmeCommand := "sudo nvme list"
			nvmeInfo, err := client.Run(nvmeCommand)
			lines := strings.Split(string(nvmeInfo), "\n")

			if len(lines) > 2 {
				nvme := utils.Nvme{}
				parts := strings.Fields(lines[2])

				nvme.Vendor = parts[2]
				nvme.Serial = parts[1]
				nvme.Logicalname = parts[0]

				if nvme.Vendor == "Samsung" {
					nvme.Product = parts[3] + " " + parts[4] + " " + parts[5]
					nvme.Size = parts[8] + parts[9]
				} else if nvme.Vendor == "INTEL" {
					nvme.Product = parts[3]
					nvme.Size = parts[5] + parts[6]
				}

				hw.Nvmes = append(hw.Nvmes, nvme)
			}
			utils.HandleErr(err)
			//pp.Print(string(nvmeInfo))

			//if err2 != nil {
			//	fmt.Println(err2)
			//}
			//
			//result["network"] = make(map[string]string)
			//result["network"]["ip"] = string(out)
			//result["cpu"] = make(map[string]string)
			//result["cpu"]["name"] = string(cpuName)
			//result["cpu"]["Core"] = string(cpuCore)
			//result["cpu"]["Count"] = string(cpuCount)
			//
			//fmt.Println(result)

			//pp.Print(hw.Nvmes)

			pp.Print(hw)

			wg.Done()
		}
	}
}

func collectAPI(c chan string) {
	for {
		select {
		case apiHosts := <-c:
			fmt.Println(apiHosts)
			wg.Done()
		}
	}
}

func runServers() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Run as Server mode")
	})

	utils.HandleErr(app.Listen(":4000"))
}

func main() {
	runServer := flag.Bool("server", false, "Run as Server mode")

	flag.Parse()

	if *runServer {
		runServers()
	} else {
		collectInformation()
	}
}
