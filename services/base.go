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

func createYamlQuotes(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	// used to update a YAML pipeline if quotes are defined.  Will add the orgIdentifier and projectIdentifier if not found
	var out string

	if strings.Contains(yaml, "orgIdentifier: ") {
		out = strings.ReplaceAll(yaml, "orgIdentifier: \""+sourceOrg+"\"", "orgIdentifier: \""+targetOrg+"\"")
	} else {
		out = fmt.Sprintln(yaml, " orgIdentifier:", targetOrg)
	}

	if strings.Contains(yaml, "projectIdentifier: ") {
		out = strings.ReplaceAll(out, "projectIdentifier: \""+sourceProject, "projectIdentifier: \""+targetProject)
	} else {
		out = fmt.Sprintln(yaml, " projectIdentifier:", targetProject)
	}

	return out
}

func updateYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject string) string {
	// used to only update a YAML pipeline. Will NOT add the orgIdentifier and projectIdentifier if not found
	var out string
	// var orgIdOut string
	// var projIdOut string

	if strings.Contains(yaml, "\"orgIdentifier\": \""+sourceOrg+"\"") || strings.Contains(yaml, "orgIdentifier: "+sourceOrg) {
		orgIdOut := strings.ReplaceAll(yaml, "\"orgIdentifier\": \""+sourceOrg+"\"", "\"orgIdentifier\": \""+targetOrg+"\"")
		out = strings.ReplaceAll(orgIdOut, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	}

	if strings.Contains(yaml, "\"projectIdentifier\": \""+sourceProject+"\"") || strings.Contains(yaml, "projectIdentifier: "+sourceProject) {
		projIdOut := strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
		out = strings.ReplaceAll(projIdOut, "\"projectIdentifier\": \""+sourceProject+"\"", "\"projectIdentifier\": \""+targetProject+"\"")
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
