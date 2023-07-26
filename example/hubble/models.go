package hubble

type HubbleDataQuery struct {
	WorkloadName string   `json:"workloadName,omitempty"`
	TraceId      string   `json:"traceId,omitempty"`
	Namespace    []string `json:"namespace,omitempty"`

	// TODO, next version should be an array of verdicts
	Verdict int64 `json:"verdict,omitempty"`
}
