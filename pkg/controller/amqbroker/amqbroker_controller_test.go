package amqbroker

import (
	"context"
	"testing"

	jacobseev1alpha1 "github.com/jacobsee/amq-broker-operator/pkg/apis/jacobsee/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TestAmqbrokerController tests the AMQ Broker Controller... obviously.
func TestAmqbrokerController(t *testing.T) {

	// We're going to reuse these a few times throughout this test.
	var (
		name      = "amqbroker"
		namespace = "amqbrokernamespace"
		username  = "test_user"
		password  = "test_pass"
	)

	// This is the custom resource that we expect the operator to pick up and do some things with.
	brokerResource := &jacobseev1alpha1.AMQBroker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: jacobseev1alpha1.AMQBrokerSpec{
			Username: username,
			Password: password,
		},
	}

	// The faked kube API will choke without v1.Route added to the API, since the operator automatically adds a route.
	routeResource := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	// We want the fake "cluster" to be pre-seeded with these.
	objs := []runtime.Object{
		brokerResource,
	}

	// But the API needs to know about all of these.
	scheme := scheme.Scheme
	scheme.AddKnownTypes(jacobseev1alpha1.SchemeGroupVersion, brokerResource)
	scheme.AddKnownTypes(routev1.SchemeGroupVersion, routeResource)

	client := fake.NewFakeClient(objs...)

	// Actually create the thing that can do the reconciling.
	reconciler := &ReconcileAMQBroker{client: client, scheme: scheme}

	// Construct the reconcile event...
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	// ...and send it.
	res, err := reconciler.Reconcile(request)

	// Did the operator behave?
	if err != nil {
		t.Fatalf("reconcile error: (%v)", err)
	}
	if res.Requeue {
		t.Error("reconcile requeued unexpectedly")
	}

	// And did the deployment actually show up in the API?
	dep := &appsv1.Deployment{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "amqbroker-deployment", Namespace: namespace}, dep)
	if err != nil {
		t.Fatalf("failure getting deployment: (%v)", err)
	}

}
