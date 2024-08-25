package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/spaolacci/murmur3"
)

func main() {
	// Initialize logging
	initializeLogging()

	// Add a check for the --help flag or no arguments
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "--help") {
		printUsage()
		return
	}

	// Parse command-line flags
	topic, partitionKey, broker, hashAlgorithm := parseFlags()

	// Check if the topic is __consumer_offsets
	if topic == "__consumer_offsets" {
		partition := calculateConsumerOffsetsPartition(partitionKey)
		fmt.Printf("Using StringHashCode: The partition key '%s' belongs to partition %d\n", partitionKey, partition)
	} else {
		// Get partitions for the topic using sarama
		partitions := getTopicPartitions(broker, topic)

		// Calculate and display partitions for the partition key
		displayPartitionResults(partitionKey, partitions, hashAlgorithm)
	}
}

// initializeLogging sets up logging for the application.
func initializeLogging() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// printUsage displays the usage instructions for the CLI tool.
func printUsage() {
	fmt.Println("Usage: partikey --topic <topic> --partitionkey <partitionkey> --broker <broker> [--hash Murmur2|Murmur3]")
	fmt.Println("\nOptions:")
	fmt.Println("  --topic         The Kafka topic to query.")
	fmt.Println("  --partitionkey  The Kafka partition key to calculate the partition for.")
	fmt.Println("  --broker        The Kafka broker address (host:port).")
	fmt.Println("  --hash          The hash algorithm to use (Murmur2 or Murmur3). Default is Murmur2.")
	fmt.Println("  --help          Display this help message.")
}

// parseFlags parses and validates command-line flags.
func parseFlags() (string, string, string, string) {
	topic := flag.String("topic", "", "Kafka topic to query")
	partitionKey := flag.String("partitionkey", "", "Kafka partition key to calculate the partition for")
	broker := flag.String("broker", "", "Kafka broker address (host:port)")
	hashAlgorithm := flag.String("hash", "murmur2", "Hash algorithm to use (Murmur2 or Murmur3). Default is Murmur2.")

	flag.Parse()

	if *topic == "" || *partitionKey == "" || *broker == "" {
		log.Fatalf("All -topic, -partitionkey, and -broker must be specified")
	}

	return *topic, *partitionKey, *broker, *hashAlgorithm
}

// calculateConsumerOffsetsPartition calculates the partition for the given partition key using the special logic for __consumer_offsets.
func calculateConsumerOffsetsPartition(groupID string) int {
	consumerOffsetsPartitionCount := 50
	partition := abs(stringHashCode(groupID) % int32(consumerOffsetsPartitionCount))
	return int(partition) // Convert int32 to int before returning
}

// stringHashCode calculates a simple hash code for a string.
func stringHashCode(s string) int32 {
	var h int32 = 0
	for _, ch := range s {
		h = 31*h + ch
	}
	return h
}

// abs returns the absolute value of an integer.
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func getTopicPartitions(broker, topic string) int {
	// Replace port 9093 with 9092 in the broker string
	broker = replacePort(broker, "9093", "9092")

	config := sarama.NewConfig()

	client, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer func(client sarama.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("Failed to close Kafka client: %v", err)
		}
	}(client)

	partitions, err := client.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions for topic: %v", err)
	}

	return len(partitions)
}

// replacePort replaces the old port with a new port in the broker string.
func replacePort(broker, oldPort, newPort string) string {
	return strings.Replace(broker, ":"+oldPort, ":"+newPort, 1)
}

// displayPartitionResults calculates and displays the partition for the given partition key using the specified hash algorithm.
func displayPartitionResults(partitionKey string, numPartitions int, hashAlgorithm string) {
	var partition int

	if hashAlgorithm == "Murmur3" {
		partition = calculatePartitionMurmur3(partitionKey, numPartitions)
		fmt.Printf("Using Murmur3: The partition key '%s' belongs to partition %d\n", partitionKey, partition)
	} else {
		partition = calculatePartitionMurmur2(partitionKey, numPartitions)
		fmt.Printf("Using Murmur2: The partition key '%s' belongs to partition %d\n", partitionKey, partition)
	}
}

// calculatePartitionMurmur2 calculates the partition for the given partition key using the Murmur2 hashing algorithm.
func calculatePartitionMurmur2(partitionKey string, numPartitions int) int {
	hashValue := murmur2([]byte(partitionKey))
	return int(hashValue) % numPartitions
}

// calculatePartitionMurmur3 calculates the partition for the given partition key using the Murmur3 hashing algorithm.
func calculatePartitionMurmur3(partitionKey string, numPartitions int) int {
	hashValue := murmur3.Sum32([]byte(partitionKey))
	partition := int(int32(hashValue)) % numPartitions

	// Ensure partition is non-negative
	if partition < 0 {
		partition = -partition
	}

	return partition
}

// murmur2 is an implementation of the Murmur2 hashing algorithm as used in Kafka.
func murmur2(data []byte) uint32 {
	const m uint32 = 0x5bd1e995
	const r = 24

	length := len(data)
	h := uint32(length)

	for length >= 4 {
		k := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		k *= m
		k ^= k >> r
		k *= m

		h *= m
		h ^= k

		data = data[4:]
		length -= 4
	}

	switch length {
	case 3:
		h ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint32(data[0])
		h *= m
	}

	h ^= h >> 13
	h *= m
	h ^= h >> 15

	return h & 0x7fffffff
}
