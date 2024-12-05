package services

import (
	"encoding/json"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"harness-copy-project/model"
)

const ENVGROUPLIST = "/ng/api/environmentGroup/list"
const ENVGROUP = "/ng/api/environmentGroup"

type EnvGroupContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
	showPB        bool
}

func NewEnvGroupOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger, showPB bool) EnvGroupContext {
	return EnvGroupContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
		showPB:        showPB,
	}
}

func (c EnvGroupContext) Copy() error {

	c.logger.Info("Copying environment groups for ",
		zap.String("project: ", c.sourceProject),
	)

	// Leveraging listPipelines func from pipeline.go file
	envGroups, err := c.api.listEnvGroups(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive environment groups",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	var bar *progressbar.ProgressBar

	if c.showPB {
		bar = progressbar.Default(int64(len(envGroups)), "Environment Groups:    ")
	}

	for _, eg := range envGroups {

		IncrementEnvironmentGroupsTotal()

		c.logger.Info("Processing environment group",
			zap.String("environment group", eg.EnvGroup.Name),
			zap.String("target project", c.targetProject),
		)

		e := model.CreateEnvGroup{}

		newYaml := updateYaml(eg.EnvGroup.YAML, c.targetOrg, c.targetProject)
		e.OrgIdentifier = c.targetOrg
		e.ProjectIdentifier = c.targetProject
		e.Color = eg.EnvGroup.Color
		e.Identifier = eg.EnvGroup.Identifier
		e.YAML = newYaml

		err = c.api.createEnvGroup(e, c.logger)

		if err != nil {
			c.logger.Error("Failed to create environment group",
				zap.String("environment group", eg.EnvGroup.Name),
				zap.Error(err),
			)
		} else {
			IncrementEnvironmentGroupsMoved()
		}
		if c.showPB {
			bar.Add(1)
		}
	}
	if c.showPB {
		bar.Finish()
	}

	return nil
}

func (api *ApiRequest) listEnvGroups(org, project string, logger *zap.Logger) ([]model.EnvGroupContent, error) {

	logger.Info("Fetching environment grups",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "100",
		}).
		Post(api.BaseURL + ENVGROUPLIST)
	if err != nil {
		logger.Error("Failed to request to list of environment groups",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing environment groups",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetEnvGroupResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) createEnvGroup(envGroup model.CreateEnvGroup, logger *zap.Logger) error {

	logger.Info("Creating environment group",
		zap.String("org", envGroup.OrgIdentifier),
		zap.String("project", envGroup.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(envGroup).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + ENVGROUP)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("environment group", envGroup.Identifier),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate environment group found, ignoring error",
					zap.String("environment group", envGroup.Identifier),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("environment group", envGroup.Identifier),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
