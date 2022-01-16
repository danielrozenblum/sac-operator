package sac

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"bitbucket.org/accezz-io/sac-operator/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (s *SecureAccessCloudClientImpl) CreateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications/"

	var createdApplicationDTO dto.ApplicationDTO

	err := s.performModifyRequest(http.MethodPost, endpoint, applicationDTO, createdApplicationDTO)
	if err != nil {
		return nil, err
	}

	return &createdApplicationDTO, nil
}

func (s *SecureAccessCloudClientImpl) UpdateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications/" + applicationDTO.ID.String()

	var createdApplicationDTO dto.ApplicationDTO

	err := s.performModifyRequest(http.MethodPut, endpoint, applicationDTO, createdApplicationDTO)
	if err != nil {
		return nil, err
	}

	return &createdApplicationDTO, nil
}

func (s *SecureAccessCloudClientImpl) FindApplicationByID(id uuid.UUID) (*dto.ApplicationDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications/" + id.String()

	var applications dto.ApplicationPageDTO
	err := s.performGetRequest(endpoint, &applications)

	if err != nil {
		return &dto.ApplicationDTO{}, err
	}

	if applications.NumberOfElements == 0 {
		return &dto.ApplicationDTO{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return &applications.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) FindApplicationByName(name string) (*dto.ApplicationDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications" + "?filter=" + url.QueryEscape(name)

	var applications dto.ApplicationPageDTO
	err := s.performGetRequest(endpoint, &applications)

	if err != nil {
		return &dto.ApplicationDTO{}, err
	}

	if applications.NumberOfElements == 0 {
		return &dto.ApplicationDTO{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return &applications.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) DeleteApplication(id uuid.UUID) error {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications/" + id.String()

	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the GET request
	response, err := client.NewRequest().Delete(endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK && response.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Policy API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) FindPolicyByName(name string) (dto.PolicyDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/policies" + "?filter=" + url.QueryEscape(name)

	var policies dto.PoliciesPageDTO
	err := s.performGetRequest(endpoint, &policies)

	if err != nil {
		return dto.PolicyDTO{}, err
	}

	if policies.NumberOfElements == 0 {
		return dto.PolicyDTO{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return policies.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) FindPoliciesByNames(names []string) ([]dto.PolicyDTO, error) {
	var results []dto.PolicyDTO

	for _, name := range names {
		policyDTO, err := s.FindPolicyByName(name)
		if err != nil {
			return results, err
		}

		results = append(results, policyDTO)
	}

	return results, nil
}

func (s *SecureAccessCloudClientImpl) UpdatePolicies(applicationId uuid.UUID, applicationType model.ApplicationType, policies []uuid.UUID) error {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/policies/by-app-id/" + applicationId.String()

	applicationToPoliciesBindingRequest := applicationToPoliciesBinding{
		ApplicationType: applicationType,
		PolicyIDs:       utils.ToStringArray(policies),
	}

	err := s.performModifyRequest(http.MethodPut, endpoint, applicationToPoliciesBindingRequest, nil)
	if err != nil {
		return err
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Site API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) FindSiteByName(name string) (*dto.SiteDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/sites" + "?filter=" + url.QueryEscape(name)

	var sites dto.SitePageDTO
	err := s.performGetRequest(endpoint, &sites)

	if err != nil {
		return &dto.SiteDTO{}, err
	}

	if sites.NumberOfElements == 0 {
		return &dto.SiteDTO{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return &sites.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) BindApplicationToSite(applicationId uuid.UUID, siteId uuid.UUID) error {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/applications/" + applicationId.String() + "/site-binding/" + siteId.String()
	return s.performModifyRequest(http.MethodPut, endpoint, nil, nil)
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

func (s *SecureAccessCloudClientImpl) performModifyRequest(method string, endpoint string, requestObj interface{}, responseObj interface{}) error {
	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the request
	request := client.NewRequest()
	var response *resty.Response
	var err error

	// 2.1. Marshal the request
	if requestObj != nil {
		body, err := json.Marshal(requestObj)
		if err != nil {
			return err
		}

		request = request.SetBody(body).SetHeader("Content-Type", "application/json")
	}

	switch method {
	case http.MethodPost:
		response, err = request.Post(endpoint)
	case http.MethodPut:
		response, err = request.Put(endpoint)
	default:
		return errors.New("unsupported http method: " + method)
	}

	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusCreated && response.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	// 4. Unmarshal response body
	if responseObj != nil {
		err = json.Unmarshal(response.Body(), &responseObj)
		if err != nil {
			return err
		}
	}

	return nil
}

type applicationToPoliciesBinding struct {
	ApplicationType model.ApplicationType
	PolicyIDs       []string
}
