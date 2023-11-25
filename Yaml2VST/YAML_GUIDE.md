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

- `id`: A unique identifier for the test.
- `name`: The name of the test.
- `unit`: The unit to which the test belongs.
- `created`: The creation date of the test in YYYY-MM-DD format.

Example:

```yaml
id: "12345"
name: "TestName"
unit: "TestUnit"
created: "2023-01-01"
```

## 2. Imports

The `imports` section allows you to specify external libraries or packages that your test depends on. You can use aliases for imports.

Example:

```yaml
imports:
  - alias: "Endpoint"
    path: "github.com/example/endpoint"
  - alias: "Network"
    path: "github.com/example/network"
  - path: "fmt"
```

## 3. Embedded Files

The `embeddedFiles` section is used to embed external files into the test. These files can be accessed within your code.

Example:

```yaml
embeddedFiles:
  - name: "file1"
    content: |
      //go:embed file1.txt
      var file1 []byte
  - name: "file2"
    content: |
      //go:embed file2.txt
      var file2 []byte
```

## 4. Endpoint Calls

The `endpointCalls` section defines calls to Endpoint functions. Each call has the following properties:

- `function`: The name of the Endpoint function to call.
- `arguments`: A list of arguments for the function.
- `shellCommands`: Optional shell commands to execute.
- `testFunction`: The name of the test function (if applicable).
- `cleanFunction`: The name of the clean-up function (if applicable).

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
```

## 5. Network Calls

The `networkCalls` section defines calls to Network functions. Each call has the following properties:

- `function`: The name of the Network function to call.
- `url`: The URL for the request (if applicable).
- `headers`: HTTP headers (if applicable).
- `queryParams`: Query parameters (if applicable).
- `body`: Request body (if applicable).
- Other function-specific properties.

Example:

```yaml
networkCalls:
  - function: "GET"
    url: "https://example.com"
    headers:
      Content-Type:
        - "application/json"
  - function: "POST"
    url: "https://example.com/api"
    headers:
      Content-Type:
        - "application/json"
    body: "{\"key\": \"value\"}"
```

## 6. Function Documentation

### Supported Functions

Here is documentation for the supported functions and their usage in the test:

- **`Start`**
  - **Arguments**: None
  - **Description**: Start the test and execute the test and clean-up functions.

- **`Say`**
  - **Arguments**: A list of messages to display.
  - **Description**: Display messages in the test.

- **`Shell`**
  - **Arguments**: List of shell commands to execute.
  - **Description**: Execute shell commands within the test.

- **`Stop`**
  - **Arguments**: An error code.
  - **Description**: Stop the test execution with the specified error code.

- **`Find`**
  - **Arguments**: A file type to search for.
  - **Description**: Search for files with the specified file type.

- **`Read`**
  - **Arguments**: The path to the file to read.
  - **Description**: Read the contents of a file.

- **`Write`**
  - **Arguments**: The filename and content to write.
  - **Description**: Write content to a file.

- **`Exists`**
  - **Arguments**: The file path to check.
  -

 **Description**: Check if a file exists.

- **`Quarantined`**
  - **Arguments**: The filename and quarantine type.
  - **Description**: Simulate file quarantine and check for malicious files.

- **`Remove`**
  - **Arguments**: The file path to remove.
  - **Description**: Remove a file.

- **`ExecuteRandomCommand`**
  - **Arguments**: List of random shell commands to execute.
  - **Description**: Execute random shell commands.

- **`IsAvailable`**
  - **Arguments**: List of tools or commands to check for availability.
  - **Description**: Check if tools or commands are available.

- **`IsSecure`**
  - **Arguments**: List of security checks to perform.
  - **Description**: Check if the system is secure.

- **`GET`**
  - **Arguments**: URL, headers, query parameters.
  - **Description**: Execute an HTTP GET request.

- **`POST`**
  - **Arguments**: URL, headers, request body.
  - **Description**: Execute an HTTP POST request.

- **`TCP`**
  - **Arguments**: Host, port, message.
  - **Description**: Execute a TCP connection.

- **`UDP`**
  - **Arguments**: Host, port, message.
  - **Description**: Execute a UDP connection.

- **`ScanPort`**
  - **Arguments**: Protocol, hostname, port.
  - **Description**: Scan a port for availability.

- **`MultiplePortScan`**
  - **Arguments**: List of ports to scan.
  - **Description**: Scan multiple ports for availability.
