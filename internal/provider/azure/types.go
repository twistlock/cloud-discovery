package azure

// Options are options for getting cloud data for Azure services
type Options struct {
	TenantID       string `json:"tenantId"`
	ClientID       string `json:"clientId"`
	Secret         string `json:"clientSecret"`
	SubscriptionID string `json:"subscriptionId"`
}
