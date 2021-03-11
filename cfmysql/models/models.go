package models

type ServiceInstance struct {
	Name       string
	Guid       string
	SpaceGuid  string
	ServiceUrl string
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
