package main

import (
	"flag"
	"log"
	"os"
	"text/template"
	"time"

	"fmt"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// ImportConfig represents an import configuration.
type ImportConfig struct {
	Alias         string            `yaml:"alias,omitempty"`
	Path          string            `yaml:"path"`
	Filename      string            `yaml:"filename,omitempty"`
	EmbeddedFiles map[string][]byte `yaml:"embeddedFiles,omitempty"`
	EmbeddedFS    map[string]string `yaml:"embeddedFS,omitempty"`
}

// EndpointCall represents a single call to an Endpoint function.
type EndpointCall struct {
	Function      string     `yaml:"function"`
	Arguments     []string   `yaml:"arguments,omitempty"`
	ShellCommands [][]string `yaml:"shellCommands,omitempty"`
	TestFunction  string     `yaml:"testFunction,omitempty"`
	CleanFunction string     `yaml:"cleanFunction,omitempty"`
}

// NetworkCall represents a single call to a Network function.
type NetworkCall struct {
	Function    string              `yaml:"function"`
	URL         string              `yaml:"url,omitempty"`
	Headers     map[string][]string `yaml:"headers,omitempty"`
	QueryParams map[string][]string `yaml:"queryParams,omitempty"`
	Body        string              `yaml:"body,omitempty"`
	Encoding    string              `yaml:"encoding,omitempty"`
	AuthType    string              `yaml:"authType,omitempty"`
	Credential  string              `yaml:"credential,omitempty"`
	Timeout     int                 `yaml:"timeout,omitempty"`
	UserAgent   string              `yaml:"userAgent,omitempty"`
	Host        string              `yaml:"host,omitempty"`
	Port        string              `yaml:"port,omitempty"`
	Message     string              `yaml:"message,omitempty"`
	Protocol    string              `yaml:"protocol,omitempty"`
	Hostname    string              `yaml:"hostname,omitempty"`
	Ports       []int               `yaml:"ports,omitempty"`
}

// EmbeddedFile represents an embedded file configuration.
type EmbeddedFile struct {
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

// Config represents the main configuration structure.
type Config struct {
	ID            string         `yaml:"id"`
	Name          string         `yaml:"name"`
	Unit          string         `yaml:"unit"`
	Created       string         `yaml:"created"`
	Imports       []ImportConfig `yaml:"imports,omitempty"`
	EmbeddedFiles []EmbeddedFile `yaml:"embeddedFiles,omitempty"`
	EndpointCalls []EndpointCall `yaml:"endpointCalls"`
	NetworkCalls  []NetworkCall  `yaml:"networkCalls"`
}

const templateText = `/*
ID: {{.ID}}
NAME: {{.Name}}
UNIT: {{.Unit}}
CREATED: {{.Created}}
*/
package main

// Import necessary packages based on configuration.
{{- if .Imports}}
import (
{{- range .Imports}}
    {{- if .Alias}}
    {{.Alias}} "{{.Path}}"
    {{- else}}
    "{{.Path}}"
    {{- end}}
{{- end}}
)
{{- end}}

// Embedded files content.
{{- range .EmbeddedFiles}}
{{ .Content }}
{{- end}}

func test() {
    {{- range $.EndpointCalls}}
        {{- if eq .Function "Say"}}
        // Say function logic for Endpoint.
        Endpoint.Say({{index .Arguments 0 | printf "%q"}})
        {{- end}}
		
        {{- if eq .Function "Shell"}}
        // Shell function logic for Endpoint.
            {{- range $cmdSet := .ShellCommands}}
            command := []string{ {{- range $i, $cmd := $cmdSet}}{{if $i}}, {{end}}{{printf "%q" $cmd}}{{- end}} }
            Endpoint.Shell(command)
            {{- end}}
        {{- end}}
    
        {{- if eq .Function "Find"}}
        // Find function logic for Endpoint.
        Endpoint.Say("Starting scan for files")
        fileType := {{index .Arguments 0 | printf "%q"}}
        files := Endpoint.Find(fileType)
        if len(files) == 0 {
            Endpoint.Stop(104)
        }
        {{- end}}
    
        {{- if eq .Function "Read"}}
        // Read function logic for Endpoint.
        Endpoint.Say("Reading file")
        file := {{index .Arguments 0 | printf "%q"}}
        contents := Endpoint.Read(file)
        {{- end}}
    
        {{- if eq .Function "Write"}}
        // Write function logic for Endpoint.
        Endpoint.Say("Writing file")
        Endpoint.Write({{index .Arguments 0 | printf "%q"}}, []byte({{index .Arguments 1 | printf "%q"}}))
        {{- end}}
    
        {{- if eq .Function "Exists"}}
        // Exists function logic for Endpoint.
        Endpoint.Say("Checking if file exists")
        exists := Endpoint.Exists({{index .Arguments 0 | printf "%q"}})
        {{- end}}
    
        {{- if eq .Function "Quarantined"}}
        // Quarantined function logic for Endpoint.
        Endpoint.Say("Extracting file for quarantine test")
        filename := {{index .Arguments 0 | printf "%q"}}
        if Endpoint.Quarantined(filename, malicious) {
            Endpoint.Say("Malicious file was caught!")
            Endpoint.Stop(105)
        } else {
            Endpoint.Say("Malicious file was not caught")
        }
        {{- end}}
    {{- end}}

    // Handling Network calls.
    {{- range $.NetworkCalls}}
        {{- if eq .Function "GET"}}
        // GET function logic for Network.
        Endpoint.Say("Executing GET Request")
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
        // POST function logic for Network.
        Endpoint.Say("Executing POST Request")
        requestOptions := Network.RequestParameters{
            Headers: map[string][]string{"Content-Type": {"application/json"}},
            Body: []byte("{{.Body}}"),
        }
        requester := Network.NewHTTPRequest("{{.URL}}", nil)
        response, err := requester.POST(requestOptions)
        if err != nil {
            Endpoint.Say("POST Error: " + err.Error())
        } else {
            Endpoint.Say("POST Response: " + string(response.Body))
        }
    
        {{- else if eq .Function "TCP"}}
        // TCP function logic for Network.
        Endpoint.Say("Executing TCP connection")
        Network.TCP("{{.Host}}", "{{.Port}}", []byte("{{.Message}}"))
        
        {{- else if eq .Function "UDP"}}
        // UDP function logic for Network.
        Endpoint.Say("Executing UDP connection")
        Network.UDP("{{.Host}}", "{{.Port}}", []byte("{{.Message}}"))
        
        {{- else if eq .Function "ScanPort"}}
        // ScanPort function logic for Network.
        Endpoint.Say("Executing Port Scan")
        isOpen := Network.ScanPort("{{.Protocol}}", "{{.Hostname}}", {{.Port}})
        Endpoint.Say(fmt.Sprintf("ScanPort: Port %d open: %v", {{.Port}}, isOpen))
        
        {{- else if eq .Function "MultiplePortScan"}}
        // MultiplePortScan function logic for Network.
        Endpoint.Say("Executing Multi Port Scan")
        fmt.Println("Ports:", []int{ {{- range $i, $port := .Ports}}{{if $i}}, {{end}}{{ $port }}{{- end}} })
        for _, port := range []int{ {{- range $i, $port := .Ports}}{{if $i}}, {{end}}{{ $port }}{{- end}} } {
            isOpen := Network.ScanPort("{{.Protocol}}", "{{.Hostname}}", port)
            if isOpen {
                fmt.Printf("Port %d is open!\n", port)
            } else {
                fmt.Printf("Port %d is closed!\n", port)
            }
        }
        {{- end}}
    {{- end}}
}

// Function for cleanup operations.
func clean() {
    // Clean logic here.
    Endpoint.Say("Cleaning up")
}

// Main function to start the application.
func main() {
    // Starting the application with necessary endpoint calls.
    {{- range .EndpointCalls}}
        {{- if eq .Function "Start"}}
        Endpoint.Start(test, clean)
        {{- end}}
    {{- end}}
}
`

func main() {
	yamlFilePath := flag.String("yaml", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	yamlFile, err := os.ReadFile(*yamlFilePath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v\n", err)
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v\n", err)
	}

	// Generate UUID for ID if not provided
	if conf.ID == "" {
		conf.ID = uuid.New().String()
	}

	// Generate timestamp for Created if not provided
	if conf.Created == "" {
		conf.Created = time.Now().Format("2006-01-02 15:04:05")
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll("generated_code", 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v\n", err)
	}

	// Create a file with a UTC timestamp in the name
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("generated_code/generated-%s.go", timestamp)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v\n", err)
	}
	defer file.Close()

	// Create and execute the template
	t := template.Must(template.New("goFile").Parse(templateText))
	err = t.Execute(file, conf)
	if err != nil {
		log.Fatalf("Failed to execute template: %v\n", err)
	}

	fmt.Println("File generated:", filename)
}
