package main

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the YAML configuration structure
type Config struct {
	ID            string         `yaml:"id"`
	Name          string         `yaml:"name"`
	Unit          string         `yaml:"unit"`
	Created       string         `yaml:"created"`
	EndpointCalls []EndpointCall `yaml:"endpointCalls"`
	NetworkCalls  []NetworkCall  `yaml:"networkCalls"`
}

// EndpointCall represents a single call to an Endpoint function
type EndpointCall struct {
	Function      string   `yaml:"function"`
	Arguments     []string `yaml:"arguments,omitempty"`
	TestFunction  string   `yaml:"testFunction,omitempty"`
	CleanFunction string   `yaml:"cleanFunction,omitempty"`
}

// NetworkCall represents a single call to a Network function
type NetworkCall struct {
	Function    string              `yaml:"function"`
	URL         string              `yaml:"url,omitempty"`
	Headers     map[string][]string `yaml:"headers,omitempty"`
	QueryParams map[string][]string `yaml:"queryParams,omitempty"`
	Body        string              `yaml:"body,omitempty"`
	Encoding    string              `yaml:"encoding,omitempty"`
	Host        string              `yaml:"host,omitempty"`
	Port        string              `yaml:"port,omitempty"`
	Message     string              `yaml:"message,omitempty"`
	Protocol    string              `yaml:"protocol,omitempty"`
	Hostname    string              `yaml:"hostname,omitempty"`
}

// Template text for generating the Go file
const templateText = `/*
ID: {{.ID}}
NAME: {{.Name}}
UNIT: {{.Unit}}
CREATED: {{.Created}}
*/
package main

import (
	Endpoint "github.com/preludeorg/libraries/go/tests/endpoint"
	Network "github.com/preludeorg/libraries/go/tests/network"
	"fmt"
	"time"
)

func test() {
    // Endpoint calls
    {{- range $.EndpointCalls}}
    {{- if and (ne .Function "Start") (ne .Function "Stop")}}
    Endpoint.{{ .Function }}({{range $i, $arg := .Arguments}}{{if $i}}, {{end}}{{printf "%q" $arg}}{{end}})
    {{- end}}
    {{- end}}

    // Network calls
	{{- range $.NetworkCalls}}
    {{- if eq .Function "GET"}}
        requestOptions := Network.RequestParameters{
            Headers: map[string][]string{"Content-Type": {"application/json"}},
            QueryParams: map[string][]string{
                {{- range $key, $values := .QueryParams}}
                "{{ $key }}": { {{- range $i, $v := $values}}{{if $i}}, {{end}}{{printf "%q" $v}}{{- end}} },
                {{- end}}
            },
            Body: []byte("{{.Body}}"),
        }
        requester := Network.NewHTTPRequest("{{.URL}}", nil)
        response, err := requester.GET(requestOptions)
        if err != nil {
            Endpoint.Say("GET Error: " + err.Error())
        } else {
            Endpoint.Say("GET Response: " + string(response.Body))
        }
    {{- else if eq .Function "POST"}}
        // Similar implementation for POST
    {{- else if eq .Function "TCP"}}
        Network.TCP("{{.Host}}", "{{.Port}}", []byte("{{.Message}}"))
    {{- else if eq .Function "UDP"}}
		Network.UDP("{{.Host}}", "{{.Port}}", []byte("{{.Message}}"))
    {{- else if eq .Function "ScanPort"}}
        isOpen := Network.ScanPort("{{.Protocol}}", "{{.Hostname}}", {{.Port}})
        Endpoint.Say(fmt.Sprintf("ScanPort: Port %d open: %v", {{.Port}}, isOpen))
    {{- end}}
    {{- end}}
}

func clean() {
    Endpoint.Say("Cleaning up")
}

func main() {
    {{- range .EndpointCalls}}
    {{- if eq .Function "Start"}}
    Endpoint.Start(test, clean)
    {{- end}}
    {{- end}}
}
`

func main() {
	// Read YAML configuration
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	// Create generated_code directory if it doesn't exist
	err = os.MkdirAll("generated_code", 0755)
	if err != nil {
		panic(err)
	}

	// Create a file with a UTC timestamp in the name
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("generated_code/generate-%s.go", timestamp)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create and execute the template
	t := template.Must(template.New("goFile").Parse(templateText))
	err = t.Execute(file, conf)
	if err != nil {
		panic(err)
	}

	fmt.Println("File generated:", filename)
}
