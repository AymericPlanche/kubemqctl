package operator

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type Deployment struct {
	Name      string
	Namespace string
	CRDs      []*v1beta1.CustomResourceDefinition
	*appsv1.Deployment
	*rbac.Role
	*rbac.RoleBinding
	OperatorServiceAccount *apiv1.ServiceAccount
	ClusterServiceAccount  *apiv1.ServiceAccount
}

func CreateDeployment(name, namespace string) (*Deployment, error) {
	kubemqClusterCrd, err := CreateKubemqClustersCRD(namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, kubemq cluster crd error: %s", err.Error())
	}
	kubemqDashboardCrd, err := CreateKubemqDashboardCRD(namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, kubemq dashboard crd error: %s", err.Error())
	}

	role, err := CreateRole(name, namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, role error: %s", err.Error())
	}

	roleBinding, err := CreateRoleBinding(name, namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, role binding error: %s", err.Error())
	}

	serviceAccount, err := CreateServiceAccount(name, namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, service account error: %s", err.Error())
	}
	clusterServiceAccount, err := CreateServiceAccount("kubemq-cluster", namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, service account error: %s", err.Error())
	}

	operator, err := CreateOperator(name, namespace).Get()
	if err != nil {
		return nil, fmt.Errorf("error create operator bundle, operator deployment error: %s", err.Error())
	}

	return &Deployment{
		Name:                   name,
		Namespace:              namespace,
		CRDs:                   []*v1beta1.CustomResourceDefinition{kubemqClusterCrd, kubemqDashboardCrd},
		Deployment:             operator,
		Role:                   role,
		RoleBinding:            roleBinding,
		OperatorServiceAccount: serviceAccount,
		ClusterServiceAccount:  clusterServiceAccount,
	}, nil
}

func (b *Deployment) IsValid() error {
	if b.CRDs == nil {
		return fmt.Errorf("no crd exsits or defined")
	}
	if b.Deployment == nil {
		return fmt.Errorf("no operator deployment exsits or defined")
	}
	if b.Role == nil {
		return fmt.Errorf("no role exsits or defined")
	}

	if b.RoleBinding == nil {
		return fmt.Errorf("no role binding exsits or defined")
	}
	if b.OperatorServiceAccount == nil {
		return fmt.Errorf("no service account exsits or defined")
	}
	if b.ClusterServiceAccount == nil {
		return fmt.Errorf("no service account exsits or defined")
	}

	return nil
}
