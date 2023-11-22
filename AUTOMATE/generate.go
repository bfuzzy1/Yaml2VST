package main

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

type ImportConfig struct {
	Alias string `yaml:"alias,omitempty"`
	Path  string `yaml:"path"`
}

type Config struct {
	ID            string         `yaml:"id"`
	Name          string         `yaml:"name"`
	Unit          string         `yaml:"unit"`
	Created       string         `yaml:"created"`
	Imports       []ImportConfig `yaml:"imports,omitempty"`
	EndpointCalls []EndpointCall `yaml:"endpointCalls"`
	NetworkCalls  []NetworkCall  `yaml:"networkCalls"`
}

// EndpointCall represents a single call to an Endpoint function
type EndpointCall struct {
	Function  string   `yaml:"function"`
	Arguments []string `yaml:"arguments,omitempty"`

	// For file operations
	FilePath    string `yaml:"filePath,omitempty"`
	FileContent string `yaml:"fileContent,omitempty"` // Could be base64 encoded for binary data
	FileType    string `yaml:"fileType,omitempty"`    // e.g., for filtering in Find function

	// For encryption/decryption
	EncryptionKey string `yaml:"encryptionKey,omitempty"`
	DataToEncrypt string `yaml:"dataToEncrypt,omitempty"`
	DataToDecrypt string `yaml:"dataToDecrypt,omitempty"`

	// For shell commands, if multiple commands are to be executed
	ShellCommands []string `yaml:"shellCommands,omitempty"`

	// Special functions
	TestFunction  string `yaml:"testFunction,omitempty"`
	CleanFunction string `yaml:"cleanFunction,omitempty"`
	ErrorCode     int    `yaml:"errorCode,omitempty"` // For use with Stop function
}

// NetworkCall represents a single call to a Network function
type NetworkCall struct {
	Function    string              `yaml:"function"`
	URL         string              `yaml:"url,omitempty"`
	Method      string              `yaml:"method,omitempty"` // GET, POST, HEAD, DELETE
	Headers     map[string][]string `yaml:"headers,omitempty"`
	QueryParams map[string][]string `yaml:"queryParams,omitempty"`
	Body        string              `yaml:"body,omitempty"`
	Encoding    string              `yaml:"encoding,omitempty"` // e.g., "gzip"

	// Authentication
	AuthType   string `yaml:"authType,omitempty"` // e.g., "Basic", "Bearer"
	Credential string `yaml:"credential,omitempty"`

	// Request Options
	Timeout   int    `yaml:"timeout,omitempty"` // in seconds
	UserAgent string `yaml:"userAgent,omitempty"`

	// TCP/UDP Specific
	Host    string `yaml:"host,omitempty"`
	Port    string `yaml:"port,omitempty"`
	Message string `yaml:"message,omitempty"`

	// Port Scanning
	Protocol string `yaml:"protocol,omitempty"` // tcp, udp
	Hostname string `yaml:"hostname,omitempty"`
	Ports    []int  `yaml:"ports,omitempty"` // for scanning multiple ports

	// Internal IP retrieval (no additional fields needed)
}

// Template text for generating the Go file
const templateText = `/*
ID: {{.ID}}
NAME: {{.Name}}
UNIT: {{.Unit}}
CREATED: {{.Created}}
*/
package main

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
        // Similar implementation for POST
	Endpoint.Say("Executing POST Request")
    {{- else if eq .Function "TCP"}}
    	Endpoint.Say("Executing TCP connection")
        Network.TCP("{{.Host}}", "{{.Port}}", []byte("{{.Message}}"))
    {{- else if eq .Function "UDP"}}
    	Endpoint.Say("Executing UDP connection")
        Network.UDP("{{.Host}}", "{{.Port}}", []byte("{{.Message}}"))
    {{- else if eq .Function "ScanPort"}}
    	Endpoint.Say("Executing Port Scan")
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
	filename := fmt.Sprintf("generated_code/generated-%s.go", timestamp)
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
