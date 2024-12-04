package services

import (
	"encoding/json"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"harness-copy-project/model"
)

const LISTUSER = "/ng/api/user/aggregate"
const ADDUSER = "/ng/api/user/users"
const CURRENTUSER = "/ng/api/user/currentUser"
const REMOVEUSER = "/ng/api/user"

type UserScopeContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
	logger        *zap.Logger
	showPB        bool
}

type RemoveUserScopeContext struct {
	api           *ApiRequest
	targetOrg     string
	targetProject string
	logger        *zap.Logger
}

func NewUserScopeOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string, logger *zap.Logger, showPB bool) UserScopeContext {
	return UserScopeContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
		showPB:        showPB,
	}
}

func RemoveCurrentUserOperation(api *ApiRequest, targetOrg, targetProject string, logger *zap.Logger) RemoveUserScopeContext {
	return RemoveUserScopeContext{
		api:           api,
		targetOrg:     targetOrg,
		targetProject: targetProject,
		logger:        logger,
	}
}

func (c UserScopeContext) Copy() error {

	c.logger.Info("Copying users",
		zap.String("project", c.sourceProject),
	)

	users, err := c.api.listUsers(c.sourceOrg, c.sourceProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive users",
			zap.String("Project", c.sourceProject),
			zap.Error(err),
		)
		return err
	}

	var bar *progressbar.ProgressBar

	if c.showPB {
		bar = progressbar.Default(int64(len(users)), "Users    ")
	}

	for _, u := range users {

		IncrementUsersTotal()

		c.logger.Info("Processing user",
			zap.String("user", u.Name),
			zap.String("sourceProject", c.sourceProject),
		)

		userToAdd := &model.UserEmail{
			EmailAddress:      []string{u.Email},
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
		}

		err = c.api.addUserToScope(userToAdd, c.logger)

		if err != nil {
			c.logger.Error("Failed to create user",
				zap.String("user", u.Name),
				zap.Error(err),
			)
		} else {
			IncrementUsersMoved()
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

func (c RemoveUserScopeContext) Copy() error {

	c.logger.Info("Removing current user assignment",
		zap.String("project", c.targetProject),
	)

	currentUser, err := c.api.getCurrentUserInfo(c.targetOrg, c.targetProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to retrive current user",
			zap.String("Project", c.targetProject),
			zap.Error(err),
		)
		return err
	}

	err = c.api.RemoveCurrentUserAccess(*currentUser, c.targetOrg, c.targetProject, c.logger)
	if err != nil {
		c.logger.Error("Failed to remove current user",
			zap.String("Project", c.targetProject),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (api *ApiRequest) listUsers(org, project string, logger *zap.Logger) ([]*model.User, error) {

	logger.Info("Fetching users",
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
		Post(api.BaseURL + LISTUSER)
	if err != nil {
		logger.Error("Failed to request to list of users",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when listing users",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetUserResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	users := []*model.User{}
	for _, c := range result.Data.Content {
		newUser := c.User
		users = append(users, &newUser)
	}

	return users, nil
}

func (api *ApiRequest) addUserToScope(user *model.UserEmail, logger *zap.Logger) error {

	logger.Info("Creating user",
		zap.String("user", user.EmailAddress[0]),
		zap.String("project", user.ProjectIdentifier),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     user.OrgIdentifier,
			"projectIdentifier": user.ProjectIdentifier,
		}).
		Post(api.BaseURL + ADDUSER)

	if err != nil {
		logger.Error("Failed to send request to create ",
			zap.String("user", user.EmailAddress[0]),
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorResponse); err == nil {
			if code, ok := errorResponse["code"].(string); ok && code == "DUPLICATE_FIELD" {
				// Log as a warning and skip the error
				logger.Info("Duplicate user found, ignoring error",
					zap.String("user", user.EmailAddress[0]),
				)
				return nil
			}
		} else {
			logger.Error(
				"Error response from API when creating ",
				zap.String("user", user.EmailAddress[0]),
				zap.String("response",
					resp.String(),
				),
			)
		}
		return handleErrorResponse(resp)
	}

	return nil
}

func (api *ApiRequest) getCurrentUserInfo(org, project string, logger *zap.Logger) (*model.User, error) {

	logger.Info("Fetching current user info",
		zap.String("org", org),
		zap.String("project", project),
	)

	IncrementApiCalls()

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Get(api.BaseURL + CURRENTUSER)
	if err != nil {
		logger.Error("Failed to request to get current user",
			zap.Error(err),
		)
		return nil, err
	}
	if resp.IsError() {
		logger.Error("Error response from API when getting current user",
			zap.String("response",
				resp.String(),
			),
		)
		return nil, handleErrorResponse(resp)
	}

	result := model.GetCurrentUserResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return nil, err
	}

	currentUser := model.User{}
	currentUser = result.Data
	return &currentUser, nil
}

func (api *ApiRequest) RemoveCurrentUserAccess(user model.User, org, project string, logger *zap.Logger) error {

	logger.Info("Removing user access",
		zap.String("user", user.Name),
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
		}).
		Delete(api.BaseURL + REMOVEUSER + "/" + user.UUID)
	if err != nil {
		logger.Error("Failed to remove user from project",
			zap.Error(err),
		)
		return err
	}
	if resp.IsError() {
		logger.Error("Error response from API when removing user from project",
			zap.String("response",
				resp.String(),
			),
		)
		return handleErrorResponse(resp)
	}

	result := model.RemoveUserResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logger.Error("Failed to parse response from API",
			zap.Error(err),
		)
		return err
	}

	return nil
}
