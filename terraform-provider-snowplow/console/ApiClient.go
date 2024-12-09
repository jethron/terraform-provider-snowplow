package console

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/snowplow-devops/terraform-provider-snowplow/terraform-provider-snowplow/console/apitypes"
)

type ApiClient struct {
	http    *http.Client
	jwt     string
	baseUrl string
	OrgId   string
	version string
}

type ApiClientProvider interface {
	GetApiClient() *ApiClient
}

type tokenResponse struct {
	AccessToken string
}

type ErrorResponse struct {
	Message string
	TraceId string
}

func NewApiClient(ctx context.Context, version string, host string, apiKeyId string, apiKey string, orgId string) (*ApiClient, error) {
	client := ApiClient{
		baseUrl: fmt.Sprintf("%s/api/msc/v1", host),
		http:    http.DefaultClient,
		OrgId:   orgId,
		version: version,
	}

	if err := client.authenticate(ctx, apiKeyId, apiKey); err != nil {
		return nil, err
	}

	return &client, nil
}

func (a *ApiClient) newRequest(ctx context.Context, endpoint, method string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("https://%s%s", a.baseUrl, endpoint), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-SNOWPLOW-TERRAFORM", a.version)

	if a.jwt != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.jwt))
	}

	return req, err
}

func (a *ApiClient) newOrgRequest(ctx context.Context, endpoint, method string) (*http.Request, error) {
	if a.OrgId == "" {
		return nil, fmt.Errorf("can not make organization specific api request without organization id")
	}
	return a.newRequest(ctx, fmt.Sprintf("/organizations/%s%s", a.OrgId, endpoint), method)
}

func (a *ApiClient) readBody(req *http.Request) ([]byte, error) {
	resp, err := a.http.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v:%v\n%s\n", resp.StatusCode, req.URL, body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if len(body) > 0 {
			var errorResponse ErrorResponse
			err = json.Unmarshal(body, &errorResponse)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("console api failure: %v", errorResponse)
		} else {
			return nil, fmt.Errorf("console api value: unknown error for %v response", resp.StatusCode)
		}
	} else {
		if json.Valid(body) {
			return body, nil
		} else {
			return nil, fmt.Errorf("console returned invalid json content: %v", body)
		}
	}

}

func (a *ApiClient) authenticate(ctx context.Context, apiKeyId, apiKeySecret string) error {
	var req *http.Request
	var err error

	if apiKeySecret == "" {
		return fmt.Errorf("console api key required to authenticate")
	}

	if apiKeyId == "" {
		req, err = a.newOrgRequest(ctx, "/credentials/v2/token", http.MethodGet)

		if err != nil {
			return fmt.Errorf("error requesting v2 token; do you need a key ID to use v3?: %w", err)
		}

		req.Header.Add("X-API-KEY", apiKeySecret)
	} else {
		req, err = a.newOrgRequest(ctx, "/credentials/v3/token", http.MethodGet)

		if err != nil {
			return err
		}

		req.Header.Add("X-API-KEY-ID", apiKeyId)
		req.Header.Add("X-API-KEY", apiKeySecret)
	}

	body, err := a.readBody(req)

	if err != nil {
		return err
	}

	var tokenResponse tokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return err
	} else if tokenResponse.AccessToken == "" {
		return fmt.Errorf("unable to obtain console api token: %v", body)
	}

	a.jwt = tokenResponse.AccessToken

	return nil
}

func (a *ApiClient) GetOrganizations(ctx context.Context) ([]apitypes.Organization, error) {
	req, err := a.newRequest(ctx, "/organizations", http.MethodGet)

	if err != nil {
		return nil, err
	}

	body, err := a.readBody(req)

	if err != nil {
		return nil, err
	}

	organizations := make([]apitypes.Organization, 0)

	err = json.Unmarshal(body, &organizations)

	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (a *ApiClient) GetUser(ctx context.Context, id string) (*apitypes.User, error) {
	req, err := a.newOrgRequest(ctx, fmt.Sprintf("/users/%s", id), http.MethodGet)

	if err != nil {
		return nil, err
	}

	body, err := a.readBody(req)

	if err != nil {
		return nil, err
	}

	var user apitypes.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *ApiClient) GetUsers(ctx context.Context) ([]apitypes.User, error) {
	req, err := a.newOrgRequest(ctx, "/users", http.MethodGet)

	if err != nil {
		return nil, err
	}

	body, err := a.readBody(req)

	if err != nil {
		return nil, err
	}

	users := make([]apitypes.User, 0)

	err = json.Unmarshal(body, &users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a *ApiClient) GetPipeline(ctx context.Context, id string) (*apitypes.Pipeline, error) {
	req, err := a.newOrgRequest(ctx, fmt.Sprintf("/resources/v1/pipelines/%s", id), http.MethodGet)

	if err != nil {
		return nil, err
	}

	body, err := a.readBody(req)

	if err != nil {
		return nil, err
	}

	var pipeline apitypes.Pipeline

	err = json.Unmarshal(body, &pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (a *ApiClient) GetPipelines(ctx context.Context) ([]apitypes.Pipeline, error) {
	req, err := a.newOrgRequest(ctx, "/resources/v1/pipelines", http.MethodGet)

	if err != nil {
		return nil, err
	}

	body, err := a.readBody(req)

	if err != nil {
		return nil, err
	}

	pipelines := make([]apitypes.Pipeline, 0)

	err = json.Unmarshal(body, &pipelines)

	if err != nil {
		return nil, err
	}

	return pipelines, nil
}
