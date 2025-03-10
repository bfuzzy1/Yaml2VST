# YAML Configuration Guide for Yaml2VST

This guide provides comprehensive documentation on how to structure YAML configuration files for our tests. YAML configuration files allow you to define various settings and specify actions to be performed within the test.

## Table of Contents

1. [General Configuration](#1-general-configuration)
2. [Imports](#2-imports)
3. [Embedded Files](#3-embedded-files)
4. [Endpoint Calls](#4-endpoint-calls)
5. [Network Calls](#5-network-calls)
6. [Function Documentation](#6-function-documentation)

---

## 1. General Configuration

The general configuration section defines basic information about the test.

- `id`: A unique identifier for the test. (Optional - will be auto-generated if not provided)
- `name`: The name of the test. (Required)
- `unit`: The unit to which the test belongs. (Required)
- `created`: The creation date and time. (Optional - will be auto-generated if not provided)

Example:

```yaml
id: "12345"
name: "TestName"
unit: "TestUnit"
created: "2023-01-01 12:34:56"
```

## 2. Imports

The `imports` section allows you to specify external libraries or packages that your test depends on. You can use aliases for imports.

Example:

```yaml
imports:
  - alias: "Endpoint"
    path: "github.com/preludeorg/libraries/go/tests/endpoint"
  - alias: "Network"
    path: "github.com/preludeorg/libraries/go/tests/network"
  - path: "fmt"  # No alias needed for standard packages
```

## 3. Embedded Files

The `embeddedFiles` section is used to embed external files into the test. These files can be accessed within your code.

Example:

```yaml
embeddedFiles:
  - name: "emlFile"
    content: |
      //go:embed fromrussiawithlove.eml
      var emlFile []byte
  - name: "pdfFile"
    content: |
      //go:embed fromrussiawithlove.pdf
      var pdfFile []byte
```

## 4. Endpoint Calls

The `endpointCalls` section defines calls to Endpoint functions. Each call has the following properties:

- `function`: The name of the Endpoint function to call.
- `arguments`: A list of arguments for the function.
- `shellCommands`: Optional shell commands to execute (for Shell and ExecuteRandomCommand functions).
- `testFunction`: The name of the test function (if applicable, used by Start).
- `cleanFunction`: The name of the clean-up function (if applicable, used by Start).

Example:

```yaml
endpointCalls:
  - function: "Start"
    testFunction: "test"
    cleanFunction: "clean"
  - function: "Say"
    arguments:
      - "Starting combined Endpoint and Network tests"
  - function: "Shell"
    shellCommands:
      - ["cmd.exe", "/C", "echo Hello, World"]
  - function: "ExecuteRandomCommand"
    shellCommands:
      - ["ls", "-la"]
      - ["echo", "Hello World"]
```

## 5. Network Calls

The `networkCalls` section defines calls to Network functions. Each call has specific properties depending on the function type:

### Common Network Properties

- `function`: The name of the Network function to call.

### HTTP Request Properties (GET, POST)

- `url`: The URL for the request. (Required)
- `headers`: HTTP headers as a map of string arrays.
- `queryParams`: Query parameters as a map of string arrays.
- `body`: Request body as a string.
- `encoding`: Content encoding. (Optional)
- `authType`: Authentication type. (Optional)
- `credential`: Authentication credential. (Optional)
- `timeout`: Request timeout in seconds. (Optional)
- `userAgent`: User-Agent header. (Optional)

### Socket Communication Properties (TCP, UDP)

- `host`: The host to connect to. (Required)
- `port`: The port to connect to. (Required)
- `message`: The message to send. (Required)

### Port Scanning Properties (ScanPort, MultiplePortScan)

- `protocol`: The protocol to use ("tcp" or "udp"). (Required)
- `hostname`: The hostname to scan. (Required)
- `port`: The single port to scan (for ScanPort). (Required)
- `ports`: An array of ports to scan (for MultiplePortScan). (Required)

Example:

```yaml
networkCalls:
  - function: "GET"
    url: "https://example.com"
    headers:
      Content-Type: ["application/json"]
    queryParams:
      key: ["value"]
  - function: "POST"
    url: "https://example.com/api"
    headers:
      Content-Type: ["application/json"]
    body: "{\"key\": \"value\"}"
  - function: "TCP"
    host: "127.0.0.1"
    port: "8080"
    message: "Hello TCP"
  - function: "UDP"
    host: "10.0.0.1"
    port: "8081"
    message: "Hello UDP"
  - function: "ScanPort"
    protocol: "tcp"
    hostname: "localhost"
    port: 80
  - function: "MultiplePortScan"
    protocol: "tcp"
    hostname: "localhost"
    ports:
      - 22
      - 80
      - 443
```

## 6. Function Documentation

### Supported Functions

Here is documentation for the supported functions and their usage in the test:

#### Endpoint Functions

- **`Start`**
  - **Description**: Initializes the test and executes the test and clean-up functions.
  - **Properties**: 
    - `testFunction`: The name of the function to execute for testing.
    - `cleanFunction`: The name of the function to execute for cleanup.

- **`Say`**
  - **Description**: Displays messages in the test log.
  - **Arguments**: A message string to display.

- **`Shell`**
  - **Description**: Executes shell commands within the test.
  - **Properties**:
    - `shellCommands`: A list of command arrays, where each array represents a command and its arguments.

- **`Stop`**
  - **Description**: Stops the test execution with the specified error code.
  - **Arguments**: An error code (integer).

- **`Find`**
  - **Description**: Searches for files with the specified file type.
  - **Arguments**: A file extension to search for (e.g., ".txt").
  - **Behavior**: Stops with code 104 if no files are found.

- **`Read`**
  - **Description**: Reads the contents of a file.
  - **Arguments**: The path to the file to read.

- **`Write`**
  - **Description**: Writes content to a file.
  - **Arguments**: 
    - The filename to write to.
    - The content to write.

- **`Exists`**
  - **Description**: Checks if a file exists.
  - **Arguments**: The file path to check.

- **`Quarantined`**
  - **Description**: Simulates file quarantine and checks for malicious files.
  - **Arguments**: The filename to check.
  - **Behavior**: Stops with code 105 if a malicious file is detected.

- **`Remove`**
  - **Description**: Removes a file.
  - **Arguments**: The file path to remove.

- **`ExecuteRandomCommand`**
  - **Description**: Executes random shell commands from a provided list.
  - **Properties**:
    - `shellCommands`: A list of command arrays to randomly choose from.

- **`IsAvailable`**
  - **Description**: Checks if tools or commands are available.
  - **Arguments**: List of tools or commands to check.

- **`IsSecure`**
  - **Description**: Performs security checks on the system.
  - **Arguments**: None.

#### Network Functions

- **`GET`**
  - **Description**: Executes an HTTP GET request.
  - **Properties**: See HTTP Request Properties above.

- **`POST`**
  - **Description**: Executes an HTTP POST request.
  - **Properties**: See HTTP Request Properties above.

- **`TCP`**
  - **Description**: Establishes a TCP connection and sends a message.
  - **Properties**: See Socket Communication Properties above.

- **`UDP`**
  - **Description**: Establishes a UDP connection and sends a message.
  - **Properties**: See Socket Communication Properties above.

- **`ScanPort`**
  - **Description**: Scans a single port to check if it's open.
  - **Properties**: See Port Scanning Properties above.

- **`MultiplePortScan`**
  - **Description**: Scans multiple ports to check if they're open.
  - **Properties**: See Port Scanning Properties above.
  - **Note**: Must include protocol, hostname, and an array of ports.
