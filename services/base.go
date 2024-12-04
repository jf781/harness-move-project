package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
	"harness-copy-project/model"
)

// const BaseURL = "https://app.harness.io"

type ApiRequest struct {
	Client  *resty.Client
	Token   string
	Account string
	BaseURL string
}

type Operation interface {
	Copy() error
}

func updateYamlKeyValues(node *yaml.Node, updates map[string]interface{}) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, contentNode := range node.Content {
			updateYamlKeyValues(contentNode, updates)
		}
	case yaml.MappingNode:
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			if newValue, exists := updates[keyNode.Value]; exists {
				valueNode.Value = fmt.Sprintf("%v", newValue)
			} else {
				updateYamlKeyValues(valueNode, updates)
			}
		}
	case yaml.SequenceNode:
		for _, itemNode := range node.Content {
			updateYamlKeyValues(itemNode, updates)
		}
	case yaml.AliasNode:
		updateYamlKeyValues(node.Alias, updates)
	}
}

func updateYaml(inputYaml, targetOrg, targetProject string) string {
	var data yaml.Node

	yamlBytes := []byte(inputYaml)

	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		panic(fmt.Errorf("failed to parse YAML: %v", err))
	}

	updatedKeys := map[string]interface{}{
		"orgIdentifier":     targetOrg,
		"projectIdentifier": targetProject,
	}

	updateYamlKeyValues(&data, updatedKeys)

	var output strings.Builder
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(&data); err != nil {
		panic(fmt.Errorf("failed to marshal updated YAML: %v", err))
	}
	encoder.Close()

	return output.String()
}

func handleErrorResponse(resp *resty.Response) error {
	result := model.ErrorResponse{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}
	if result.Code == "DUPLICATE_FIELD" {
		return nil
	}
	if strings.Contains(result.Message, "already exists") {
		return nil
	}
	return fmt.Errorf("%s: %s", result.Code, removeNewLine(result.Message))
}

func removeNewLine(value string) string {
	return strings.ReplaceAll(value, "\n", "")
}

func reportFailed(failed []string, description string) {
	if len(failed) > 0 {
		fmt.Println(color.RedString(fmt.Sprintf("Failed %s %d", description, len(failed))))
		fmt.Println(color.RedString(strings.Join(failed, "\n")))
	}
}

func ValidateCopy(booleans []bool) bool {
	for _, b := range booleans {
		if !b {
			// false value not found
			return false
		}
	}
	// false value found
	return true
}
