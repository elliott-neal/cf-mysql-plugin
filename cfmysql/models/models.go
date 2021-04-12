package models

type ServiceInstance struct {
	Name               string
	Guid               string
	SpaceGuid          string
	ServiceUrl         string
	ServiceBindingsUrl string
}

type ServiceInstanceType struct {
	Type 	string
}

type ServiceKey struct {
	ServiceInstanceGuid string
	Uri                 string
	DbName              string
	Hostname            string
	Port                string
	Username            string
	Password            string
	CaCert              string
}

type ServiceBindings struct {
	AppUrl string
}

type MysqlCredentials struct {
	Uri		 string
	Hostname string
	Port 	 string
	Database string
	Username string
	Password string
}
