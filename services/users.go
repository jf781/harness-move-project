package services

import (
	"encoding/json"
	"fmt"

	"github.com/jf781/harness-move-project/model"
	"github.com/schollz/progressbar/v3"
)

const USERS = "/ng/api/user/aggregate"

type UserScopeContext struct {
	api           *ApiRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewUserScopeOperation(api *ApiRequest, sourceOrg, sourceProject, targetOrg, targetProject string) UserScopeContext {
	return UserScopeContext{
		api:           api,
		sourceOrg:     sourceOrg,
		sourceProject: sourceProject,
		targetOrg:     targetOrg,
		targetProject: targetProject,
	}
}

func (c UserScopeContext) Move() error {

	users, err := c.api.listUsers(c.sourceOrg, c.sourceProject)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(users)), "Roles")
	var failed []string

	for _, u := range users {

		emails := &model.UserEmail{
			EmailAddress:      u.Email,
			OrgIdentifier:     c.targetOrg,
			ProjectIdentifier: c.targetProject,
		}

		err = c.api.AddUserToScope(emails)

		if err != nil {
			failed = append(failed, fmt.Sprintln(u.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "Roles:")
	return nil
}

func (api *ApiRequest) listUsers(org, project string) ([]*model.User, error) {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     org,
			"projectIdentifier": project,
			"size":              "1000",
		}).
		Get(BaseURL + ROLE)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.GetUserResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	users := []*model.User{}
	for _, c := range result.Data.Content {
		users = append(users, &c.User)

	}

	return users, nil
}

func (api *ApiRequest) AddUserToScope(user *model.UserEmail) error {

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
			"orgIdentifier":     user.OrgIdentifier,
			"projectIdentifier": user.ProjectIdentifier,
		}).
		Post(BaseURL + ROLE)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
