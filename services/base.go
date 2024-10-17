package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
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

func createYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	// used to update a YAML pipeline.  Will add the orgIdentifier and projectIdentifier if not found
	var out string

	if strings.Contains(yaml, "orgIdentifier: ") {
		out = strings.ReplaceAll(yaml, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	} else {
		out = fmt.Sprintln(yaml, " orgIdentifier:", targetOrg)
	}

	if strings.Contains(yaml, "projectIdentifier: ") {
		out = strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
	} else {
		out = fmt.Sprintln(yaml, " projectIdentifier:", targetProject)
	}

	return out
}

func updateYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	// Initialize out with the original YAML content
	out := yaml

	// Handle JSON-style with quoted key and quoted value for orgIdentifier
	if strings.Contains(yaml, "\"orgIdentifier\": \""+sourceOrg+"\"") {
		out = strings.ReplaceAll(out, "\"orgIdentifier\": \""+sourceOrg+"\"", "\"orgIdentifier\": \""+targetOrg+"\"")
	}

	// Handle JSON-style with quoted key and unquoted value for orgIdentifier
	if strings.Contains(yaml, "\"orgIdentifier\": "+sourceOrg) {
		out = strings.ReplaceAll(out, "\"orgIdentifier\": "+sourceOrg, "\"orgIdentifier\": "+targetOrg)
	}

	// Handle YAML-style with unquoted key and quoted value for orgIdentifier
	if strings.Contains(yaml, "orgIdentifier: \""+sourceOrg+"\"") {
		out = strings.ReplaceAll(out, "orgIdentifier: \""+sourceOrg+"\"", "orgIdentifier: \""+targetOrg+"\"")
	}

	// Handle YAML-style with unquoted key and unquoted value for orgIdentifier
	if strings.Contains(yaml, "orgIdentifier: "+sourceOrg) {
		out = strings.ReplaceAll(out, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	}

	// Handle JSON-style with quoted key and quoted value for projectIdentifier
	if strings.Contains(yaml, "\"projectIdentifier\": \""+sourceProject+"\"") {
		out = strings.ReplaceAll(out, "\"projectIdentifier\": \""+sourceProject+"\"", "\"projectIdentifier\": \""+targetProject+"\"")
	}

	// Handle JSON-style with quoted key and unquoted value for projectIdentifier
	if strings.Contains(yaml, "\"projectIdentifier\": "+sourceProject) {
		out = strings.ReplaceAll(out, "\"projectIdentifier\": "+sourceProject, "\"projectIdentifier\": "+targetProject)
	}

	// Handle YAML-style with unquoted key and quoted value for projectIdentifier
	if strings.Contains(yaml, "projectIdentifier: \""+sourceProject+"\"") {
		out = strings.ReplaceAll(out, "projectIdentifier: \""+sourceProject+"\"", "projectIdentifier: \""+targetProject+"\"")
	}

	// Handle YAML-style with unquoted key and unquoted value for projectIdentifier
	if strings.Contains(yaml, "projectIdentifier: "+sourceProject) {
		out = strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
	}

	return out
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
