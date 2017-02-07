package goarubacloud

import (
	"github.com/arubacloud/goarubacloud/models"
	"github.com/docker/machine/libmachine/log"
	"fmt"
)

type API struct {
	client *Client
}

type SetEnqueueServerCreation struct {
	Username      string `json:"Username"`
	Password      string `json:"Password"`
	Server        struct {
			      AdministratorPassword string `json:"AdministratorPassword"`
			      Name                  string `json:"Name"`
			      SmartVMWarePackageID  int    `json:"SmartVMWarePackageID"`
			      Note                  string `json:"Note"`
			      OSTemplateId          int    `json:"OSTemplateId"`
		      }
}

type GetServerDetailsRequest struct {
	Username      string        `json:"Username,omitempty"`
	Password      string        `json:"Password,omitempty"`
	ServerId      int        `json:"ServerId,omitempty"`
}

type SetEnqueueServerDeletion struct {
	Username      string        `json:"Username,omitempty"`
	Password      string        `json:"Password,omitempty"`
	ServerId      int        `json:"ServerId,omitempty"`
}

type Response struct {
	Success bool        `json:"Success,omitempty"`
	Value   string        `json:"Value,omitempty"`
}

type SetEnqueueServerStart struct {
	Username      string        `json:"Username,omitempty"`
	Password      string        `json:"Password,omitempty"`
	ServerId      int        `json:"ServerId,omitempty"`
}

type SetEnqueueServerStop struct {
	Username      string        `json:"Username,omitempty"`
	Password      string        `json:"Password,omitempty"`
	ServerId      int        `json:"ServerId,omitempty"`
}

type GetServersRequest struct {
	Username      string        `json:"Username"`
	Password      string        `json:"Password"`
}

type GetTemplatesRequest struct {
	Username      string        `json:"Username"`
	Password      string        `json:"Password"`
}

type GetPreconfiguredPackagesRequest struct {
	Username      string        `json:"Username"`
	Password      string        `json:"Password"`
	HypervisorType int 			`json:"HypervisorType"`
}

func NewAPI(endpoint, username, password string) (api *API, err error) {
	client, err := NewClient(endpoint, username, password)
	if err != nil {
		return nil, err
	}
	return &API{client}, nil
}

func (a *API) GetTemplates() (hypervisorTypes []*models.GetHypvervisorTypeResponse, err error) {
	var getTemplatesRequest = GetTemplatesRequest{}
	getTemplatesRequest.Username = a.client.Username
	getTemplatesRequest.Password = a.client.Password
	err = a.client.Post("/GetHypervisors", getTemplatesRequest, &hypervisorTypes)
	if err != nil {
		return nil, err
	}

	return hypervisorTypes, nil
}

func (a *API) GetPreconfiguredPackages() (packagesTypes []*models.CloudPackage, err error) {
	var getpackageRequest = GetPreconfiguredPackagesRequest{}
	getpackageRequest.Username = a.client.Username
	getpackageRequest.Password = a.client.Password
	getpackageRequest.HypervisorType = 4
	err = a.client.Post("/GetPreconfiguredPackages", getpackageRequest, &packagesTypes)
	if err != nil {
		return nil, err
	}

	return packagesTypes, nil
}

func (a *API) GetPreconfiguredPackage(packagename string) (cloudpackage *models.CloudPackage, err error) {
	cloudpackages, err := a.GetPreconfiguredPackages()
	if err != nil {
		return nil, err
	}
	for _, cloudpackage := range cloudpackages {
			for _, description := range cloudpackage.Descriptions {
				if description.LanguageID == 2 && description.Text == packagename {
					return cloudpackage, nil
				}
			}
		
	}

	return nil,
		fmt.Errorf("No package found with Name: %s found on datacenter: %s", packagename, a.client.EndPoint)
}

func (a *API) GetTemplate(templatename string) (template *models.Template, err error) {
	templates, err := a.GetTemplates()
	if err != nil {
		return nil, err
	}
	for _, hv := range templates {
		if hv.HypervisorType == 4 {
			for _, template := range hv.Templates {
				if template.Name == templatename {
					return &template, nil
				}
			}
		}
	}

	return nil,
		fmt.Errorf("No template with Name: %s found on datacenter: %s", templatename, a.client.EndPoint)
}



func (a *API) GetServers() (servers []*models.Server, err error) {
	var getServersRequest = GetServersRequest{}
	getServersRequest.Username = a.client.Username
	getServersRequest.Password = a.client.Password
	err = a.client.Post("/GetServers", getServersRequest, &servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (a *API) GetServer(serverId int) (server *models.Server, err error) {
	var getServerDetailsRequest = GetServerDetailsRequest{}
	getServerDetailsRequest.Username = a.client.Username
	getServerDetailsRequest.Password = a.client.Password
	getServerDetailsRequest.ServerId = serverId

	err = a.client.Post("/GetServerDetails", getServerDetailsRequest, &server)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (a *API) CreateServer(name, admin_password string, package_id, os_template_id int) (server *models.Server, err error) {
	var createRequest SetEnqueueServerCreation
	createRequest.Username = a.client.Username
	createRequest.Password = a.client.Password
	createRequest.Server.AdministratorPassword = admin_password
	createRequest.Server.Name = name
	createRequest.Server.OSTemplateId = os_template_id
	createRequest.Server.SmartVMWarePackageID = package_id

	log.Debug("Post CreateServer Request.")
	err = a.client.Post("/SetEnqueueServerCreation", createRequest, &server)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (a *API) DeleteServer(server_id int) (err error) {
	var deleteServer SetEnqueueServerDeletion
	deleteServer.Username = a.client.Username
	deleteServer.Password = a.client.Password
	deleteServer.ServerId = server_id

	err = a.client.Post("/SetEnqueueServerDeletion", deleteServer, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) StartServer(server_id int) (err error) {
	var startServer SetEnqueueServerStart
	startServer.Username = a.client.Username
	startServer.Password = a.client.Password
	startServer.ServerId = server_id

	err = a.client.Post("/SetEnqueueServerStart", startServer, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) StopServer(server_id int) (err error) {
	var stopServer SetEnqueueServerStop
	stopServer.Username = a.client.Username
	stopServer.Password = a.client.Password
	stopServer.ServerId = server_id

	err = a.client.Post("/SetEnqueueServerStop", stopServer, nil)
	if err != nil {
		return err
	}

	return nil
}
