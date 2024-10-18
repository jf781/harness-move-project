package services

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"harness-copy-project/model"
	"time"
)

const FREEZEPROJECT = "/ng/api/freeze"
const FREEZEPROJECTSTATUS = "/ng/api/freeze/updateFreezeStatus"

type FreezeSourceProjectContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	logger        *zap.Logger
}

func FreezeSourceProjectOperation(api *ApiRequest, sourceOrg, sourceProject string, logger *zap.Logger) FreezeSourceProjectContext {
	return FreezeSourceProjectContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		logger:        logger,
	}
}

func (c FreezeSourceProjectContext) Copy() error {

	timeZone, _ := time.LoadLocation("America/Los_Angeles")
	currentTime := time.Now().In(timeZone)

	fmt.Printf("Freezing source project %s\n", c.sourceProject)

	c.logger.Info("Freezing project",
		zap.String("project", c.sourceProject),
	)

	freeze := model.FreezeRequest{
		Freeze: model.Freeze{
			Name:       "Harness Copy Project Freeze",
			Identifier: "hrns_copy_prj_freeze",
			EntityConfigs: []model.EntityConfig{
				{
					Name: "hrns_copy_prj_freeze",
					Entities: []model.Entity{
						{Type: "Service", FilterType: "All"},
						{Type: "EnvType", FilterType: "All"},
					},
				},
			},
			Status:            "Disabled",
			OrgIdentifier:     c.sourceOrg,
			ProjectIdentifier: c.sourceProject,
			Windows: []model.Window{
				{
					TimeZone:  "America/Los_Angeles",
					StartTime: currentTime.Format("2006-01-02 03:04 PM"),
					Duration:  "365d",
				},
			},
		},
	}

	freezeYaml, err := yaml.Marshal(&freeze)
	freezeRequest := string(freezeYaml)
	fmt.Println(freezeRequest)

	if err != nil {
		c.logger.Error("Failed to marshal freeze request",
			zap.String("project", c.sourceProject),
			zap.Error(err),
		)
	}

	freezeId, err := c.api.createProjectFreeze(&freezeRequest, c.sourceOrg, c.sourceProject, c.logger)

	if freezeId == nil {
		return fmt.Errorf("failed to create freeze for project %s", c.sourceProject)
	}
	freezeSlc := []string{*freezeId}

	if err != nil {
		c.logger.Error("Failed to create project freeze",
			zap.String("project", c.sourceOrg),
			zap.Error(err),
		)
	} else {
		err := c.api.enableProjectFreeze(&freezeSlc, c.sourceOrg, c.sourceProject, c.logger)

		if err != nil {
			c.logger.Error("Failed to enable project freeze ",
				zap.String("project", c.sourceOrg),
				zap.Error(err),
			)
		}
	}
	return nil
}

func (api *ApiRequest) createProjectFreeze(freeze *string, org string, project string, logger *zap.Logger) (*string, error) {

	api.Client.SetDebug(true)

	logger.Info("Creating project freeze",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/yaml").
		SetBody(*freeze).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
		}).
		Post(api.BaseURL + FREEZEPROJECT)
	if err != nil {
		logger.Error("Failed to request project freeze creation",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when creating project freeze",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.FreezeResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	if result.Status != "SUCCESS" {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return &result.Data.Identifier, nil
}

func (api *ApiRequest) enableProjectFreeze(freezeId *[]string, org string, project string, logger *zap.Logger) error {

	logger.Info("Enabling project freeze",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(freezeId).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"status":            "Enabled",
		}).
		Post(api.BaseURL + FREEZEPROJECTSTATUS)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Project", project),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate project freeze.",
					zap.String("project", project),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when setting project freeze status ",
				zap.String("project", project),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
