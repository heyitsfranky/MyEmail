package MyEmail

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/heyitsfranky/MyConfig"
	"github.com/segmentio/kafka-go"
)

var Data *InitData

type InitData struct {
	KafkaBroker  string `yaml:"kafka-broker"`
	ClientOrigin string `yaml:"client-origin"`
}

type EmailData struct {
	ID            string                 `json:"id"`
	Operation     string                 `json:"operation"`
	TargetAddress string                 `json:"target_address"`
	Subject       string                 `json:"subject"`
	Content       map[string]interface{} `json:"content"`
}

func Init(configPath string) error {
	if Data == nil {
		err := MyConfig.Init(configPath, &Data)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateSendEmailEvent(id string, operation string, targetAddress string, subject string, content map[string]interface{}, async bool) error {
	if Data == nil {
		return fmt.Errorf("must first call Init() to initialize the kafka settings")
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{Data.KafkaBroker},
		Topic:    "send_email",
		Balancer: &kafka.LeastBytes{},
	})
	if !async { //necessary as otherwise the writer is closed as soon as the function ends
		defer writer.Close()
	}
	emailData := EmailData{
		ID:            id,
		Operation:     operation,
		TargetAddress: targetAddress,
		Subject:       subject,
		Content:       content,
	}
	jsonString, err := json.Marshal(emailData)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return err
	}
	message := kafka.Message{
		Key:   []byte(fmt.Sprintf("key-%d", time.Now().Unix())),
		Value: []byte(jsonString),
	}
	if async {
		go func() {
			err := writer.WriteMessages(context.Background(), message)
			if err != nil {
				fmt.Println("Error sending message asynchronously:", err)
			}
			writer.Close()
		}()
	} else {
		err := writer.WriteMessages(context.Background(), message)
		if err != nil {
			fmt.Println("Error sending message synchronously:", err)
			return err
		}
	}
	return nil
}
