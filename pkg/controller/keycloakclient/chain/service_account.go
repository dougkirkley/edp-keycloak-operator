package chain

import (
	v1v1alpha1 "github.com/epmd-edp/keycloak-operator/pkg/apis/v1/v1alpha1"
	"github.com/pkg/errors"
)

type ServiceAccount struct {
	BaseElement
	next Element
}

func (el *ServiceAccount) Serve(keycloakClient *v1v1alpha1.KeycloakClient) error {
	if keycloakClient.Spec.ServiceAccount == nil || !keycloakClient.Spec.ServiceAccount.Enabled {
		return el.NextServeOrNil(el.next, keycloakClient)
	}

	if keycloakClient.Spec.ServiceAccount != nil && keycloakClient.Spec.Public {
		return errors.New("service account can not be configured with public client")
	}

	clientRoles := make(map[string][]string)
	for _, v := range keycloakClient.Spec.ServiceAccount.ClientRoles {
		clientRoles[v.ClientID] = v.Roles
	}

	if err := el.State.AdapterClient.SyncServiceAccountRoles(keycloakClient.Spec.TargetRealm,
		keycloakClient.Status.ClientID, keycloakClient.Spec.ServiceAccount.RealmRoles, clientRoles); err != nil {
		return errors.Wrap(err, "unable to sync service account roles")
	}

	return el.NextServeOrNil(el.next, keycloakClient)
}