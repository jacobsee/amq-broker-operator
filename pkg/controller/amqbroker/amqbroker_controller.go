package amqbroker

import (
	"context"

	jacobseev1alpha1 "github.com/jacobsee/amq-broker-operator/pkg/apis/jacobsee/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_amqbroker")

// Add creates a new AMQBroker Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAMQBroker{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("amqbroker-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AMQBroker
	err = c.Watch(&source.Kind{Type: &jacobseev1alpha1.AMQBroker{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AMQBroker
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jacobseev1alpha1.AMQBroker{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAMQBroker{}

// ReconcileAMQBroker reconciles a AMQBroker object
type ReconcileAMQBroker struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AMQBroker object and makes changes based on the state read
// and what is in the AMQBroker.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAMQBroker) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AMQBroker")

	// Fetch the AMQBroker instance
	instance := &jacobseev1alpha1.AMQBroker{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Deployment object
	deployment := newDeploymentForCR(instance)
	service := newServiceForCR(instance)
	route := newRouteForCR(instance)

	// Set AMQBroker instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this deployment & service already exist
	deploymentCreated := createDeploymentIfNotExists(deployment, r)
	serviceCreated := createServiceIfNotExists(service, r)
	routeCreated := createRouteIfNotExists(route, r)

	if deploymentCreated != nil {
		return reconcile.Result{}, deploymentCreated
	}
	if serviceCreated != nil {
		return reconcile.Result{}, serviceCreated
	}
	if routeCreated != nil {
		return reconcile.Result{}, routeCreated
	}
	return reconcile.Result{}, nil
}

func createDeploymentIfNotExists(deployment *appsv1.Deployment, r *ReconcileAMQBroker) (err error) {
	reqLogger := log.WithValues("Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return err
		}
		// Deployment created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return nil
}

func createServiceIfNotExists(service *corev1.Service, r *ReconcileAMQBroker) (err error) {
	reqLogger := log.WithValues("Service.Namespace", service.Namespace, "Service.Name", service.Name)
	found := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return err
		}
		// Service created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}
	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	return nil
}

func createRouteIfNotExists(route *routev1.Route, r *ReconcileAMQBroker) (err error) {
	reqLogger := log.WithValues("Route.Namespace", route.Namespace, "Route.Name", route.Name)
	found := &routev1.Route{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: route.Name, Namespace: route.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = r.client.Create(context.TODO(), route)
		if err != nil {
			return err
		}
		// Route created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}
	reqLogger.Info("Skip reconcile: Route already exists", "Route.Namespace", found.Namespace, "Route.Name", found.Name)
	return nil
}

// newDeploymentForCR returns an AMQ Broker deployment with the same name/namespace as the cr
func newDeploymentForCR(cr *jacobseev1alpha1.AMQBroker) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}
	numReplicas := int32(1)
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-deployment",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "amq-broker-72-openshift",
						Image: "registry.access.redhat.com/amq-broker-7/amq-broker-72-openshift",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8161,
							}, {
								ContainerPort: 5672,
							}, {
								ContainerPort: 1883,
							}, {
								ContainerPort: 61613,
							}, {
								ContainerPort: 61616,
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "AMQ_USER",
								Value: cr.Spec.Username,
							},
							{
								Name:  "AMQ_PASSWORD",
								Value: cr.Spec.Password,
							},
						},
					}},
				},
			},
		},
	}
	return &dep
}

// newServiceForCR returns a service pointing to the AMQ Broker deployment with the same name/namespace as the cr
func newServiceForCR(cr *jacobseev1alpha1.AMQBroker) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			// Selector: appsv1.Route,
			Ports: []corev1.ServicePort{
				{
					Name:     "jolokia",
					Protocol: "TCP",
					Port:     8161,
				}, {
					Name:     "amqp",
					Protocol: "TCP",
					Port:     5672,
				}, {
					Name:     "mqtt",
					Protocol: "TCP",
					Port:     1883,
				}, {
					Name:     "stomp",
					Protocol: "TCP",
					Port:     61613,
				}, {
					Name:     "openwire",
					Protocol: "TCP",
					Port:     61616,
				},
			},
		},
	}
	return &service
}

func newRouteForCR(cr *jacobseev1alpha1.AMQBroker) *routev1.Route {
	labels := map[string]string{
		"app": cr.Name,
	}
	route := routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-route",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: cr.Name + "-service",
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(8161),
			},
		},
	}
	return &route
}
