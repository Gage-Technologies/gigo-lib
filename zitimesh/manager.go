package zitimesh

import (
	"context"
	"crypto/x509"
	"fmt"
	"go.uber.org/atomic"
	"time"

	"github.com/gage-technologies/gigo-lib/config"
	"github.com/openziti/edge-api/rest_management_api_client"
	"github.com/openziti/edge-api/rest_management_api_client/identity"
	"github.com/openziti/edge-api/rest_management_api_client/service"
	"github.com/openziti/edge-api/rest_management_api_client/service_policy"
	"github.com/openziti/edge-api/rest_model"
	rest_model_edge "github.com/openziti/edge-api/rest_model"
	"github.com/openziti/edge-api/rest_util"
)

var (
	ErrIdentityExists = fmt.Errorf("identity already exists")
)

// Manager
//
// Manages the ziti mesh by creating identities, services, and service policies
type Manager struct {
	edge *atomic.Pointer[rest_management_api_client.ZitiEdgeManagement]
}

// NewManager
//
// Creates a new manager for the ziti mesh
func NewManager(cfg config.ZitiConfig) (*Manager, error) {
	// abort if there are no schemes
	if len(cfg.EdgeSchemes) == 0 {
		return nil, fmt.Errorf("no schemes defined")
	}

	// create address
	ctrlAddress := fmt.Sprintf("%s://%s", cfg.EdgeSchemes[0], cfg.EdgeHost)

	// retrieve the certs
	caCerts, err := rest_util.GetControllerWellKnownCas(ctrlAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get well-known CA certs: %w", err)
	}

	// create the CA pool
	caPool := x509.NewCertPool()
	for _, ca := range caCerts {
		caPool.AddCert(ca)
	}

	// create the edge client
	ok, err := rest_util.VerifyController(ctrlAddress, caPool)
	if err != nil {
		return nil, fmt.Errorf("failed to verify controller: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("controller verification failed")
	}

	// create the edge client
	client, err := rest_util.NewEdgeManagementClientWithUpdb(cfg.ManagementUser, cfg.ManagementPass, ctrlAddress, caPool)
	if err != nil {
		return nil, fmt.Errorf("failed to create edge client: %w", err)
	}

	// list identities to validate that we have a good connection
	limit := int64(1)
	params := &identity.ListIdentitiesParams{
		Context: context.Background(),
		Limit:   &limit,
	}
	_, err = client.Identity.ListIdentities(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to validate connection: %w", err)
	}

	// create a new manager
	manager := &Manager{
		edge: &atomic.Pointer[rest_management_api_client.ZitiEdgeManagement]{},
	}

	// auth the client
	err = manager.authClient(cfg.ManagementUser, cfg.ManagementPass, ctrlAddress, caPool)
	if err != nil {
		return nil, fmt.Errorf("failed to auth client: %w", err)
	}

	// launch the renewal client auth routine
	go manager.renewClientAuthRoutine(cfg.ManagementUser, cfg.ManagementPass, ctrlAddress, caPool)

	return manager, nil
}

// authClient
//
// Authenticate with the ziti manager
func (m *Manager) authClient(username string, password string, address string, caPool *x509.CertPool) error {
	// create the edge client
	client, err := rest_util.NewEdgeManagementClientWithUpdb(username, password, address, caPool)
	if err != nil {
		return fmt.Errorf("failed to create edge client: %w", err)
	}

	// list identities to validate that we have a good connection
	limit := int64(1)
	params := &identity.ListIdentitiesParams{
		Context: context.Background(),
		Limit:   &limit,
	}
	_, err = client.Identity.ListIdentities(params, nil)
	if err != nil {
		return fmt.Errorf("failed to validate connection: %w", err)
	}

	// swap the edge client
	m.edge.Store(client)

	return nil
}

// renewClientAuthRoutine
//
// Automatically renews the client auth every 20m
func (m *Manager) renewClientAuthRoutine(username string, password string, address string, caPool *x509.CertPool) {
	// create a ticker to renew the auth every 20m
	ticker := time.NewTicker(20 * time.Minute)
	for {
		select {
		case <-ticker.C:
			err := m.authClient(username, password, address, caPool)
			if err != nil {
				fmt.Printf("failed to renew client auth: %v\n", err)
			}
		}
	}
}

// CreateAgent
//
// Creates an agent in the ziti mesh
func (m *Manager) CreateAgent(id int64) (string, string, error) {
	// create our variables for the identity of the agent
	isAdmin := false
	name := fmt.Sprintf("gigo-ws-agent-%d", id)
	identityType := rest_model_edge.IdentityTypeDevice

	// create filter to search for service
	searchParam := identity.NewListIdentitiesParams()
	filter := fmt.Sprintf("name=\"%s\"", name)
	searchParam.Filter = &filter

	// query to see if a service already exists
	list, err := m.edge.Load().Identity.ListIdentities(searchParam, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to query for service: %w", err)
	}

	// create a variable to hold the id of the identity
	var identityId string

	// if the identity exists then we need to fail since agents can't be recreated
	if list != nil && len(list.Payload.Data) > 0 {
		return "", "", ErrIdentityExists
	}

	// create the request for the identity
	createIdentityReq := identity.NewCreateIdentityParams()
	createIdentityReq.Identity = &rest_model_edge.IdentityCreate{
		Enrollment:          &rest_model_edge.IdentityCreateEnrollment{Ott: true},
		IsAdmin:             &isAdmin,
		Name:                &name,
		RoleAttributes:      &rest_model.Attributes{"gigo-agents"},
		ServiceHostingCosts: nil,
		Tags:                nil,
		Type:                &identityType,
	}
	createIdentityReq.SetTimeout(10 * time.Second)

	// create the identity
	createIdentityRes, err := m.edge.Load().Identity.CreateIdentity(createIdentityReq, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create identity: %w", err)
	}

	// set the identity id
	identityId = createIdentityRes.Payload.Data.ID

	// retrieve the token for the identity
	params := &identity.DetailIdentityParams{
		Context: context.Background(),
		ID:      identityId,
	}
	params.SetTimeout(10 * time.Second)
	resp, err := m.edge.Load().Identity.DetailIdentity(params, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve identity: %w", err)
	}
	return identityId, resp.Payload.Data.Enrollment.Ott.JWT, nil
}

// DeleteAgent
//
// Deletes an agent from the ziti mesh
func (m *Manager) DeleteAgent(id int64) error {
	// create our variables for the identity of the agent
	name := fmt.Sprintf("gigo-ws-agent-%d", id)

	// create filter to search for service
	searchParam := identity.NewListIdentitiesParams()
	filter := fmt.Sprintf("name=\"%s\"", name)
	searchParam.Filter = &filter

	// query to see if a service already exists
	list, err := m.edge.Load().Identity.ListIdentities(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service: %w", err)
	}

	// if the identity doesn't exist, return
	if list == nil || len(list.Payload.Data) == 0 {
		return nil
	}

	// delete the identity
	params := &identity.DeleteIdentityParams{
		Context: context.Background(),
		ID:      *list.Payload.Data[0].ID,
	}
	params.SetTimeout(10 * time.Second)
	_, err = m.edge.Load().Identity.DeleteIdentity(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete identity: %w", err)
	}

	return nil
}

// CreateServer
//
// Creates a server in the ziti mesh
func (m *Manager) CreateServer(id int64) (string, string, error) {
	// create our variables for the identity of the server
	isAdmin := false
	name := fmt.Sprintf("gigo-server-%d", id)
	identityType := rest_model_edge.IdentityTypeDevice

	// create filter to search for service
	searchParam := identity.NewListIdentitiesParams()
	filter := fmt.Sprintf("name=\"%s\"", name)
	searchParam.Filter = &filter

	// query to see if a service already exists
	list, err := m.edge.Load().Identity.ListIdentities(searchParam, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to query for service: %w", err)
	}

	// create a variable to hold the id of the identity
	var identityId string

	// if the identity exists then we need to fail since agents can't be recreated
	if list != nil && len(list.Payload.Data) > 0 {
		return "", "", ErrIdentityExists
	}

	// create the request for the identity
	createIdentityReq := identity.NewCreateIdentityParams()
	createIdentityReq.Identity = &rest_model_edge.IdentityCreate{
		Enrollment:          &rest_model_edge.IdentityCreateEnrollment{Ott: true},
		IsAdmin:             &isAdmin,
		Name:                &name,
		RoleAttributes:      &rest_model.Attributes{"gigo-servers"},
		ServiceHostingCosts: nil,
		Tags:                nil,
		Type:                &identityType,
	}
	createIdentityReq.SetTimeout(10 * time.Second)

	// create the identity
	createIdentityRes, err := m.edge.Load().Identity.CreateIdentity(createIdentityReq, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create identity: %w", err)
	}

	// set the identity id
	identityId = createIdentityRes.Payload.Data.ID

	// retrieve the token for the identity
	params := &identity.DetailIdentityParams{
		Context: context.Background(),
		ID:      identityId,
	}
	params.SetTimeout(10 * time.Second)
	resp, err := m.edge.Load().Identity.DetailIdentity(params, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve identity: %w", err)
	}
	return identityId, resp.Payload.Data.Enrollment.Ott.JWT, nil
}

// DeleteServer
//
// Deletes a server from the ziti mesh
func (m *Manager) DeleteServer(id int64) error {
	// create our variables for the identity of the agent
	name := fmt.Sprintf("gigo-server-%d", id)

	// create filter to search for service
	searchParam := identity.NewListIdentitiesParams()
	filter := fmt.Sprintf("name=\"%s\"", name)
	searchParam.Filter = &filter

	// query to see if a service already exists
	list, err := m.edge.Load().Identity.ListIdentities(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service: %w", err)
	}

	// if the identity doesn't exist, return
	if list == nil || len(list.Payload.Data) == 0 {
		return nil
	}

	// delete the identity
	params := &identity.DeleteIdentityParams{
		Context: context.Background(),
		ID:      *list.Payload.Data[0].ID,
	}
	params.SetTimeout(10 * time.Second)
	_, err = m.edge.Load().Identity.DeleteIdentity(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete identity: %w", err)
	}

	return nil
}

// CreateWorkspaceService
//
// Creates a service in the ziti mesh
func (m *Manager) CreateWorkspaceService(agentId int64) (string, error) {
	// create filter to search for service
	name := fmt.Sprintf("gigo-workspace-access-%d", agentId)
	searchParam := service.NewListServicesParams()
	filter := fmt.Sprintf("name=\"%s\"", name)
	searchParam.Filter = &filter

	// query to see if a service already exists
	id, err := m.edge.Load().Service.ListServices(searchParam, nil)
	if err != nil {
		return "", fmt.Errorf("failed to query for service: %w", err)
	}
	if id != nil && len(id.Payload.Data) > 0 {
		return name, nil
	}

	// create a new service since no service exists
	encryptOn := true
	serviceCreate := &rest_model.ServiceCreate{
		EncryptionRequired: &encryptOn,
		Name:               &name,
		RoleAttributes:     rest_model.Roles{"gigo-workspace-access"},
	}
	serviceParams := &service.CreateServiceParams{
		Service: serviceCreate,
		Context: context.Background(),
	}
	serviceParams.SetTimeout(30 * time.Second)
	_, err = m.edge.Load().Service.CreateService(serviceParams, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create service: %w", err)
	}

	return name, nil
}

// DeleteWorkspaceService
//
// Deletes a service from the ziti mesh
func (m *Manager) DeleteWorkspaceService(agentId int64) error {
	// create filter to search for service
	name := fmt.Sprintf("gigo-workspace-access-%d", agentId)
	searchParam := service.NewListServicesParams()
	filter := fmt.Sprintf("name=\"%s\"", name)
	searchParam.Filter = &filter

	// query to see if a service already exists
	id, err := m.edge.Load().Service.ListServices(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service: %w", err)
	}
	if id == nil || len(id.Payload.Data) == 0 {
		return nil
	}

	// delete the service
	params := &service.DeleteServiceParams{
		Context: context.Background(),
		ID:      *id.Payload.Data[0].ID,
	}
	params.SetTimeout(10 * time.Second)
	_, err = m.edge.Load().Service.DeleteService(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

// CreateWorkspaceServicePolicy
//
// Creates a service policy in the ziti mesh
func (m *Manager) CreateWorkspaceServicePolicy() error {
	// create filter to search for the service policy enabling agents to bind to the service
	searchParam := service_policy.NewListServicePoliciesParams()
	filter := "name=\"gigo-workspace-access-bind\""
	searchParam.Filter = &filter

	// query to see if a service policy already exists
	id, err := m.edge.Load().ServicePolicy.ListServicePolicies(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service policy: %w", err)
	}
	if id == nil || len(id.Payload.Data) == 0 {
		// create a new service policy for agents to bind to the service
		name := "gigo-workspace-access-bind"
		servType := rest_model.DialBindBind
		defaultSemantic := rest_model.SemanticAnyOf
		servicePolicy := &rest_model.ServicePolicyCreate{
			IdentityRoles: rest_model.Roles{"#gigo-agents"},
			Name:          &name,
			Semantic:      &defaultSemantic,
			ServiceRoles:  rest_model.Roles{"#gigo-workspace-access"},
			Type:          &servType,
		}
		params := &service_policy.CreateServicePolicyParams{
			Policy:  servicePolicy,
			Context: context.Background(),
		}
		params.SetTimeout(30 * time.Second)
		_, err := m.edge.Load().ServicePolicy.CreateServicePolicy(params, nil)
		if err != nil {
			return fmt.Errorf("failed to create service policy: %w", err)
		}
	}

	// create filter to search for the service policy enabling agents to bind to the service
	searchParam = service_policy.NewListServicePoliciesParams()
	filter = "name=\"gigo-workspace-access-dial\""
	searchParam.Filter = &filter

	// query to see if a service policy already exists
	id, err = m.edge.Load().ServicePolicy.ListServicePolicies(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service policy: %w", err)
	}
	if id == nil || len(id.Payload.Data) == 0 {
		// create a new service policy for agents to bind to the service
		name := "gigo-workspace-access-dial"
		servType := rest_model.DialBindDial
		defaultSemantic := rest_model.SemanticAllOf
		servicePolicy := &rest_model.ServicePolicyCreate{
			IdentityRoles: rest_model.Roles{"#gigo-servers"},
			Name:          &name,
			Semantic:      &defaultSemantic,
			ServiceRoles:  rest_model.Roles{"#gigo-workspace-access"},
			Type:          &servType,
		}
		params := &service_policy.CreateServicePolicyParams{
			Policy:  servicePolicy,
			Context: context.Background(),
		}
		params.SetTimeout(30 * time.Second)
		_, err := m.edge.Load().ServicePolicy.CreateServicePolicy(params, nil)
		if err != nil {
			return fmt.Errorf("failed to create service policy: %w", err)
		}
	}

	return nil
}

// DeleteWorkspaceServicePolicy
//
// Deletes a service policy from the ziti mesh
func (m *Manager) DeleteWorkspaceServicePolicy() error {
	// create filter to search for service
	searchParam := service_policy.NewListServicePoliciesParams()
	filter := "name=\"gigo-workspace-access-bind\""
	searchParam.Filter = &filter

	// query to see if a service policy already exists
	id, err := m.edge.Load().ServicePolicy.ListServicePolicies(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service policy: %w", err)
	}
	if id != nil && len(id.Payload.Data) > 0 {
		// delete the service policy
		params := &service_policy.DeleteServicePolicyParams{
			Context: context.Background(),
			ID:      *id.Payload.Data[0].ID,
		}
		params.SetTimeout(10 * time.Second)
		_, err = m.edge.Load().ServicePolicy.DeleteServicePolicy(params, nil)
		if err != nil {
			return fmt.Errorf("failed to delete service policy: %w", err)
		}
	}

	// create filter to search for service
	searchParam = service_policy.NewListServicePoliciesParams()
	filter = "name=\"gigo-workspace-access-dial\""
	searchParam.Filter = &filter

	// query to see if a service policy already exists
	id, err = m.edge.Load().ServicePolicy.ListServicePolicies(searchParam, nil)
	if err != nil {
		return fmt.Errorf("failed to query for service policy: %w", err)
	}
	if id != nil && len(id.Payload.Data) > 0 {
		// delete the service policy
		params := &service_policy.DeleteServicePolicyParams{
			Context: context.Background(),
			ID:      *id.Payload.Data[0].ID,
		}
		params.SetTimeout(10 * time.Second)
		_, err = m.edge.Load().ServicePolicy.DeleteServicePolicy(params, nil)
		if err != nil {
			return fmt.Errorf("failed to delete service policy: %w", err)
		}
	}

	return nil
}
