# Transbuilder

**Transbuilder** is an automated translation tool developed in Go that uses the ChatGPT API to automatically translate various file types into one or more target languages. It supports popular file formats such as XLIFF, JSON, YAML, and CSV.

## Features

- Automatically translates files using the ChatGPT API.
- Supports multiple file formats:
  - **XLIFF** (`.xlf`, `.xliff`)
  - **YAML** (`.yaml`, `.yml`)
  - **JSON** (`.json`)
  - **CSV** (`.csv`)
- Allows translation into multiple target languages at once.
- Generates translated files with proper formatting for each file type.

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) version 1.15 or later.
- An API key for the ChatGPT API (OpenAI).

### Installation Steps

1. Clone the project repository:
    ```bash
    git clone https://github.com/ICTools/transbuilder.git
    cd transbuilder
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Build the project:
    ```bash
    go build -o transbuilder ./cmd
    ```

## Usage

### Arguments

- `--file`: Path to the file to be translated (required).
- `--langs`: Comma-separated list of target languages (e.g., `fr,de,es`) (required).
- `--api-key`: Your ChatGPT API key (required).

### Example Usage

1. To translate an XLIFF file into French and German:

    ```bash
    ./transbuilder --file input.xliff --langs fr,de --api-key YOUR_API_KEY
    ```

2. To translate a JSON file into Spanish:

    ```bash
    ./transbuilder --file input.json --langs es --api-key YOUR_API_KEY
    ```

The translated files will be generated with the original input file name, followed by the target language code, for example: `input_fr.xliff`.

### Supported File Formats

| Format   | Extensions       | Description                                   |
|----------|------------------|-----------------------------------------------|
| XLIFF    | `.xlf`, `.xliff`  | Localization exchange file in XML             |
| YAML     | `.yaml`, `.yml`   | Readable text format for configuration or data|
| JSON     | `.json`           | JSON data file                               |
| CSV      | `.csv`            | Spreadsheet-like data file, each cell can be translated |

## Contributing

If you'd like to contribute to this project, please submit a pull request or open an issue on the GitHub repository.