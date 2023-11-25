/*
ID: 12345
NAME: TestName
UNIT: TestUnit
CREATED: 2023-01-01
*/
package main
import (
        Endpoint "github.com/preludeorg/libraries/go/tests/endpoint"
        Network "github.com/preludeorg/libraries/go/tests/network"
        "fmt"
)
//go:embed fromrussiawithlove.eml
var emlFile []byte

//go:embed fromrussiawithlove.pdf
var pdfFile []byte


func test() {
    // Endpoint calls
            Endpoint.Say("Starting combined Endpoint and Network tests")
                command := []string{"cmd.exe", "/C", "wmic path win32_logicaldisk get caption,filesystem,freespace,size,volumename" }
                Endpoint.Shell(command)
            Endpoint.Say("Starting scan for files")
            fileType := ".txt"
            files := Endpoint.Find(fileType)
            if len(files) == 0 {
                Endpoint.Stop(104)
            }
            data := Endpoint.Read("path/to/file.txt")
            Endpoint.Write("filename.txt", []byte("Hello, World!"))
            exists := Endpoint.Exists("path/to/file.txt")
            isQuarantined := Endpoint.Quarantined("quarantine.txt", []byte("Sample Content"))
            success := Endpoint.Remove("path/to/remove.txt")
            commands := [][]string{{"ls", "-la" },{"echo", "Hello World" } }
            output, err := Endpoint.ExecuteRandomCommand(commands)
            available := Endpoint.IsAvailable("bash", "python")
            secure := Endpoint.IsSecure()

    // Network calls
            Endpoint.Say("Executing GET Request")
            requestOptions := Network.RequestParameters{
                Headers: map[string][]string{"Content-Type": {"application/json"}},
                QueryParams: map[string][]string{
                        "key": {"value" },
                },
                Body: []byte(""),
            }
            requester := Network.NewHTTPRequest("https://example.com", nil)
            response, err := requester.GET(requestOptions)
            if err != nil {
                Endpoint.Say("GET Error: " + err.Error())
            } else {
                Endpoint.Say("GET Response: " + string(response.Body))
            }
            Endpoint.Say("Executing POST Request")
            Endpoint.Say("Executing TCP connection")
            Network.TCP("127.0.0.1", "8080", []byte("Hello TCP"))
            Endpoint.Say("Executing UDP connection")
            Network.UDP("10.0.0.1", "8081", []byte("Hello UDP"))
            Endpoint.Say("Executing Port Scan")
            isOpen := Network.ScanPort("tcp", "localhost", 80)
            Endpoint.Say(fmt.Sprintf("ScanPort: Port %d open: %v", 80, isOpen))
            Endpoint.Say("Executing Multi Port Scan")
            fmt.Println("Ports:", []int{22, 80, 443 })
            for _, port := range []int{22, 80, 443 } {
                isOpen := Network.ScanPort("", "", port)
                if isOpen {
                    fmt.Printf("Port %d is open!\\n", port)
                } else {
                    fmt.Printf("Port %d is closed!\\n", port)
                }
            }
}

func clean() {
    Endpoint.Say("Cleaning up")
}

func main() {
            Endpoint.Start(test, clean)
}
