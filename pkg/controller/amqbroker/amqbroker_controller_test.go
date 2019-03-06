package amqbroker

import (
	"testing"

	jacobseev1alpha1 "github.com/jacobsee/amq-broker-operator/pkg/apis/jacobsee/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TestAmqbrokerController tests the AMQ Broker Controller... obviously.
func TestAmqbrokerController(t *testing.T) {

	var (
		name      = "amqbroker-operator"
		namespace = "amqbroker"
		username  = "test_user"
		password  = "test_pass"
	)

	brokerResource := &jacobseev1alpha1.AMQBroker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "AMQBroker",
			Namespace: "AMQBrokerNamespace",
		},
		Spec: jacobseev1alpha1.AMQBrokerSpec{
			Username: username,
			Password: password,
		},
	}

	objs := []runtime.Object{
		brokerResource,
	}

	scheme := scheme.Scheme
	scheme.AddKnownTypes(jacobseev1alpha1.SchemeGroupVersion, brokerResource)

	client := fake.NewFakeClient(objs...)

	reconciler := &ReconcileAMQBroker{client: client, scheme: scheme}

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	res, err := reconciler.Reconcile(request)
	if err != nil {
		t.Fatalf("reconcile error: (%v)", err)
	}
	if res.Requeue {
		t.Error("reconcile requeued unexpectedly")
	}

}
