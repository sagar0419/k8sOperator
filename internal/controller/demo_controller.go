/*
Copyright 2023.

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

package controller

import (
	"context"

	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"

	// apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8soperatorv1 "k8sOperator/api/v1"
)

// DemoReconciler reconciles a Demo object
type DemoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8soperator.sagar.com,resources=demoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8soperator.sagar.com,resources=demoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8soperator.sagar.com,resources=demoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Demo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *DemoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)

	log := r.Log.WithValues("Demo", req.NamespacedName)
	// TODO(user): your logic here

	// Here k8soperatorv1 is api group and Demo is custom resource
	// Initialize a pointer variable "operator" to an empty "Demo" structure.
	operator := &k8soperatorv1.Demo{}
	// Using get method to populate it with the data retrived from cluster.
	// If the resource is not found it will store the error in "err" variable
	err := r.Client.Get(ctx, req.NamespacedName, operator)
	// checking if there is an error while fetching the details
	if err != nil {
		// checking if error is resource "Is not found"
		if errors.IsNotFound(err) {
			log.Info("Operator is not found. Ignoring since object is deleted")
			// since the error is not found it means it has been deleted so now we will not redeploy it. We will stop the cntroller from redeploying the app.
			// and we will terminate the program here only
			return ctrl.Result{}, nil
		}
		// If error is other than not found, that means crashloopback, Imagepull back off than we need to retry the deployment of the app.
		// Here we are logging the error.
		log.Error(err, "Failed to get opertaor")

		// when we return the err than the reconciler will get to know that there is an issue with the request and the controller shoul  requeue the request to be processed again.
		// In short the reconciler will retry to deploy the application
		return ctrl.Result{}, err
	}

	// Initialize a pointer variable "found" to an empty "Deployment" structure.
	found := &appsv1.Deployment{}
	// Using get method to populate it with the data retrived from cluster.
	err = r.Get(ctx, types.NamespacedName{Name: operator.Name, Namespace: operator.Namespace}, found)
	// checking if there is an error while getting the details and checking if error is resource "Is not found"
	if err != nil && errors.IsNotFound(err) {
		// Calling deploymentForOperator method to create and return the deployment object and storing it in "dep" variable.
		dep := r.deploymentForOperator(operator)

		// Getting log info  of creating new deployment
		log.Info("Creating new deployment")
		// Creating new deployment if it doent exist.
		err = r.Create(ctx, dep)
		//  if error occur or deployment is uncessful while creating the deployment
		if err != nil {
			// Sending the logs of the error
			log.Error(err, "Failed to create the deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			// returning the error that has been occured while deploying.
			return ctrl.Result{}, err
		}
		// return the result which will tell the operator that we want to requeue the request.
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil { //checking if error is not a "IsNotFound".
		// Log the error
		log.Error(err, "failed to get deployment")
		// return the error with empty reult object.
		return ctrl.Result{}, err
	}

	// Updating Deployment
	// Creating deployment object and assigning the value to "deploy" variable
	deploy := r.deploymentForOperator(operator)
	// if the template do not match update the object with the template from the dpeloy object

	// deploy.Spec.Template is desired pod spec
	// found.Spec.Template is observed pod spec template
	if !equality.Semantic.DeepDerivative(deploy.Spec.Template, found.Spec.Template) {
		// update the found object with the deploy object.
		found = deploy
		// adding log info that we are updating deployment
		log.Info("Updating deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		// updating the found object
		err := r.Update(ctx, found)
		// if there is error in updating
		if err != nil {
			// log the error if deployment is unsucessful
			log.Error(err, "Failed to update the deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			// return error which we get while updating
			return ctrl.Result{}, err
		}
		// else requeue the object.
		return ctrl.Result{Requeue: true}, nil
	}

	// checking replicas of deployment is matching
	// getting  replicas mentioned in deployment replicas
	size := operator.Spec.Size

	// comparing with the deployed replicas
	if *found.Spec.Replicas != size {
		// if not equal replica is udpated with the desired size
		found.Spec.Replicas = &size
		// deployment is updated
		err = r.Update(ctx, found)
		// if there is an error in deployment
		if err != nil {
			// logging the error
			log.Error(err, "Failed to update the replicas", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// if the update is succesfull operation will return requeue true to trigger function again.
		// indicating that there are still changes to reconcile.
		return ctrl.Result{Requeue: true}, nil
	}

	// Verifying service
	foundService := &corev1.Service{}
	// getting service using r.get func using the provided types.namespacedName that represent the name and namespace of the service
	err = r.Get(ctx, types.NamespacedName{Name: operator.Name, Namespace: operator.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		dep := r.serviceForOperator(operator)
		log.Info("Creating Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to Create the Servcie", "Servcie.Namespace", dep.Namespace, "Servcie.Name", dep.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "failed to get service")
		return ctrl.Result{}, err
	}

	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(found.Namespace),
		client.MatchingLabels(map[string]string{"app": found.Name, "Labels": found.Name}),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)
	if !reflect.DeepEqual(podNames, operator.Status.PodList) {
		operator.Status.PodList = podNames
		err := r.Status().Update(ctx, operator)
		if err != nil {
			log.Error(err, "Failed to update the pod list status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// Func to get pod name
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// func to create deployment
func (r *DemoReconciler) deploymentForOperator(m *k8soperatorv1.Demo) *appsv1.Deployment {
	ls := map[string]string{
		"app":    m.Name,
		"labels": m.Name,
	}
	replicas := m.Spec.Size
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  m.Spec.AppContainerName,
							Image: m.Spec.AppImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: m.Spec.AppPort,
								},
							},
						}, {
							Name:    m.Spec.MonitorContainerName,
							Image:   m.Spec.AppImage,
							Command: []string{"sh", "-c", m.Spec.MonitorCommand},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: m.Spec.AppPort,
								},
							},
						},
					},
				},
			},
		},
	}
	// ctrl.SetControllerReference this functionis used to set the parenbt child relation between 2 object and allow kubernetes two properly manage the objects
	// m is owner object, dep is deployment object, r.scheme is the runtime.Scheme object. Runtime scheme is used to encode and decode the kubernetes object.
	// when demo will get deleted the kubernetes will delete everything which is related to it and created by operator.
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// func to create service
func (r *DemoReconciler) serviceForOperator(m *k8soperatorv1.Demo) *corev1.Service {
	ls := map[string]string{
		"app":    m.Name,
		"labels": m.Name,
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Type:     corev1.ServiceType(m.Spec.Service.Type),
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.Protocol(m.Spec.Service.Protocol),
					TargetPort: intstr.FromInt(int(m.Spec.Service.TargetPort)),
					Name:       m.Spec.Service.Name,
					Port:       m.Spec.Service.Port,
					NodePort:   m.Spec.Service.NodePort,
				},
			},
		},
	}
	// ctrl.SetControllerReference this functionis used to set the parenbt child relation between 2 object and allow kubernetes two properly manage the objects
	// m is owner object, dep is deployment object, r.scheme is the runtime.Scheme object. Runtime scheme is used to encode and decode the kubernetes object.
	// when demo will get deleted the kubernetes will delete everything which is related to it and created by operator.
	ctrl.SetControllerReference(m, svc, r.Scheme)
	return svc
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8soperatorv1.Demo{}).
		Complete(r)
}
