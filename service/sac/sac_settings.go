package sac

type SecureAccessCloudSettings struct {
	ClientID     string
	ClientSecret string
	TenantDomain string
}

func (s *SecureAccessCloudSettings) BuildAPIPrefixURL() string {
	return "https://api." + s.TenantDomain
}

func (s *SecureAccessCloudSettings) BuildOAuthTokenURL() string {
	return s.BuildAPIPrefixURL() + "/v1/oauth/token"
}
