package controller

import (
	"github.com/niqdev/lastpass-operator/pkg/controller/lastpass"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, lastpass.Add)
}
