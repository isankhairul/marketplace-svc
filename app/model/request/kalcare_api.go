package request

type OrderCancel struct {
	OrderNo               string `json:"order_no"`
	OrderStatus           string `json:"order_status"`
	OrderCreated          string `json:"order_created"`
	StatusExpired         string `json:"status_expired"`
	PaymentExpiredHours   int    `json:"payment_expired_hours"`
	PaymentExpiredMinutes int    `json:"payment_expired_minutes"`
	QueueType             string `json:"queue_type,omitempty"`
	NotifyUrl             string `json:"notify_url,omitempty"`
	NotifyHeader          string `json:"notify_header,omitempty"`
	NotifyCode            string `json:"notify_code,omitempty"`
	NotifyPayload         string `json:"notify_payload,omitempty"`
	MerchantId            string `json:"merchant_id,omitempty"`
}

type RegisteredWebhook struct {
	ID                int    `json:"id"`
	Url               string `json:"url"`
	Header            string `json:"header"`
	Name              string `json:"webhook_name"`
	Code              string `json:"webhook_code"`
	ThirdParty        string `json:"third_party_name"`
	ModuleIntegration string `json:"module_integration_slug"`
}

type AuthorizationData struct {
	Data []string `json:"records"`
}

type OrderCancelKafka struct {
	Body       string   `json:"body"`
	Properties []string `json:"properties"`
	Headers    []string `json:"headers"`
}

type OrderCancelArr []OrderCancel
type RegisteredWebhookArr []RegisteredWebhook
