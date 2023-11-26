# Yaml2VST

[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](LICENSE)
![CodeQL](https://github.com/bfuzzy1/Yaml2VST/workflows/CodeQL/badge.svg)

Yaml2VST is a tool for converting YAML configurations into [Verified Security Test](https://www.preludesecurity.com/products/detect) (VST) code. It simplifies the process of generating security tests based on YAML-defined configurations, making it easy to automate and verify security-related tasks.

## Usage

Yaml2VST is straightforward to use. You provide a YAML configuration file, and it generates VST code based on the hardcoded template provided in the source code. Here's an example of how to use it:

```shell
./yaml2vst -yaml config.yaml
```

- `-yaml`: Specify the path to your YAML configuration file.

## Examples

To help you get started quickly, here's an example YAML configuration:

[Example config.yaml](https://github.com/bfuzzy1/Yaml2VST/blob/main/Yaml2VST/config.yaml)

This YAML configuration, when processed by Yaml2VST, will generate VST code based on on a hardcoded template defined in the source code.

## Installation

To get started with Yaml2VST, you need to install it on your system. You can do this using the following steps:

1. Clone the repository:

   ```shell
   git clone https://github.com/bfuzzy1/Yaml2VST.git
   ```

2. Build the binary:

   ```shell
   cd Yaml2VST
   go build
   ```

3. Run the binary:

   ```shell
   ./yaml2vst
   ```

## License

Yaml2VST is open-source and available under the [BSD 3 Clause License](LICENSE).

## Contact

If you have any questions, suggestions, or feedback, feel free to reach out.

- GitHub Issues: [Submit an Issue](https://github.com/bfuzzy1/Yaml2VST/issues)
