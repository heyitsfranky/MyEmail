package MyEmail

import (
	"testing"
)

const CONFIG_DIR_PATH = "configs/"

func Test_Init(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		expectErr bool
	}{
		{"ValidConfig", "config.yaml", false},
		{"MissingClientOrigin", "missing_clientorigin.yaml", true},
		{"MissingKafkaBroker", "missing_kafkabroker.yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Data = nil
			err := Init(CONFIG_DIR_PATH + tt.config)
			if err == nil && tt.expectErr {
				t.Errorf("Expected an error but got nil")
			} else if err != nil && !tt.expectErr {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func Test_CreateSendEmailEvent(t *testing.T) {
	Init(CONFIG_DIR_PATH + "config.yaml")
	testEmail := EmailData{ID: "testID", Operation: "testOperation", TargetAddress: "hello@world.com", Subject: "hey!", Content: map[string]interface{}{"test": "hi", "again": "yo", "numbers": 123}}
	tests := []struct {
		name      string
		email     EmailData
		async     bool
		expectErr bool
	}{
		{"ValidEventSync", testEmail, false, false},
		{"ValidEventAsync", testEmail, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateSendEmailEvent(tt.name, tt.email.ID, tt.email.Operation, tt.email.TargetAddress, tt.email.Content, tt.async)
			if err == nil && tt.expectErr {
				t.Errorf("Expected an error but got nil")
			} else if err != nil && !tt.expectErr {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
