package resources

import (
	"code.cloudfoundry.org/cli/cf/api/resources"
	"github.com/elliott-neal/cf-mysql-plugin/cfmysql/models"
	"strconv"
	"strings"
)

//type MysqlCredentials struct {
//	Uri      string `json:"uri"`
//	DbName   string `json:"name"`
//	Hostname string `json:"hostname"`
//	Port     string
//	RawPort  json.RawMessage `json:"port"`
//	Username string          `json:"username"`
//	Password string          `json:"password"`
//	Tls      TlsResource     `json:"tls"`
//}

type TlsResource struct {
	Cert TlsCertResource `json:"cert"`
}

type TlsCertResource struct {
	Ca string `json:"ca"`
}

type PaginatedServiceInstanceResources struct {
	TotalResults int    `json:"total_results"`
	NextUrl      string `json:"next_url"`
	Resources    []ServiceInstanceResource
}

type PaginatedAwsRdsMysqlServiceBindingsResources struct {
	TotalResults int    `json:"total_results"`
	NextUrl      string `json:"next_url"`
	Resources    []ServiceAwsRdsMysqlBindingsResource
}

type ServiceAwsRdsMysqlBindingsResource struct {
	resources.Resource
	Entity ServiceAwsRdsMysqlBindingsEntity
}

type ServiceAwsRdsMysqlBindingsEntity struct {
	Credentials ServiceAwsRdsMysqlBindingsCredentials
}

type ServiceAwsRdsMysqlBindingsCredentials struct {
	Uri		 string `json:"uri"`
	Database string `json:"database"`
	Hostname string `json:"hostname"`
	Port	 int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type PaginatedRdsMysqlServiceBindingsResources struct {
	TotalResults int    `json:"total_results"`
	NextUrl      string `json:"next_url"`
	Resources    []ServiceRdsMysqlBindingsResource
}

type ServiceRdsMysqlBindingsResource struct {
	resources.Resource
	Entity ServiceRdsMysqlBindingsEntity
}

type ServiceRdsMysqlBindingsEntity struct {
	Credentials ServiceRdsMysqlBindingsCredentials
}

type ServiceRdsMysqlBindingsCredentials struct {
	Database string `json:"DB_NAME"`
	Hostname string `json:"ENDPOINT_ADDRESS"`
	Port	 string `json:"PORT"`
	Username string `json:"MASTER_USERNAME"`
	Password string `json:"MASTER_PASSWORD"`
}

type PaginatedPMysqlServiceBindingsResources struct {
	TotalResults int    `json:"total_results"`
	NextUrl      string `json:"next_url"`
	Resources    []ServicePMysqlBindingsResource
}

type ServicePMysqlBindingsResource struct {
	resources.Resource
	Entity ServicePMysqlBindingsEntity
}

type ServicePMysqlBindingsEntity struct {
	Credentials ServicePMysqlBindingsCredentials
}

type ServicePMysqlBindingsCredentials struct {
	Uri      string `json:"uri"`
	Database string `json:"name"`
	Hostname string `json:"hostname"`
	Port     int	`json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServiceInstanceTypeResource struct {
	Entity ServiceInstanceTypeEntity
}

type ServiceInstanceTypeEntity struct {
	Label string `json:"label"`
}

type ServiceInstanceResource struct {
	resources.Resource
	Entity ServiceInstanceEntity
}

type ServiceInstanceEntity struct {
	Name            	string                         		`json:"name"`
	DashboardURL    	string                         		`json:"dashboard_url"`
	Tags            	[]string                       		`json:"tags"`
	ServiceKeys     	[]resources.ServiceKeyResource 		`json:"service_keys"`
	ServicePlan     	resources.ServicePlanResource  		`json:"service_plan"`
	LastOperation   	resources.LastOperation        		`json:"last_operation"`
	SpaceUrl        	string                         		`json:"space_url"`
	ServiceUrl			string 						   		`json:"service_url"`
	ServiceBindingsUrl	string								`json:"service_bindings_url"`
}

func (self *PaginatedServiceInstanceResources) ToModel() []models.ServiceInstance {
	var convertedModels []models.ServiceInstance

	for _, resource := range self.Resources {
		model := models.ServiceInstance{}
		model.Guid = resource.Metadata.GUID
		model.Name = resource.Entity.Name
		model.ServiceUrl = resource.Entity.ServiceUrl
		model.ServiceBindingsUrl = resource.Entity.ServiceBindingsUrl

		pathParts := strings.Split(resource.Entity.SpaceUrl, "/")
		model.SpaceGuid = pathParts[len(pathParts)-1]

		convertedModels = append(convertedModels, model)
	}

	return convertedModels
}

func (self *PaginatedAwsRdsMysqlServiceBindingsResources) ToModel() []models.MysqlCredentials {
	var convertedModels []models.MysqlCredentials

	for _, resource := range self.Resources {
		model := models.MysqlCredentials{}
		model.Uri = resource.Entity.Credentials.Uri
		model.Hostname = resource.Entity.Credentials.Hostname
		model.Port = strconv.Itoa(resource.Entity.Credentials.Port)
		model.Database = resource.Entity.Credentials.Database
		model.Username = resource.Entity.Credentials.Username
		model.Password = resource.Entity.Credentials.Password

		convertedModels = append(convertedModels, model)
	}

	return convertedModels
}

func (self *PaginatedRdsMysqlServiceBindingsResources) ToModel() []models.MysqlCredentials {
	var convertedModels []models.MysqlCredentials

	for _, resource := range self.Resources {
		model := models.MysqlCredentials{}
		model.Hostname = resource.Entity.Credentials.Hostname
		model.Port = resource.Entity.Credentials.Port
		model.Database = resource.Entity.Credentials.Database
		model.Username = resource.Entity.Credentials.Username
		model.Password = resource.Entity.Credentials.Password

		convertedModels = append(convertedModels, model)
	}

	return convertedModels
}

func (self *PaginatedPMysqlServiceBindingsResources) ToModel() []models.MysqlCredentials {
	var convertedModels []models.MysqlCredentials

	for _, resource := range self.Resources {
		model := models.MysqlCredentials{}
		model.Uri = resource.Entity.Credentials.Uri
		model.Hostname = resource.Entity.Credentials.Hostname
		model.Port = strconv.Itoa(resource.Entity.Credentials.Port)
		model.Database = resource.Entity.Credentials.Database
		model.Username = resource.Entity.Credentials.Username
		model.Password = resource.Entity.Credentials.Password

		convertedModels = append(convertedModels, model)
	}

	return convertedModels
}

func (self *ServiceInstanceTypeResource) ToModel() (models.ServiceInstanceType, error) {
	return models.ServiceInstanceType {
		Type: self.Entity.Label,
	}, nil
}
