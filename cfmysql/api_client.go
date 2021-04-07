package cfmysql

import (
	"encoding/json"
	"fmt"
	"os"

	"bytes"
	"code.cloudfoundry.org/cli/plugin"
	sdkModels "code.cloudfoundry.org/cli/plugin/models"
	pluginModels "github.com/elliott-neal/cf-mysql-plugin/cfmysql/models"
	"github.com/elliott-neal/cf-mysql-plugin/cfmysql/resources"
	"io"
	"net/url"
)

//go:generate counterfeiter . ApiClient
type ApiClient interface {
	GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error)
	GetService(cliConnection plugin.CliConnection, spaceGuid string, name string) (pluginModels.ServiceInstance, error)
	GetServiceKey(cliConnection plugin.CliConnection, serviceInstanceGuid string, serviceInstanceUrl string, keyName string) (key pluginModels.ServiceKey, found bool, err error)
	CreateServiceKey(cliConnection plugin.CliConnection, serviceInstanceGuid string, serviceInstanceUrl string, keyName string) (pluginModels.ServiceKey, error)
	GetInstanceType(cliConnection plugin.CliConnection, serviceUrl string) (string, error)
}

func NewApiClient(httpClient HttpWrapper) *apiClient {
	return &apiClient{
		httpClient: httpClient,
	}
}

type apiClient struct {
	httpClient HttpWrapper
	cliConfig  *CliConfig
}

func (self *apiClient) GetInstanceType(cliConnection plugin.CliConnection, serviceUrl string) (string, error) {
	instanceTypeResponse, err := self.getFromCfApi(serviceUrl, cliConnection)
	if err != nil {
		return "", fmt.Errorf("error retrieving service instance: %s", err)
	}

	instance, err := deserializeServiceInstanceType(instanceTypeResponse)
	if err != nil {
		return "", fmt.Errorf("error deserializing service instance type: %s", err)
	}

	return instance.Type, nil
}

func (self *apiClient) GetService(cliConnection plugin.CliConnection, spaceGuid string, name string) (pluginModels.ServiceInstance, error) {
	path := fmt.Sprintf(
		"/v2/spaces/%s/service_instances?return_user_provided_service_instances=true&q=name%%3A%s",
		spaceGuid,
		url.QueryEscape(name),
	)

	instanceResponse, err := self.getFromCfApi(path, cliConnection)
	if err != nil {
		return pluginModels.ServiceInstance{}, fmt.Errorf("error retrieving service instance: %s", err)
	}

	_, instances, err := deserializeInstances(instanceResponse)
	if err != nil {
		return pluginModels.ServiceInstance{}, fmt.Errorf("error deserializing service instances: %s", err)
	}

	if len(instances) == 0 {
		return pluginModels.ServiceInstance{}, fmt.Errorf("%s not found in current space", name)
	}

	return instances[0], nil
}

func (self *apiClient) GetServiceKey(cliConnection plugin.CliConnection, serviceInstanceGuid string, serviceInstanceUrl string, keyName string) (pluginModels.ServiceKey, bool, error) {
	path := fmt.Sprintf(
		"/v2/service_instances/%s/service_keys?q=name%%3A%s",
		serviceInstanceGuid,
		url.QueryEscape(keyName),
	)

	keyResponse, err := self.getFromCfApi(path, cliConnection)
	if err != nil {
		return pluginModels.ServiceKey{}, false, fmt.Errorf("error retrieving service key: %s", err)
	}

	instancetype, err := self.GetInstanceType(cliConnection, serviceInstanceUrl)
	if err != nil {
		return pluginModels.ServiceKey{}, false, fmt.Errorf("unable to determine service type: %s", err)
	}

	fmt.Println(instancetype)

	var serviceKeys []pluginModels.ServiceKey

	if instancetype == "p.mysql" {
		serviceKeys, err = deserializePMysqlServiceKeys(keyResponse)
	} else if instancetype == "aws-rds-mysql" {
		serviceKeys, err = deserializeAwsRdsMysqlServiceKeys(keyResponse)
	} else if instancetype == "rdsmysql" {
		serviceKeys, err = deserializeRdsMysqlServiceKeys(keyResponse)
	} else {
		return pluginModels.ServiceKey{}, false, fmt.Errorf("unsupported service type: %s", instancetype)
	}

	fmt.Fprintf(os.Stdout,"Retrieving service Key:%s\n", serviceKeys)

	if err != nil {
		return pluginModels.ServiceKey{}, false, fmt.Errorf("error deserializing service key response: %s", err)
	}

	if len(serviceKeys) == 0 {
		return pluginModels.ServiceKey{}, false, nil
	}

	return serviceKeys[0], true, nil
}

func (self *apiClient) CreateServiceKey(cliConnection plugin.CliConnection, serviceInstanceGuid string, serviceInstanceUrl string, keyName string) (pluginModels.ServiceKey, error) {
	content := ServiceKeyRequest{
		Name:                keyName,
		ServiceInstanceGuid: serviceInstanceGuid,
	}

	body, err := json.Marshal(content)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error serializing request body: %s", err)
	}

	response, err := self.postToCfApi("/v2/service_keys", bytes.NewBuffer(body), cliConnection)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error creating service key: %s", err)
	}

	instancetype, err := self.GetInstanceType(cliConnection, serviceInstanceUrl)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("unable to determine service type: %s", err)
	}

	fmt.Fprintf(os.Stdout,"Response received %s, for instance: %s", response, instancetype)

	var serviceKey pluginModels.ServiceKey

	if instancetype == "p.mysql" {
		serviceKey, err = deserializePMysqlServiceKey(response)
	} else if instancetype == "aws-rds-mysql" {
		serviceKey, err = deserializeAwsRdsMysqlServiceKey(response)
	} else if instancetype == "rdsmysql" {
		serviceKey, err = deserializeRdsMysqlServiceKey(response)
	} else {
		return pluginModels.ServiceKey{}, fmt.Errorf("unsupported service type: %s", instancetype)
	}
	fmt.Fprintf(os.Stdout,"Printing service Key:%s\n", serviceKey)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error deserializing service key response: %s", err)
	}

	return serviceKey, nil
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

func deserializePMysqlServiceKeys(keyResponse []byte) ([]pluginModels.ServiceKey, error) {
	paginatedResources := new(resources.PaginatedPMysqlServiceKeyResources)
	err := json.Unmarshal(keyResponse, paginatedResources)
	if err != nil {
		return nil, fmt.Errorf("error deserializing service key response: %s", err)
	}

	serviceKeys, err := paginatedResources.ToModel()
	if err != nil {
		return nil, fmt.Errorf("error converting service key response: %s", err)
	}

	return serviceKeys, nil
}

func deserializeAwsRdsMysqlServiceKeys(keyResponse []byte) ([]pluginModels.ServiceKey, error) {
	paginatedResources := new(resources.PaginatedAwsRdsMysqlServiceKeyResources)
	err := json.Unmarshal(keyResponse, paginatedResources)
	if err != nil {
		return nil, fmt.Errorf("error deserializing service key response: %s", err)
	}

	serviceKeys, err := paginatedResources.ToModel()
	if err != nil {
		return nil, fmt.Errorf("error converting service key response: %s", err)
	}

	return serviceKeys, nil
}

func deserializeRdsMysqlServiceKeys(keyResponse []byte) ([]pluginModels.ServiceKey, error) {
	paginatedResources := new(resources.PaginatedRdsMysqlServiceKeyResources)
	err := json.Unmarshal(keyResponse, paginatedResources)
	if err != nil {
		return nil, fmt.Errorf("error deserializing service key response: %s", err)
	}

	serviceKeys, err := paginatedResources.ToModel()
	if err != nil {
		return nil, fmt.Errorf("error converting service key response: %s", err)
	}

	return serviceKeys, nil
}

func deserializePMysqlServiceKey(keyResponse []byte) (pluginModels.ServiceKey, error) {
	resource := new(resources.PMysqlServiceKeyResource)
	err := json.Unmarshal(keyResponse, resource)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error deserializing service key response: %s", err)
	}

	serviceKey, err := resource.ToModel()
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error converting service key response: %s", err)
	}

	return serviceKey, nil
}

func deserializeAwsRdsMysqlServiceKey(keyResponse []byte) (pluginModels.ServiceKey, error) {
	resource := new(resources.AwsRdsMysqlServiceKeyResource)
	err := json.Unmarshal(keyResponse, resource)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error deserializing service key response: %s", err)
	}

	serviceKey, err := resource.ToModel()
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error converting service key response: %s", err)
	}

	return serviceKey, nil
}

func deserializeRdsMysqlServiceKey(keyResponse []byte) (pluginModels.ServiceKey, error) {
	resource := new(resources.RdsMysqlServiceKeyResource)
	err := json.Unmarshal(keyResponse, resource)
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error deserializing service key response: %s", err)
	}

	serviceKey, err := resource.ToModel()
	if err != nil {
		return pluginModels.ServiceKey{}, fmt.Errorf("error converting service key response: %s", err)
	}

	return serviceKey, nil
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
