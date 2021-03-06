package cfmysql

import (
	"code.cloudfoundry.org/cli/plugin"
	sdkModels "code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	pluginModels "github.com/elliott-neal/cf-mysql-plugin/cfmysql/models"
	"io"
)

//go:generate counterfeiter . CfService
type CfService interface {
	GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error)
	OpenSshTunnel(cliConnection plugin.CliConnection, toService MysqlService, apps []sdkModels.GetAppsModel, localPort int)
	GetService(connection plugin.CliConnection, name string) (MysqlService, error)
}

func NewCfService(apiClient ApiClient, runner SshRunner, waiter PortWaiter, httpClient HttpWrapper, randWrapper RandWrapper, logWriter io.Writer) *cfService {
	return &cfService{
		apiClient:   apiClient,
		sshRunner:   runner,
		portWaiter:  waiter,
		httpClient:  httpClient,
		randWrapper: randWrapper,
		logWriter:   logWriter,
	}
}

const ServiceKeyName = "cf-mysql"

type MysqlService struct {
	Name     string
	Hostname string
	Port     string
	DbName   string
	Username string
	Password string
	CaCert   string
}

type cfService struct {
	apiClient   ApiClient
	httpClient  HttpWrapper
	portWaiter  PortWaiter
	sshRunner   SshRunner
	randWrapper RandWrapper
	logWriter   io.Writer
}

func (self *cfService) GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error) {
	return self.apiClient.GetStartedApps(cliConnection)
}

func (self *cfService) OpenSshTunnel(cliConnection plugin.CliConnection, toService MysqlService, apps []sdkModels.GetAppsModel, localPort int) {
	throughAppIndex := self.randWrapper.Intn(len(apps))
	throughApp := apps[throughAppIndex].Name
	go self.sshRunner.OpenSshTunnel(cliConnection, toService, throughApp, localPort)

	self.portWaiter.WaitUntilOpen(localPort)
}

func (self *cfService) GetService(connection plugin.CliConnection, name string) (MysqlService, error) {
	space, err := connection.GetCurrentSpace()
	if err != nil {
		return MysqlService{}, fmt.Errorf("unable to retrieve current space: %s", err)
	}

	credentials, err := self.apiClient.GetService(connection, space.Guid, name)
	if err != nil {
		return MysqlService{}, fmt.Errorf("unable to retrieve metadata for service %s: %s", name, err)
	}

	return toServiceModel(name, credentials), nil

	//if err != nil {
	//	return MysqlService{}, fmt.Errorf("unable to retrieve metadata for service %s: %s", name, err)
	//}
	//
	//serviceKey, found, err := self.apiClient.GetServiceKey(connection, instance.Guid, instance.ServiceUrl, ServiceKeyName)
	//if err != nil {
	//	return MysqlService{}, fmt.Errorf("unable to retrieve service key: %s", err)
	//}
	//
	//if found {
	//	return toServiceModel(name, serviceKey), nil
	//}
	//
	//fmt.Fprintf(self.logWriter, "Creating new service key %s for %s...\n", ServiceKeyName, name)
	//serviceKey, err = self.apiClient.CreateServiceKey(connection, instance.Guid, instance.ServiceUrl, ServiceKeyName)
	//if err != nil {
	//	return MysqlService{}, fmt.Errorf("unable to create service key: %s", err)
	//}
	//
}

func toServiceModel(name string, serviceKey pluginModels.MysqlCredentials) MysqlService {
	return MysqlService{
		Name:     name,
		Hostname: serviceKey.Hostname,
		Port:     serviceKey.Port,
		DbName:   serviceKey.Database,
		Username: serviceKey.Username,
		Password: serviceKey.Password,
	}
}
