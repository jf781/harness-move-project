package services

import (
	"encoding/json"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"harness-copy-project/model"
)

const CONNECTORLOOKUP = "/ng/api/connectors/listV2"
const CONNECTORCREATE = "/ng/api/connectors"

type ConnectorContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
	showPB        bool
}

func NewConnectorOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger, showPB bool) ConnectorContext {
	return ConnectorContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
		showPB:        showPB,
	}
}

func (c ConnectorContext) Copy() error {

	c.logger.Info("Copying Connectors",
		zap.String("project", c.sourceProject),
	)

	connectors, err := c.api.listConnectors(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive connectors",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	var bar *progressbar.ProgressBar

	if c.showPB {
		bar = progressbar.Default(int64(len(connectors)), "Connectors    ")
	}

	for _, cn := range connectors {

		IncrementConnectorsTotal()

		c.logger.Info("Processing connector",
			zap.String("connector", cn.Connector.Name),
			zap.String("targetProject", c.targetProject),
		)

		cn.Connector.OrgIdentifier = c.targetOrg
		cn.Connector.ProjectIdentifier = c.targetProject

		err = c.api.addConnector(cn, c.logger)

		if err != nil {
			c.logger.Error("Failed to create connector",
				zap.String("connector", cn.Connector.Name),
				zap.Error(err),
			)
		} else {
			IncrementConnectorsMoved()
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

func (api *ApiRequest) listConnectors(org, project string, logger *zap.Logger) ([]*model.ConnectorContent, error) {

	logger.Info("Fetching connectors",
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
			"pageSize":          "100",
		}).
		Post(api.BaseURL + CONNECTORLOOKUP)
	if err != nil {
		logger.Error("Failed to request to list of connectors",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing connectors",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.ConnectorListResult{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	connectors := []*model.ConnectorContent{}
	for _, c := range result.Data.Content {
		if !c.HarnessManaged {
			newConnectors := c
			connectors = append(connectors, &newConnectors)
		} else {
			logger.Warn("Skipping connector because it is managed by Harness.",
				zap.String("connector", c.Connector.Name),
				zap.String("status", c.Status.Status),
				zap.Bool("harnessManaged", c.HarnessManaged),
			)
		}
	}

	return connectors, nil
}

func (api *ApiRequest) addConnector(connector *model.ConnectorContent, logger *zap.Logger) error {

	logger.Info("Creating connector",
		zap.String("connector", connector.Connector.Name),
		zap.String("project", connector.Connector.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(connector).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.BaseURL + CONNECTORCREATE)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Connector", connector.Connector.Name),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate connector found, ignoring error",
					zap.String("connector", connector.Connector.Name),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("Connector", connector.Connector.Name),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
