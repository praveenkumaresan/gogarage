# go-garage

go utilities repo for ephemeral golang code/cli tools

## Installation for the local development environment

- Clone the repository:
  ```sh
  git clone https://github.com/praveenkumaresan/gogarage.git
    cd gogarage
    ```

- Build the Go binary:
  ```sh
  go build -o gogarage main.go
  ```
    
## Usage

- Run the tool with the required command-line flags:
  ```sh
  ./gogarage --topic <topic_name> --key <key_value> --broker <broker_address>
  ```
  