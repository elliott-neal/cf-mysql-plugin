package resources_test

import (
	. "github.com/elliott-neal/cf-mysql-plugin/cfmysql/resources"

	"code.cloudfoundry.org/cli/cf/api/resources"
	"encoding/json"
	//"errors"
	"github.com/elliott-neal/cf-mysql-plugin/cfmysql/models"
	"github.com/elliott-neal/cf-mysql-plugin/cfmysql/test_resources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resources", func() {
	Describe("Service instances", func() {
		Context("Deserializing JSON", func() {
			It("can get service instance names and guids", func() {
				paginatedResources := new(PaginatedServiceInstanceResources)
				err := json.Unmarshal(test_resources.LoadResource("../test_resources/service_instances.json"), paginatedResources)

				Expect(err).To(BeNil())
				Expect(paginatedResources.Resources).To(HaveLen(4))
				Expect(paginatedResources.NextUrl).To(Equal("/v2/service_instances?page=2"))

				Expect(paginatedResources.Resources[0].Entity.Name).To(Equal("database-a"))
				Expect(paginatedResources.Resources[0].Metadata.GUID).To(Equal("service-instance-guid-a"))
				Expect(paginatedResources.Resources[0].Entity.SpaceUrl).To(Equal("/v2/spaces/space-guid"))

				Expect(paginatedResources.Resources[1].Entity.Name).To(Equal("database-b"))
				Expect(paginatedResources.Resources[1].Metadata.GUID).To(Equal("service-instance-guid-b"))
				Expect(paginatedResources.Resources[1].Entity.SpaceUrl).To(Equal("/v2/spaces/space-guid"))
			})
		})

		Context("Converting to models", func() {
			It("Converts very nicely", func() {
				resourceInstances := &PaginatedServiceInstanceResources{
					Resources: []ServiceInstanceResource{
						{
							Resource: resources.Resource{
								Metadata: resources.Metadata{
									GUID: "fine-guid",
								},
							},
							Entity: ServiceInstanceEntity{
								Name:     "outstanding-service-name",
								SpaceUrl: "/v2/spaces/outer-space-guid",
							},
						},
						{
							Resource: resources.Resource{
								Metadata: resources.Metadata{
									GUID: "better-guid",
								},
							},
							Entity: ServiceInstanceEntity{
								Name:     "best-service-name",
								SpaceUrl: "/v2/spaces/inner-space-guid",
							},
						},
					},
				}

				instanceModels := resourceInstances.ToModel()

				Expect(instanceModels).To(HaveLen(2))
				Expect(instanceModels[0]).To(Equal(models.ServiceInstance{
					Name:      "outstanding-service-name",
					Guid:      "fine-guid",
					SpaceGuid: "outer-space-guid",
				}))
				Expect(instanceModels[1]).To(Equal(models.ServiceInstance{
					Name:      "best-service-name",
					Guid:      "better-guid",
					SpaceGuid: "inner-space-guid",
				}))
			})
		})
	})

	Describe("Service Bindings", func() {
		Context("Deserializing JSON", func() {
			It("can get credentials", func() {
				paginatedResources := new(PaginatedPMysqlServiceBindingsResources)
				err := json.Unmarshal(test_resources.LoadResource("../test_resources/service_bindings.json"), paginatedResources)


				Expect(err).To(BeNil())
				Expect(paginatedResources.Resources).To(HaveLen(2))
				Expect(paginatedResources.NextUrl).To(Equal("next-url"))
				//Expect(paginatedResources.Resources[0].Entity.ServiceInstanceGuid).To(Equal("service-instance-guid"))
				Expect(paginatedResources.Resources[0].Entity.Credentials.Uri).To(Equal("uri"))
				Expect(paginatedResources.Resources[0].Entity.Credentials.Database).To(Equal("db-name"))
				Expect(paginatedResources.Resources[0].Entity.Credentials.Hostname).To(Equal("hostname"))
				Expect(paginatedResources.Resources[0].Entity.Credentials.Username).To(Equal("username"))
				Expect(paginatedResources.Resources[0].Entity.Credentials.Password).To(Equal("password"))
				Expect(paginatedResources.Resources[0].Entity.Credentials.Port).To(Equal(3306))
				Expect(paginatedResources.Resources[1].Entity.Credentials.Port).To(Equal(2342))
				//Expect(paginatedResources.Resources[0].Entity.Credentials.Tls.Cert.Ca).To(Equal("ca-certificate"))

				//var portString string
				//err = json.Unmarshal(paginatedResources.Resources[0].Entity.Credentials.Port, &portString)
				//Expect(err).To(BeNil())
				//Expect(portString).To(Equal("3306"))

				//var portInt int
				//err = json.Unmarshal(paginatedResources.Resources[1].Entity.Credentials.Port, &portInt)
				//Expect(err).To(BeNil())
				//Expect(portInt).To(Equal(2342))
			})
		})

		Context("Converting to models", func() {
			It("Converts very nicely if the port is string or int", func() {
				resourceBindings := &PaginatedPMysqlServiceBindingsResources{
					Resources: []ServicePMysqlBindingsResource{
						{
							Entity: ServicePMysqlBindingsEntity{
								//ServiceInstanceGuid: "service-instance-guid-a",
								Credentials: ServicePMysqlBindingsCredentials{
									Uri:      "uri-a",
									Database:   "db-name-a",
									Hostname: "hostname-a",
									//Port:  []byte("\"1234\""),
									Port: 1234,
									Username: "username-a",
									Password: "password-a",
									//Tls: TlsResource{
									//	Cert: TlsCertResource{
									//		Ca: "ca-certificate-a",
									//	},
									//},
								},
							},
						},
						{
							Entity: ServicePMysqlBindingsEntity{
								//ServiceInstanceGuid: "service-instance-guid-b",
								Credentials: ServicePMysqlBindingsCredentials{
									Uri:      "uri-b",
									Database:   "db-name-b",
									Hostname: "hostname-b",
									//Port:  []byte("2345"),
									Port: 2345,
									Username: "username-b",
									Password: "password-b",
								},
							},
						},
					},
				}

				//serviceKeys, err := resourceBindings.ToModel()
				serviceBindings := resourceBindings.ToModel()

				//Expect(err).To(BeNil())

				Expect(serviceBindings).To(HaveLen(2))
				Expect(serviceBindings[0]).To(Equal(models.MysqlCredentials{
					//ServiceInstanceGuid: "service-instance-guid-a",
					Uri:                 "uri-a",
					Database:              "db-name-a",
					Hostname:            "hostname-a",
					Port:                "1234",
					Username:            "username-a",
					Password:            "password-a",
					//CaCert:              "ca-certificate-a",
				}))
				Expect(serviceBindings[1]).To(Equal(models.MysqlCredentials{
					//ServiceInstanceGuid: "service-instance-guid-b",
					Uri:                 "uri-b",
					Database:              "db-name-b",
					Hostname:            "hostname-b",
					Port:                "2345",
					Username:            "username-b",
					Password:            "password-b",
					//CaCert:              "",
				}))
			})

			//It("Returns an error if the port is not string or int", func() {
			//	keyResources := &PaginatedPMysqlServiceBindingsResources{
			//		Resources: []ServicePMysqlBindingsResource{
			//			{
			//				Entity: ServicePMysqlBindingsEntity{
			//					//ServiceInstanceGuid: "service-instance-guid-b",
			//					Credentials: ServicePMysqlBindingsCredentials{
			//						Uri:      "uri-b",
			//						Database:   "db-name-b",
			//						Hostname: "hostname-b",
			//						//Port:  []byte("false"),
			//						Port: 1234,
			//						Username: "username-b",
			//						Password: "password-b",
			//					},
			//				},
			//			},
			//		},
			//	}
			//
			//	//serviceKeys, err := keyResources.ToModel()
			//	serviceBindings := keyResources.ToModel()
			//
			//	//Expect(err).To(Equal(errors.New("unable to deserialize port in service key: 'false'")))
			//	Expect(serviceBindings).To(BeNil())
			//})
		})
	})
})
