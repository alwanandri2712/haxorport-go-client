package model

// AuthResponse adalah struktur untuk respons validasi token
type AuthResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    AuthData `json:"data"`
	Meta    AuthMeta `json:"meta"`
}

// AuthData adalah struktur untuk data pengguna dalam respons validasi token
type AuthData struct {
	UserID       string       `json:"user_id"`
	Fullname     string       `json:"fullname"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	Subscription Subscription `json:"subscription"`
}

// AuthMeta adalah struktur untuk metadata dalam respons validasi token
type AuthMeta struct {
	HeaderStatusCode int `json:"header_status_code"`
}

// Subscription adalah struktur untuk informasi langganan pengguna
type Subscription struct {
	Name     string            `json:"name"`
	Limits   SubscriptionLimits `json:"limits"`
	Features SubscriptionFeatures `json:"features"`
}

// SubscriptionLimits adalah struktur untuk batasan langganan
type SubscriptionLimits struct {
	Tunnels    ResourceLimit `json:"tunnels"`
	Ports      ResourceLimit `json:"ports"`
	Bandwidth  ResourceLimit `json:"bandwidth"`
	Requests   ResourceLimit `json:"requests"`
}

// ResourceLimit adalah struktur untuk batasan sumber daya
type ResourceLimit struct {
	Limit   int  `json:"limit"`
	Used    int  `json:"used"`
	Reached bool `json:"reached"`
}

// SubscriptionFeatures adalah struktur untuk fitur langganan
type SubscriptionFeatures struct {
	CustomDomains    bool `json:"customDomains"`
	Analytics        bool `json:"analytics"`
	PrioritySupport  bool `json:"prioritySupport"`
}
