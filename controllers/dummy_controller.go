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
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	anyninesv1 "github.com/mmertdogann/dummy-operator/api/v1"
)

// DummyReconciler reconciles a Dummy object
type DummyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=anynines.interview.com,resources=dummies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=anynines.interview.com,resources=dummies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=anynines.interview.com,resources=dummies/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)
	lg.Info("Enter Reconcile", "req", req)

	// Dummy Custom Resource
	dummy := &anyninesv1.Dummy{}

	// Try to catch the Custom Resource
	getCustomResourceErr := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, dummy)
	if getCustomResourceErr != nil {
		// If the object is under deletion, many reconciliation will happen
		// And try to catch Custom Resource many times
		// Because of that ignore the not found errors
		if errors.IsNotFound(getCustomResourceErr) {
			lg.Info("Dummy resource not found. Ignoring since the object must be deleted")
			return ctrl.Result{}, nil
		}
		lg.Error(getCustomResourceErr, "Unable to fetch Dummy resource")
		return ctrl.Result{}, getCustomResourceErr
	}

	// Log Custom Resource's information
	lg.Info(getCustomResourceInfo(dummy))

	// Copy the message value from the spec to the specEcho value in the status
	if dummy.Status.SpecEcho != dummy.Spec.Message {
		dummy.Status.SpecEcho = dummy.Spec.Message
		updateStatusErr := r.Status().Update(ctx, dummy)
		if updateStatusErr != nil {
			lg.Error(updateStatusErr, "Failed to update Dummy status")
			return ctrl.Result{}, updateStatusErr
		}
	}

	// Check if nginx pod exists, otherwise create a pod
	result, checkPodErr := r.checkNginxPod(ctx, lg, dummy, r.createNginxPod(dummy))
	if result != nil {
		return *result, checkPodErr
	}

	return ctrl.Result{}, nil
}

// checkNginxPod checks if nginx pod exists, otherwise creates nginx pod and give it an ownership
func (r *DummyReconciler) checkNginxPod(ctx context.Context,
	lg logr.Logger,
	instance *anyninesv1.Dummy,
	pod *v1.Pod,
) (*reconcile.Result, error) {
	found := &v1.Pod{}

	// Get the pod with the provided information
	getPodErr := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      pod.Name,
		Namespace: pod.Namespace,
	}, found)

	// If pod does not exist
	if getPodErr != nil && errors.IsNotFound(getPodErr) {

		lg.Info("Creating a new Pod",
			"Pod.Namespace", pod.Namespace,
			"Pod.Name", pod.Name)
		createPodErr := r.Client.Create(context.TODO(), pod)
		if createPodErr != nil {
			lg.Error(createPodErr, "Failed to create new Pod",
				"Pod.Namespace", pod.Namespace,
				"Pod.Name", pod.Name)
			return &reconcile.Result{}, createPodErr
		}

		// If the pod has the phase pending, update the CR's status
		if pod.Status.Phase == "Pending" {
			instance.Status.PodStatus = "Pending"

			updateStatusErr := r.Status().Update(ctx, instance)
			if updateStatusErr != nil {
				lg.Error(updateStatusErr, "Failed to update Dummy status")
				return &reconcile.Result{}, updateStatusErr
			}
		}

	} else if getPodErr != nil {
		lg.Error(getPodErr, "Failed to get Pod")
		return &ctrl.Result{}, getPodErr
	}

	// If the pod has the phase running, update the CR's status
	if found.Status.Phase == "Running" {
		instance.Status.PodStatus = "Running"
	}

	// If the pod has is under deletion, update the CR's status
	if found.DeletionTimestamp != nil {
		instance.Status.PodStatus = "Removing"
	}

	updateStatusErr := r.Status().Update(ctx, instance)
	if updateStatusErr != nil {
		lg.Error(updateStatusErr, "Failed to update Dummy status")
		return &reconcile.Result{}, updateStatusErr
	}

	return &ctrl.Result{}, nil
}

// createNginxPod creates nginx pod with an owner reference for the controller
func (r *DummyReconciler) createNginxPod(cr *anyninesv1.Dummy) *v1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.Name,
			Namespace:       cr.Namespace,
			Labels:          labels,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(cr, anyninesv1.GroupVersion.WithKind("Dummy"))},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  cr.Name,
					Image: "nginx:alpine",
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							Protocol:      v1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
}

// getCustomResourceInfo fetch information of the specific Custom Resource.
func getCustomResourceInfo(cr *anyninesv1.Dummy) string {
	return fmt.Sprintf(
		"\n---------- Dummy CR Info ----------\n"+
			"Name: %s\n"+
			"Namespace: %s\n"+
			"Message: %s\n"+
			"PodStatus: %s\n",
		cr.Name, cr.Namespace, cr.Spec.Message, cr.Status.PodStatus)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&anyninesv1.Dummy{}).
		Owns(&v1.Pod{}).
		Complete(r)
}
