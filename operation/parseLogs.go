package operation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type LogEntry struct {
	// Core fields
	Level    string `json:"level"`
	Message  string `json:"msg"`
	Error    string `json:"error,omitempty"`
	Response string `json:"response.message,omitempty"`
	// Harness entity fields
	Connector        string `json:"connector,omitempty"`
	Environment      string `json:"environment,omitempty"`
	EnvironmentGroup string `json:"environment group,omitempty"`
	FeatureFlag      string `json:"feature flag,omitempty"`
	FileStore        string `json:"file store,omitempty"`
	Infrastructure   string `json:"infrastructure,omitempty"`
	InputSet         string `json:"input set,omitempty"`
	Pipeline         string `json:"pipeline,omitempty"`
	Project          string `json:"project,omitempty"`
	ResourceGroup    string `json:"resource group,omitempty"`
	Role             string `json:"role,omitempty"`
	RoleAssignments  string `json:"role assignments,omitempty"`
	Service          string `json:"service,omitempty"`
	ServiceAccount   string `json:"service account,omitempty"`
	ServiceOvericde  string `json:"service override,omitempty"`
	Tags             string `json:"tags,omitempty"`
	TargetGroups     string `json:"target groups,omitempty"`
	Template         string `json:"template,omitempty"`
	Trigger          string `json:"trigger,omitempty"`
	User             string `json:"user,omitempty"`
	UserGroup        string `json:"user group,omitempty"`
	Variables        string `json:"variables,omitempty"`
}

func ParseAndPrintProjectLogs(logs, logLevel, sourceProject string) {
	logEntries := strings.Split(logs, "\n")

	entryCount := len(logEntries)
	entryIndex := int(1)

	fmt.Printf("Logs from migrating project: '%v'. Log Level: '%v'. \n -- \n", sourceProject, logLevel)
	for _, logEntry := range logEntries {
		entryIndex++
		if logEntry == "" {
			continue
		}
		var entry LogEntry
		err := json.Unmarshal([]byte(logEntry), &entry)
		if err != nil {
			fmt.Println("Failed to parse log entry:", err)
			continue
		}
		if entry.Level == logLevel {

			// fmt.Println(entry.Message)

			printLogs(entry, logLevel)

			if entryIndex == entryCount {
				fmt.Println("---")
				fmt.Printf("--- End of log entries for project: '%v' --- \n", sourceProject)
				fmt.Println("---")
			} else {
				fmt.Println("---")
			}
		}
	}
}

func ParseAndPrintGlobalLogs(logs, logLevel string) {
	logEntries := strings.Split(logs, "\n")

	entryCount := len(logEntries)
	entryIndex := int(1)

	fmt.Printf("Logs encountered when during global operations. \n -- \n")
	for _, logEntry := range logEntries {
		entryIndex++
		if logEntry == "" {
			continue
		}
		var entry LogEntry
		err := json.Unmarshal([]byte(logEntry), &entry)
		if err != nil {
			fmt.Println("Failed to parse log entry:", err)
			continue
		}
		if entry.Level == logLevel {

			// fmt.Println(entry.Message)

			printLogs(entry, logLevel)

			if entryIndex == entryCount {
				fmt.Println("---")
				fmt.Printf("--- End of log entries for global opreration: --- \n")
				fmt.Println("---")
			} else {
				fmt.Println("---")
			}
		}
	}
}

func printLogs(entry LogEntry, logLevel string) {
	v := reflect.ValueOf(entry)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Skip empty fields
		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			continue
		}

		// Print field name and value
		switch logLevel {
		case "error":
			fmt.Printf(Red + "%s: %v\n" + Reset, fieldType.Name, fieldValue.Interface())
		case "warn":
			fmt.Printf(Yellow + "%s: %v\n" + Reset, fieldType.Name, fieldValue.Interface())
		case "info":
			fmt.Printf("%s: %v\n", fieldType.Name, fieldValue.Interface())
		}
	}
}
