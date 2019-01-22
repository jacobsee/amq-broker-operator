package controller

import (
	"github.com/jacobsee/amq-broker-operator/pkg/controller/amqbroker"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, amqbroker.Add)
}
