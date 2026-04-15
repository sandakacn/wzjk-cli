package api

import "time"

// Domain represents a domain subscription
type Domain struct {
	ID            string     `json:"id"`
	DomainID      string     `json:"domainId"`
	Domain        string     `json:"domain"`
	Port          int        `json:"port"`
	CheckType     string     `json:"checkType"`
	AlertDays     int        `json:"alertDays"`
	IsActive      bool       `json:"isActive"`
	SSLIssuer     *string    `json:"sslIssuer"`
	SSLValidFrom  *time.Time `json:"sslValidFrom"`
	SSLValidTo    *time.Time `json:"sslValidTo"`
	SSLSubject    *string    `json:"sslSubject"`
	LastCheckedAt *time.Time `json:"lastCheckedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// AddDomainRequest represents the request to add a domain
type AddDomainRequest struct {
	Domain    string `json:"domain"`
	Port      int    `json:"port,omitempty"`
	CheckType string `json:"checkType,omitempty"`
	AlertDays int    `json:"alertDays,omitempty"`
}

// UpdateDomainRequest represents the request to update a domain
type UpdateDomainRequest struct {
	AlertDays int  `json:"alertDays,omitempty"`
	IsActive  *bool `json:"isActive,omitempty"`
}

// SSLInfo represents SSL certificate information
type SSLInfo struct {
	Domain          string   `json:"domain"`
	Hostname        string   `json:"hostname"`
	Port            int      `json:"port"`
	Scheme          string   `json:"scheme"`
	Path            string   `json:"path"`
	SuggestedCheckType string `json:"suggestedCheckType"`
	Issuer          string   `json:"issuer"`
	Subject         string   `json:"subject"`
	ValidFrom       string   `json:"validFrom"`
	ValidTo         string   `json:"validTo"`
	DaysUntilExpiry int      `json:"daysUntilExpiry"`
	IsValid         bool     `json:"isValid"`
	DomainMismatch  bool     `json:"domainMismatch"`
	SubjectAltNames []string `json:"subjectAltNames"`
}

// UserProfile represents the user profile
type UserProfile struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Image              string `json:"image"`
	HasWechat          bool   `json:"hasWechat"`
	HasWechatService   bool   `json:"hasWechatService"`
	IsPro              bool   `json:"isPro"`
	LoginProvider      string `json:"loginProvider"`
	NotificationSettings struct {
		Email  bool `json:"email"`
		Wechat bool `json:"wechat"`
	} `json:"notificationSettings"`
	AlertPreferences struct {
		ExpiryAlert bool `json:"expiryAlert"`
	} `json:"alertPreferences"`
}

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// DomainListResponse is the response from the domains list API
type DomainListResponse struct {
	Success bool     `json:"success"`
	Data    []Domain `json:"data"`
	Error   string   `json:"error,omitempty"`
}

// DomainResponse is the response for single domain operations
type DomainResponse struct {
	Success bool   `json:"success"`
	Data    Domain `json:"data"`
	Error   string `json:"error,omitempty"`
}

// SSLCheckResponse is the response from SSL check API
type SSLCheckResponse struct {
	Success bool    `json:"success"`
	Data    SSLInfo `json:"data"`
	Error   string  `json:"error,omitempty"`
}

// UserProfileResponse is the response from user profile API
type UserProfileResponse struct {
	Success bool        `json:"success"`
	Data    UserProfile `json:"data"`
	Error   string      `json:"error,omitempty"`
}
