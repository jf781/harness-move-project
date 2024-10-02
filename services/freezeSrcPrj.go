package services

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"harness-copy-project/model"
)

// const ROLEASSIGNMENT = "/authz/api/roleassignments" - Pulled from harness-move-project/services/roleAssignment.go

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

	fmt.Printf("Freezing source project %s\n", c.sourceProject)

	c.logger.Info("Copying role assignments",
		zap.String("project", c.sourceProject),
	)

	roleAssignments, err := c.api.listActiveRoleAssignments(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive role assignments",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	for _, r := range roleAssignments {

		// Add step where it removes the source role assignment
		// Call the createRolEAssingment function to create a new role assignment with the `_project_viewer` role

		c.logger.Info("Processing role assignment",
			zap.String("role assignment", r.RoleIdentifier),
			zap.String("sourceProject", c.sourceProject),
		)

		oldRole := &model.NewRoleAssignment{
			Identifier:              r.Identifier,
			ResourceGroupIdentifier: r.ResourceGroupIdentifier,
			RoleIdentifier:          r.RoleIdentifier,
			Principal:               r.Principal,
			OrgIdentifier:           c.sourceOrg,
			ProjectIdentifier:       c.sourceProject,
		}

		err = c.api.removeRoleAssignment(oldRole, c.logger)

		if err != nil {
			c.logger.Error("Failed to create role assignment",
				zap.String("role assignment", r.Identifier),
				zap.Error(err),
			)
		} else {
			newRole := &model.NewRoleAssignment{
				Identifier:              r.Identifier,
				ResourceGroupIdentifier: r.ResourceGroupIdentifier,
				RoleIdentifier:          "_project_viewer",
				Principal:               r.Principal,
				OrgIdentifier:           c.sourceOrg,
				ProjectIdentifier:       c.sourceProject,
			}

			err = c.api.createRoleAssignment(newRole, c.logger)

			if err != nil {
				c.logger.Error("Failed to create _project_viewer role assignment",
					zap.String("role assignment", r.Identifier),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

func (api *ApiRequest) listActiveRoleAssignments(org, project string, logger *zap.Logger) ([]*model.ExistingRoleAssignment, error) {

	logger.Info("Fetching role assignments",
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
		Get(api.BaseURL + ROLEASSIGNMENT)
	if err != nil {
		logger.Error("Failed to request to list of role assignments",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing role assignments",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetRoleAssignmentResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	roleAssignments := []*model.ExistingRoleAssignment{}
	for _, c := range result.Data.Content {
		if !c.RoleAssignment.Disabled || !c.RoleAssignment.Managed {
			newRoleAssignment := c.RoleAssignment
			roleAssignments = append(roleAssignments, &newRoleAssignment)
		}
	}

	return roleAssignments, nil
}

func (api *ApiRequest) removeRoleAssignment(role *model.NewRoleAssignment, logger *zap.Logger) error {

	logger.Info("Deleting the role assignment",
		zap.String("role assignment", role.RoleIdentifier),
		zap.String("project", role.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(role).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     role.OrgIdentifier,
			"projectIdentifier": role.ProjectIdentifier,
		}).
		Delete(api.BaseURL + ROLEASSIGNMENT + "/" + role.Identifier)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("Role assignment", role.Identifier),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate role assignment found, ignoring error",
					zap.String("role assignment", role.Identifier),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when removing  ",
				zap.String("Role assignment", role.Identifier),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}
