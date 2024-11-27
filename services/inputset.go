package services

import (
	"encoding/json"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"harness-copy-project/model"
)

type InputsetContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
	showPB        bool
}

func NewInputsetOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger, showPB bool) InputsetContext {
	return InputsetContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
		showPB:        showPB,
	}
}

func (c InputsetContext) Copy() error {
	c.logger.Info("Copying input sets",
		zap.String("project", c.sourceProject),
	)

	pipelines, err := c.api.listPipelines(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive pipelines",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	var bar *progressbar.ProgressBar

	if c.showPB {
		bar = progressbar.Default(int64(len(pipelines)), "Inputsets")
	}

	for _, pipeline := range pipelines {
		inputsets, err := c.api.listInputsets(c.sourceOrg, c.sourceProject, pipeline.Identifier, c.logger)
		if err != nil {
			c.logger.Error("Failed to retrive inputsets",
				zap.String("Project", c.sourceProject),
				zap.Error(err),
			)
			continue
		}

		if c.showPB {
			bar.ChangeMax(bar.GetMax() + len(inputsets))
		}

		for _, inputset := range inputsets {

			IncrementInputSetsTotal()

			c.logger.Info("Processing Inputset",
				zap.String("inputset", inputset.Name),
				zap.String("targetProject", c.targetProject),
				zap.String("pipeline", pipeline.Name),
			)
			is, err := c.api.getInputset(c.sourceOrg, c.sourceProject, pipeline.Identifier, inputset.Identifier, c.logger)
			if err == nil {
				newYaml := updateYaml(is.Yaml, c.targetOrg, c.targetProject)
				err = c.api.createInputset(c.targetOrg, c.targetProject, pipeline.Identifier, newYaml, c.logger)
			}
			if err != nil {
				c.logger.Error("Failed to create input set",
					zap.String("input set", inputset.Name),
					zap.Error(err),
				)
			} else {
				IncrementInputSetsMoved()
			}
			if c.showPB {
				bar.Add(1)
			}
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

func (api *ApiRequest) listInputsets(org, project, pipelineIdentifier string, logger *zap.Logger) ([]*model.ListInputsetContent, error) {

	logger.Info("Fetching inputsets",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifier":      org,
			"projectIdentifier":  project,
			"pipelineIdentifier": pipelineIdentifier,
			"inputSetType":       "ALL",
			"size":               "1000",
		}).
		Get(api.BaseURL + "/pipeline/api/inputSets")
	if err != nil {
		logger.Error("Failed to request to list of inputsets",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing input sets",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := &model.ListInputsetResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data.Content, nil
}

func (api *ApiRequest) getInputset(org, project, pipelineIdentifier, isIdentifier string, logger *zap.Logger) (*model.GetInputsetData, error) {

	logger.Info("Getting input set details",
		zap.String("input set", isIdentifier),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Load-From-Cache", "false").
		SetPathParam("inputset", isIdentifier).
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifier":      org,
			"projectIdentifier":  project,
			"pipelineIdentifier": pipelineIdentifier,
		}).
		Get(api.BaseURL + "/pipeline/api/inputSets/{inputset}")
	if err != nil {
		logger.Error("Failed to send request to get input set details ",
			zap.String("input set", isIdentifier),
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when fechting input set: "+isIdentifier,
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := &model.GetInputsetResponse{}
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	return result.Data, nil
}

func (api *ApiRequest) createInputset(org, project, pipelineIdentifier, yaml string, logger *zap.Logger) error {
	logger.Info("Creating inputset on pipeline",
		zap.String("pipeline", pipelineIdentifier),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/yaml").
		SetBody(yaml).
		SetQueryParams(map[string]string{
			"accountIdentifier":  api.Account,
			"orgIdentifier":      org,
			"projectIdentifier":  project,
			"pipelineIdentifier": pipelineIdentifier,
		}).
		Post(api.BaseURL + "/pipeline/api/inputSets")
	if err != nil {
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate pipeline input set found, ignoring error",
					zap.String("pipeline", pipelineIdentifier),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("pipeline", pipelineIdentifier),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
