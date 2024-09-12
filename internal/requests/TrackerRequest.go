package requests

type TrackerRequest struct {
    Exchange string `json:"exchange"`
    Currency string `json:"currency"`
    Side     string `json:"side"`
    Username string `json:"username"`
    IsAggregated bool `json:"is_aggregated"`
}

