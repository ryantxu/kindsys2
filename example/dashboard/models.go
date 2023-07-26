package dashboard

type DataSourceRef struct {
	Type string `json:"type,omitempty"` // eg prometheus | hubble
	UID  string `json:"uid,omitempty"`
}

type DataQuery struct {
	RefId      string        `json:"refId,omitempty"`
	Datasource DataSourceRef `json:"datasource,omitempty"`

	// more here depending on the type!!!
}

type DummyDashboard struct {
	Timezone string `json:"timezone,omitempty"`
	Editable bool   `json:"editable,omitempty"`

	// TODO, next version should be an array of verdicts
	Target any `json:"target,omitempty"`
}
