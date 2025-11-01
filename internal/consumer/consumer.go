package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/imvalerio/transcribeit/internal/utils"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	topic         string
	brokerAddress string
	log           *log.Logger
}

func NewConsumer(topic string, brokerAddress string, logger *log.Logger) *Consumer {
	return &Consumer{
		topic:         topic,
		brokerAddress: brokerAddress,
		log:           logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	fmt.Println("...Kafka consumer starting...")

	// Create a new reader with the given broker addresses, topic, and group ID.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{c.brokerAddress},
		Topic:    c.topic,
		GroupID:  "my-consumer-group",
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})

	defer func() {
		c.log.Printf("Consumer for topic '%s' stopping...", c.topic)
		if err := reader.Close(); err != nil {
			c.log.Printf("Error closing consumer for topic '%s': %v", c.topic, err)
		}
	}()

	c.log.Printf("🚀 Kafka consumer started on %s topic\n", c.topic)

	for {
		select {
		case <-ctx.Done():
			c.log.Printf("Context canceled, shutting down consumer '%s'", c.topic)
			return nil

		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				c.log.Printf("❌ error while reading message: %v", err)
				continue
			}

			inputPath := string(m.Value)
			c.log.Printf("✅ received message: topic=%s partition=%d offset=%d\n%s\n\n",
				m.Topic, m.Partition, m.Offset, inputPath)

			if c.topic == "upload-audio" {

				outputPath, err := filepath.Abs("transcriptions")
				if err != nil {
					c.log.Printf("❌ error while getting absolute path for output: %v", err)
					continue
				}
				c.transcribe(inputPath, outputPath)

			}
		}
	}
}

func (c *Consumer) transcribe(inputPath, outputPath string) error {

	c.log.Printf(" transcribing into %s...", outputPath)

	projectPath := os.Getenv("VOICE_RECOGNITION_FOLDER")
	c.log.Printf("Project path: %s", projectPath)

	pythonExe := os.Getenv("PYTHON_EXEC_PATH")
	c.log.Printf("Python exe path: %s", pythonExe)

	// Path to the script you want to run
	scriptPath := projectPath + `\main.py`

	cmd := exec.Command(pythonExe, scriptPath, inputPath, outputPath)
	cmd.Dir = projectPath // set working directory

	out, err := cmd.CombinedOutput()
	c.log.Println(string(out))

	if err == nil {
		c.log.Print("✅ transcription completed")

		utils.PushToQueue("completed-transcriptions", outputPath)
	}

	return err
}
