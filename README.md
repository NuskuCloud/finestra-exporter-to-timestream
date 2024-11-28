# README

## Project Overview

This project is designed to fetch and process data from the Finestra API provided by Salford University, and then insert that data into AWS Timestream for storage and analysis (Grafana!).

## Prerequisites

- Go 1.23.1 or later
- AWS Credentials (Access Key and Secret Key)
- Finestra API credentials (username and password provided by Salford University)
- Internet connection

## Setup

### 1. Clone the Repository

First, clone the repository to your local machine:

```sh
git clone https://github.com/nuskucloud/finestra-exporter-to-timestream.git
cd finestra-exporter-to-timestream
```

### 2. Install Dependencies

```sh
go mod tidy
```

### 3. Running the Application

```sh
go run main.go --finestra_username=<username> --finestra_password=<password> --aws_key=<aws-access-key> --aws_secret=<aws-secret> --timestream_database=<database-name> --timestream_table=<table-name> --finestra_location_id=<location-id>
```

### Command Line Flags

- `--finestra_username`: Your Finestra username (provided by Salford University).
- `--finestra_password`: Your Finestra password (provided by Salford University).
- `--aws_key`: AWS access key.
- `--aws_secret`: AWS secret key.
- `--timestream_database`: Name of the Timestream database.
- `--timestream_table`: Name of the Timestream table.
- `--finestra_location_id`: Location ID for fetching location data.

## File Breakdown

### `main.go`

The main entry point for the application. It handles argument parsing, authentication, and data processing.

### `aws_config.go`

Responsible for setting up AWS Timestream service configuration and HTTP transport settings.

### `aws_timestream.go`

Handles inserting data into AWS Timestream.

### `csv.go`

Provides functions for parsing CSV data and processing each row.

### `finestra_api_client.go`

Handles API calls to the Finestra API, including authentication and data export.

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests.

## License
This project is licensed under the MIT License.