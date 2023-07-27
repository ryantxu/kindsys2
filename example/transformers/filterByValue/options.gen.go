package filterbyvalue

type OptionsGen struct {
	Filters []FilterByValueFilter `json:"filters"`
	Match   FilterByValueMatch    `json:"match"`
	Type    FilterByValueType     `json:"type"`
}

type FilterByValueFilter struct {
	Config    map[string]interface{} `json:"config"`
	FieldName string                 `json:"fieldName"`
}

type FilterByValueMatch string

const (
	All FilterByValueMatch = "all"
	Any FilterByValueMatch = "any"
)

type FilterByValueType string

const (
	Exclude FilterByValueType = "exclude"
	Include FilterByValueType = "include"
)
