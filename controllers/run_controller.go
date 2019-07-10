/*
Copyright 2019 microsoft.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	dbazure "github.com/xinsnake/databricks-sdk-golang/azure"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	databricksv1 "github.com/microsoft/azure-databricks-operator/api/v1"
)

// RunReconciler reconciles a Run object
type RunReconciler struct {
	client.Client
	Log logr.Logger

	Recorder  record.EventRecorder
	APIClient dbazure.DBClient
}

// +kubebuilder:rbac:groups=databricks.microsoft.com,resources=runs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=databricks.microsoft.com,resources=runs/status,verbs=get;update;patch

func (r *RunReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("run", req.NamespacedName)

	// Fetch the Run instance
	instance := &databricksv1.Run{}
	err := r.Get(context.Background(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if instance.IsBeingDeleted() {
		err := r.handleFinalizer(instance)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("error when handling finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !instance.HasFinalizer(databricksv1.RunFinalizerName) {
		err = r.addFinalizer(instance)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when removing finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !instance.IsSubmitted() {
		err = r.submitRunToDatabricks(instance)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when submitting run to API: %v", err)
		}
	}

	if instance.IsSubmitted() {
		err = r.refreshDatabricksRun(instance)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when refreshing run to API: %v", err)
		}
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *RunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databricksv1.Run{}).
		Complete(r)
}
