package createHostnameResolution

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type HostnameResolution struct {
	Enabled   bool   `json:"enabled"`
	IPAddress string `json:"ip_address"`
	Hostname  string `json:"hostname"`
}

type Scope struct {
	ProjectOptions struct {
		Connections struct {
			HostnameResolution []HostnameResolution `json:"hostname_resolution"`
		} `json:"connections"`
	} `json:"project_options"`
}

type HostResolutionSlice struct {
	HostnameResolution []HostnameResolution
}

var (
	hostsPath string
	hostnames HostResolutionSlice
)

func parse(hostsContent []byte) []byte {
	scope := Scope{}
	for _, line := range strings.Split(strings.Trim(string(hostsContent), " \t\r\n"), "\n") {
		line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 1 {
			hostnameResolution := HostnameResolution{true, parts[0], parts[1]}
			hostnames.HostnameResolution = append(hostnames.HostnameResolution, hostnameResolution)
		}
	}
	scope.ProjectOptions.Connections.HostnameResolution = hostnames.HostnameResolution
	json, err := json.MarshalIndent(scope, "", "  ")
	if err != nil {
		fmt.Println("json err:", err)
	}
	return json
}

func main() {
	flag.Parse()

	hostsContent, err := ioutil.ReadFile(hostsPath)
	if err == nil {
		jsonFile := parse(hostsContent)
		_ = ioutil.WriteFile("hostname_resolution.json", jsonFile, 0644)
	}
	fmt.Println(fmt.Sprintf("File %s not found", hostsPath))
	os.Exit(1)
}

func init() {
	flag.StringVar(&hostsPath, "hosts", "hosts", "Host file")
}
