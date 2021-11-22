package sac

import (
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/resty.v1"
	"net/http"
	"net/url"
	"time"
)

var ErrorNotFound = fmt.Errorf("permission denied")

type SecureAccessCloudClientImpl struct {
	Setting *SecureAccessCloudSettings
	Client  *resty.Client
}

func NewSecureAccessCloudClientImpl(setting *SecureAccessCloudSettings) SecureAccessCloudClient {
	return &SecureAccessCloudClientImpl{Client: nil, Setting: setting}
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Application API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) CreateApplication() error {
	return nil
}

func (s *SecureAccessCloudClientImpl) UpdateApplication() error {
	return nil
}

func (s *SecureAccessCloudClientImpl) FindApplicationByName(name string) (dto.Application, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications" + "?filter=" + url.QueryEscape(name)

	var applications dto.Applications
	err := s.performGetRequest(endpoint, &applications)

	if err != nil {
		return dto.Application{}, err
	}

	if applications.NumberOfElements == 0 {
		return dto.Application{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return applications.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) DeleteApplication(id string) error {
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Policy API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) FindPolicyByName(name string) (dto.Policy, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/policies" + "?filter=" + url.QueryEscape(name)

	var policies dto.Policies
	err := s.performGetRequest(endpoint, &policies)

	if err != nil {
		return dto.Policy{}, err
	}

	if policies.NumberOfElements == 0 {
		return dto.Policy{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return policies.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) AddApplicationToPolicy() error {
	return nil
}

func (s *SecureAccessCloudClientImpl) RemoveApplicationFromPolicy() error {
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Site API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) FindSiteByName(name string) (dto.Site, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/sites" + "?filter=" + url.QueryEscape(name)

	var sites dto.Sites
	err := s.performGetRequest(endpoint, &sites)

	if err != nil {
		return dto.Site{}, err
	}

	if sites.NumberOfElements == 0 {
		return dto.Site{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return sites.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) AddApplicationToSite() error {
	return nil
}

func (s *SecureAccessCloudClientImpl) RemoveApplicationFromSite() error {
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Connector API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private Functions
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) getClient() *resty.Client {
	if s.Client != nil {
		return s.Client
	}

	cfg := clientcredentials.Config{
		ClientID:     s.Setting.ClientID,
		ClientSecret: s.Setting.ClientSecret,
		TokenURL:     s.Setting.BuildOAuthTokenURL(),
		Scopes:       []string{},
	}

	oauthClient := cfg.Client(context.Background())
	client := resty.New().SetRetryCount(0).SetTimeout(1 * time.Minute)
	// http://godoc.org/golang.org/x/oauth2 implements `httpRoundTripper` interface
	// Set the oauthClient transport
	client.SetTransport(oauthClient.Transport)

	s.Client = client
	return s.Client
}

func (s *SecureAccessCloudClientImpl) performGetRequest(endpoint string, obj interface{}) error {
	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the GET request
	response, err := client.NewRequest().Get(endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusNotFound {
		return ErrorNotFound
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	// 3. Convert to Commit model
	err = json.Unmarshal(response.Body(), &obj)
	if err != nil {
		return err
	}

	return nil
}
