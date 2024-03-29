package model

// Generated by https://quicktype.io

type StoreType string

const (
	Inline StoreType = "INLINE"
	Remote StoreType = "REMOTE"
)

type Pageable struct {
	Offset     int64 `json:"offset"`
	PageSize   int64 `json:"pageSize"`
	PageNumber int64 `json:"pageNumber"`
	Paged      bool  `json:"paged"`
	Unpaged    bool  `json:"unpaged"`
}

type ErrorResponse struct {
	Status        string `json:"status"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	CorrelationID string `json:"correlationId"`
}

type EntityValidityDetails struct {
	Valid       bool    `json:"valid"`
	InvalidYAML *string `json:"invalidYaml,omitempty"`
}
