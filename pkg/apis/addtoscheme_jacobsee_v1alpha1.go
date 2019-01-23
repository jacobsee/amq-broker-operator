package apis

import (
	"github.com/jacobsee/amq-broker-operator/pkg/apis/jacobsee/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme), routev1.Install)
}
