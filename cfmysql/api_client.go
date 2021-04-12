package cfmysql

import (
	"code.cloudfoundry.org/cli/plugin"
	sdkModels "code.cloudfoundry.org/cli/plugin/models"
	"encoding/json"
	"fmt"
	pluginModels "github.com/elliott-neal/cf-mysql-plugin/cfmysql/models"
	"github.com/elliott-neal/cf-mysql-plugin/cfmysql/resources"
	"io"
	"net/url"
)

//go:generate counterfeiter . ApiClient
type ApiClient interface {
	GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error)
	GetService(cliConnection plugin.CliConnection, spaceGuid string, name string) (pluginModels.MysqlCredentials, error)
}

func NewApiClient(httpClient HttpWrapper) *apiClient {
	return &apiClient{
		httpClient: httpClient,
	}
}

type apiClient struct {
	httpClient HttpWrapper
	cliConfig  *CliConfig
	logWriter   io.Writer
}

func (self *apiClient) GetService(cliConnection plugin.CliConnection, spaceGuid string, name string) (pluginModels.MysqlCredentials, error) {
	path := fmt.Sprintf(
		"/v2/spaces/%s/service_instances?return_user_provided_service_instances=true&q=name%%3A%s",
		spaceGuid,
		url.QueryEscape(name),
	)
	// Grab service instances by named query.
	instancesResponse, err := self.getFromCfApi(path, cliConnection)
	if err != nil {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("error retrieving service instance: %s", err)
	}

	// Deserialize instances.
	_, instances, err := deserializeInstances(instancesResponse)
	if err != nil {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("error deserializing service instances: %s", err)
	}

	// Ensure query is not none.
	if len(instances) == 0 {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("%s not found in current space", name)
	}

	// Grab service instance from service_url.
	instanceResponse, err := self.getFromCfApi(instances[0].ServiceUrl, cliConnection)
	if err != nil {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("error retrieving service instance label: %s", err)
	}

	// Deserialize instance type.
	instanceType, err := deserializeServiceInstanceType(instanceResponse)
	if err != nil {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("error deserializing service instance type: %s", err)
	}

	// Grab service bindings from service_bindings_url.
	bindingsResponse, err := self.getFromCfApi(instances[0].ServiceBindingsUrl, cliConnection)
	if err != nil {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("error retrieving service bindings: %s", err)
	}

	var mysqlCredentials []pluginModels.MysqlCredentials

	if instanceType.Type == "p.mysql" {
		_, mysqlCredentials, err = deserializePMysqlServiceBindings(bindingsResponse)
	} else if instanceType.Type == "aws-rds-mysql" {
		_, mysqlCredentials, err = deserializeAwsRdsMysqlServiceBindings(bindingsResponse)
	} else if instanceType.Type == "rdsmysql" {
		_, mysqlCredentials, err = deserializeRdsMysqlServiceBindings(bindingsResponse)
	} else {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("unsupported service type: %s", instanceType.Type)
	}

	if err != nil {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("error retrieving service bindings: %s", err)
	}

	// Ensure apps are bound.
	if len(mysqlCredentials) == 0 {
		return pluginModels.MysqlCredentials{}, fmt.Errorf("%s no bound apps", name)
	}

	//return pluginModels.MysqlCredentials{}, fmt.Errorf("TEST: %s", mysqlCredentials[0])

	return mysqlCredentials[0], nil
}

func (self *apiClient) GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error) {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve apps: %s", err)
	}

	startedApps := make([]sdkModels.GetAppsModel, 0, len(apps))

	for _, app := range apps {
		if app.State == "started" {
			startedApps = append(startedApps, app)
		}
	}

	return startedApps, nil
}

func (self *apiClient) getFromCfApi(path string, cliConnection plugin.CliConnection) ([]byte, error) {
	config, err := self.getCliConfig(cliConnection)
	if err != nil {
		return nil, err
	}

	return self.httpClient.Get(config.ApiEndpoint+path, config.AccessToken, config.SslDisabled)
}

func (self *apiClient) postToCfApi(path string, body io.Reader, cliConnection plugin.CliConnection) ([]byte, error) {
	config, err := self.getCliConfig(cliConnection)
	if err != nil {
		return nil, err
	}

	return self.httpClient.Post(config.ApiEndpoint+path, body, config.AccessToken, config.SslDisabled)
}

func (self *apiClient) getCliConfig(cliConnection plugin.CliConnection) (*CliConfig, error) {
	if self.cliConfig == nil {
		endpoint, err := cliConnection.ApiEndpoint()
		if err != nil {
			return nil, fmt.Errorf("unable to get API endpoint: %s", err)
		}

		accessToken, err := cliConnection.AccessToken()
		if err != nil {
			return nil, fmt.Errorf("unable to get access token: %s", err)
		}

		sslDisabled, err := cliConnection.IsSSLDisabled()
		if err != nil {
			return nil, fmt.Errorf("unable to check SSL status: %s", err)
		}

		self.cliConfig = &CliConfig{
			AccessToken: accessToken,
			ApiEndpoint: endpoint,
			SslDisabled: sslDisabled,
		}
	}

	return self.cliConfig, nil
}

func deserializeInstances(jsonResponse []byte) (string, []pluginModels.ServiceInstance, error) {
	paginatedResources := new(resources.PaginatedServiceInstanceResources)
	err := json.Unmarshal(jsonResponse, paginatedResources)

	if err != nil {
		return "", nil, fmt.Errorf("unable to deserialize service instances: %s", err)
	}

	return paginatedResources.NextUrl, paginatedResources.ToModel(), nil
}

func deserializeAwsRdsMysqlServiceBindings(jsonResponse []byte) (string, []pluginModels.MysqlCredentials, error) {
	paginatedResources := new(resources.PaginatedAwsRdsMysqlServiceBindingsResources)
	err := json.Unmarshal(jsonResponse, paginatedResources)

	if err != nil {
		return "", nil, fmt.Errorf("unable to deserialize binding credentials: %s", err)
	}

	return paginatedResources.NextUrl, paginatedResources.ToModel(), nil
}
func deserializeRdsMysqlServiceBindings(jsonResponse []byte) (string, []pluginModels.MysqlCredentials, error) {
	paginatedResources := new(resources.PaginatedRdsMysqlServiceBindingsResources)
	err := json.Unmarshal(jsonResponse, paginatedResources)

	if err != nil {
		return "", nil, fmt.Errorf("unable to deserialize binding credentials: %s", err)
	}

	return paginatedResources.NextUrl, paginatedResources.ToModel(), nil
}
func deserializePMysqlServiceBindings(jsonResponse []byte) (string, []pluginModels.MysqlCredentials, error) {
	paginatedResources := new(resources.PaginatedPMysqlServiceBindingsResources)
	err := json.Unmarshal(jsonResponse, paginatedResources)

	if err != nil {
		return "", nil, fmt.Errorf("unable to deserialize binding credentials: %s", err)
	}

	return paginatedResources.NextUrl, paginatedResources.ToModel(), nil
}

func deserializeServiceInstanceType(serviceResponse []byte) (pluginModels.ServiceInstanceType, error) {
	resource := new(resources.ServiceInstanceTypeResource)
	err := json.Unmarshal(serviceResponse, resource)
	if err != nil {
		return pluginModels.ServiceInstanceType{}, fmt.Errorf("error deserializing service instance type response: %s", err)
	}

	serviceInstanceType, err := resource.ToModel()
	if err != nil {
		return pluginModels.ServiceInstanceType{}, fmt.Errorf("error converting service instance type response: %s", err)
	}

	return serviceInstanceType, nil
}

type ServiceKeyRequest struct {
	Name                string `json:"name"`
	ServiceInstanceGuid string `json:"service_instance_guid"`
}

type CliConfig struct {
	AccessToken string
	ApiEndpoint string
	SslDisabled bool
}
