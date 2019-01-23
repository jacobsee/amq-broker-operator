# Operator for Red Hat AMQ Broker 7.2 on OpenShift

**THIS IS NOT READY FOR REAL USE IN ANY WAY, SHAPE, OR FORM.**

## Compiling

```
operator-sdk build docker.io/jacobsee/amq-operator
docker push docker.io/jacobsee/amq-operator
```

(if you're not me, then... use a different repository. Remember to change `deploy/operator.yaml` if you do.)

## Deploying the Operator

You're going to need `cluster-admin` for this. I'm not happy about that either.

```
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role.yaml
oc apply -f deploy/role_binding.yaml
oc apply -f deploy/crds/jacobsee_v1alpha1_amqbroker_crd.yaml
oc apply -f deploy/operator.yaml
```

## Using the Operator

Get yourself a nice `AMQBroker` resource, like this one

```
apiVersion: jacobsee.com/v1alpha1
kind: AMQBroker
metadata:
  name: example-amqbroker
spec:
  username: test_user
  password: test_pass
```

and deploy it (it's included as an example).

```
oc apply -f deploy/crds/jacobsee_v1alpha1_amqbroker_cr.yaml
```

You should see a deployment `example-amqbroker-deployment`, a service `example-amqbroker-service`, and a route `example-amqbroker-route`.

If you access the route, it should take you to the AMQ management console (after a few minutes of booting up... it's not instantaneous).

## But this doesn't have persistence/TLS/other good stuff

You're right, I haven't added those things, you should give it a try. I just did this to quickly learn about operators.

## There are a lot of files here. Where do I start?

`pkg/apis/jacobsee/v1alpha1/amqbroker_types.go`

`pkg/controller/amqbroker/amqbroker_controller.go`