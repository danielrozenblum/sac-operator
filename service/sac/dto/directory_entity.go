package dto

type DirectoryEntityDTO struct {
	DisplayName          string `json:"displayName"`
	IdentifierInProvider string `json:"identifierInProvider"`
	IdentityProviderID   string `json:"identityProviderID"`
	IdentityProviderType string `json:"identityProviderType"`
	Type                 string `json:"type"`
}
