git s# Kafka Partition Calculator CLI

This is a command-line tool called `partikey` written in Golang for calculating the Kafka partition for a given key using both Murmur2 and Murmur3 hashing algorithms. 
It connects to a Kafka broker, retrieves partition information for a specified topic, and computes which partition a given key would be assigned to according to the two hashing strategies. 
The tool also includes special logic for calculating partitions for the `__consumer_offsets` topic.


## Features

- **Kafka Partition Calculation**: Calculate the partition for a Kafka key using Murmur2 and Murmur3 hashing algorithms.
- **Default Murmur2 Hashing**: By default, the tool uses the Murmur2 hashing algorithm, which is standard for Kafka.
- **Optional Murmur3 Hashing**: You can override the default and use Murmur3 hashing instead.
- **Custom Partition Logic**: Special logic for calculating partitions for the `__consumer_offsets` topic.
- **Command-Line Interface**: Easy-to-use CLI with options to specify the Kafka topic, key, broker, and hash algorithm.

## Prerequisites

- **Go**: You need Go installed on your machine to build and run the program.
  ```sh
  brew install go
  ```

- **Kafka**: A running Kafka cluster to connect to. 

## Installation for the local development environment

- Clone the repository:
  ```sh
  git clone git@github.com:praveenkumaresan/gogarage.git
  cd cmd/partikey
  ```

- Build the Go binary:
  ```sh
  go build -o partikey main.go
  ```

## Usage

- Run the tool with the required command-line flags:
  ```sh
  ./partikey --topic <topic_name> --key <key_value> --broker <broker_address>
  ```

### Example Usage

- Calculating Partition for a Standard Kafka Topic
  ```sh
  ./partikey --topic=build-goland-over-weekend --key=kafka-is-amazing-tool --broker=foo:9093
  ```
Output:
  ```sh
   Using Murmur2: The key 'kafka-is-amazing-tool' belongs to partition 5
  ```  

- Overriding with Murmur3 Hashing
  ```sh
  ./partikey --topic=build-goland-over-weekend --key=kafka-is-amazing-tool --broker=foo:9093 --hash=Murmur3
  ```
Output:
  ```sh
  Using Murmur3: The key 'kafka-is-amazing-tool' belongs to partition 18
  ```  
  
- Special Handling for __consumer_offsets Topic
  ```sh
   ./partikey --topic=__consumer_offsets --key=build-goland-over-weekend --broker=foo:9093
  ```
Output:
  ```sh
    Using StringHashCode: The key 'build-goland-over-weekend' belongs to partition 7
  ```  

## Command-Line Options
    --topic: The Kafka topic to query.
    --key: The Kafka partition key to calculate the partition for.
    --broker: The Kafka broker address (host).
    --hash: The hash algorithm to use (Murmur2 or Murmur3). Default is Murmur2.
    --help: Display help message with usage instructions.

## Notes

- When using the tool for the __consumer_offsets topic, the partition calculation uses a special string hash code logic with a fixed partition count of 50.
Ensure that the Kafka broker address is correctly specified to connect to your Kafka cluster.