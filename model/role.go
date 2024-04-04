package model

// Generated by https://quicktype.io

type GetRoleResponse struct {
	Status        string      `json:"status"`
	Data          GetRoleData `json:"data"`
	CorrelationID string      `json:"correlationId"`
}

type GetRoleData struct {
	TotalPages    int64              `json:"totalPages"`
	TotalItems    int64              `json:"totalItems"`
	PageItemCount int64              `json:"pageItemCount"`
	PageSize      int64              `json:"pageSize"`
	Content       []*RoleListContent `json:"content"`
	PageIndex     int64              `json:"pageIndex"`
	Empty         bool               `json:"empty"`
	PageToken     interface{}        `json:"pageToken"`
}

type RoleListContent struct {
	RoleAssignment RoleAssignmentContent `json:"roleAssignment"`
	Scope          RoleAssignmentScope   `json:"scope"`
	LastModifiedAt int64                 `json:"lastModifiedAt"`
	HarnessManaged bool                  `json:"harnessManaged"`
}

type RoleAssignmentContent struct {
	Identifier              string    `json:"identifier"`
	ResourceGroupIdentifier string    `json:"resourceGroupIdentifier"`
	RoleIdentifier          string    `json:"roleIdentifier"`
	Principal               Principal `json:"principal"`
	Disabled                bool      `json:"disabled"`
	Managed                 bool      `json:"managed"`
	Internal                bool      `json:"internal"`
}

type Principal struct {
	ScopeLevel *string `json:"scopeLevel"`
	Identifier string  `json:"identifier"`
	Type       string  `json:"type"`
}

type RoleAssignmentScope struct {
	AccountIdentifier string `json:"accountIdentifier"`
	OrgIdentifier     string `json:"orgIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}
