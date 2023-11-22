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
	"time"
)

func test() {
    // Endpoint calls
    Endpoint.Say("Starting combined Endpoint and Network tests")
    Endpoint.Shell("ls", "-l")

    // Network calls
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
        // Similar implementation for POST
        Network.TCP("127.0.0.1", "8080", []byte("Hello TCP"))
		Network.UDP("127.0.0.1", "8081", []byte("Hello UDP"))
        isOpen := Network.ScanPort("tcp", "localhost", 80)
        Endpoint.Say(fmt.Sprintf("ScanPort: Port %d open: %v", 80, isOpen))
}

func clean() {
    Endpoint.Say("Cleaning up")
}

func main() {
    Endpoint.Start(test, clean)
}
