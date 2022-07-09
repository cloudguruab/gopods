/*
Copyright 2022.

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
	"sort"

	osp "github.com/gophercloud/gophercloud"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	corev1 "k8s.io/api/core/v1"

	manilav1 "github.com/cloudguruab/gopods/api/v1"
)

type realClock struct{}

func (_ realClock) Now() time.Time { return time.Now() }

// clock knows how to get the current time.
// It can be used to fake out timing for testing.
type Clock interface {
    Now() time.Time
}

// FormatReconciler reconciles a Format object
type FormatReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Clock
}

//+kubebuilder:rbac:groups=manila.example.com,resources=formats,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manila.example.com,resources=formats/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manila.example.com,resources=formats/finalizers,verbs=update
//+kubebuilder:rbac:groups=manila,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manila,resources=jobs/status,verbs=get
func (r *FormatReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var format manilav1.Format
	if err := r.Get(ctx, req.NamespacedName, &format); err != nil {
		log.Error(err, "unable to fetch formatted Share")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var shares osp.ListShares
	if err := r.List(ctx, &shares, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list Shares")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FormatReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&manilav1.Format{}).
		Complete(r)
}
