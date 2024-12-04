package model

type FreezeRequest struct {
	Freeze Freeze `yaml:"freeze"`
}

type Freeze struct {
	Name              string         `yaml:"name"`
	Identifier        string         `yaml:"identifier"`
	EntityConfigs     []EntityConfig `yaml:"entityConfigs"`
	Status            string         `yaml:"status"`
	OrgIdentifier     string         `yaml:"orgIdentifier"`
	ProjectIdentifier string         `yaml:"projectIdentifier"`
	Windows           []Window       `yaml:"windows"`
}

type EntityConfig struct {
	Name     string   `yaml:"name"`
	Entities []Entity `yaml:"entities"`
}

type Entity struct {
	Type       string `yaml:"type"`
	FilterType string `yaml:"filterType"`
}

type Window struct {
	TimeZone  string `yaml:"timeZone"`
	StartTime string `yaml:"startTime"`
	Duration  string `yaml:"duration"`
}

type FreezeResponse struct {
	Status           string                  `json:"status"`
	Code             string                  `json:"code,omitempty"`
	Message          string                  `json:"message,omitempty"`
	CorrelationId    string                  `json:"correlationId"`
	DetailedMessage  *string                 `json:"detailedMessage,omitempty"`
	ResponseMessages []FreezeResponseMessage `json:"responseMessages,omitempty"`
	Metadata         *interface{}            `json:"metadata,omitempty"`
	Data             *FreezeResponseData     `json:"data,omitempty"` // Only present in the "SUCCESS" response
	MetaData         *interface{}            `json:"metaData,omitempty"`
}

type FreezeResponseData struct {
	AccountId         string                 `json:"accountId"`
	Type              string                 `json:"type"`
	Status            string                 `json:"status"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Tags              map[string]interface{} `json:"tags"`
	OrgIdentifier     string                 `json:"orgIdentifier"`
	ProjectIdentifier string                 `json:"projectIdentifier"`
	Identifier        string                 `json:"identifier"`
	Yaml              string                 `json:"yaml"`
	CreatedAt         int64                  `json:"createdAt"`
	LastUpdatedAt     int64                  `json:"lastUpdatedAt"`
	FreezeScope       string                 `json:"freezeScope"`
}

// Struct for error response messages
type FreezeResponseMessage struct {
	Code           string                 `json:"code"`
	Level          string                 `json:"level"`
	Message        string                 `json:"message"`
	Exception      *string                `json:"exception,omitempty"`
	FailureTypes   []string               `json:"failureTypes"`
	AdditionalInfo map[string]interface{} `json:"additionalInfo"`
}
