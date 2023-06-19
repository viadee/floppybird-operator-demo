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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1alpha1 "github.com/viadee/floppybird-operator-demo/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FloppybirdReconciler reconciles a Floppybird object
type FloppybirdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.demo.viadee.de,resources=floppybirds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.demo.viadee.de,resources=floppybirds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.demo.viadee.de,resources=floppybirds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Floppybird object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *FloppybirdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here

	// get current state
	var floppybird webappv1alpha1.Floppybird
	if err := r.Client.Get(ctx, req.NamespacedName, &floppybird); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Floppybird: " + floppybird.Name)

	// define resource label
	operatorResourceLabel := map[string]string{"floppybird-instance": floppybird.ObjectMeta.Name}

	// get number of running pods
	runningPods, err := r.getListOfPods(ctx, req, operatorResourceLabel)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	// start new pod if no running pods are found
	if runningPods <= 0 {
		logger.Info("No running pods found. Creating new pod.")
		newPod := r.createPod(floppybird, operatorResourceLabel)

		// establish a controller reference between floppybird and pod
		if err := ctrl.SetControllerReference(&floppybird, newPod, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		// create in cluster
		if err := r.Create(ctx, newPod); err != nil {
			// Error occurred while creating a new pod, return error
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Running pods found: " + string(runningPods))
	}

	// check if servie exsists
	service := &corev1.Service{}
	err = r.Client.Get(ctx, req.NamespacedName, service)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	// create new service if service does not exit
	if service.Name == "" {
		// service does not exist, create new service
		logger.Info("No service found. Creating new service.")
		service := createService(req, operatorResourceLabel)

		// establish a controller reference between floppybird and service
		if err := ctrl.SetControllerReference(&floppybird, service, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		// create in cluster
		if err := r.Create(ctx, service); err != nil {
			// Error occurred while creating a new pod, return error
			return ctrl.Result{}, err
		}
	}

	// check if an ingress exist
	ingress := &networkingv1.Ingress{}
	err = r.Get(ctx, req.NamespacedName, ingress)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}
	// create new ingress if service does not exit
	if ingress.Name == "" {
		logger.Info("No ingress found. Creating new ingress.")
		ingress := createIngress(req)

		// establish a controller reference between floppybird and service
		if err := ctrl.SetControllerReference(&floppybird, ingress, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		// create in cluster
		if err := r.Create(ctx, ingress); err != nil {
			// Error occurred while creating a new pod, return error
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *FloppybirdReconciler) createPod(floppybird webappv1alpha1.Floppybird, labels map[string]string) *corev1.Pod {
	// Customize the pod creation according to your requirements
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      floppybird.Name + "-pod",
			Namespace: floppybird.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:            floppybird.Name,
				Image:           "crowdsalat/floppybird-demo:0.0.4",
				ImagePullPolicy: "IfNotPresent",
				Ports: []corev1.ContainerPort{
					{ContainerPort: 8000, Name: "http", Protocol: "TCP"},
				},
			}},
		},
	}
	return newPod
}

func createService(req ctrl.Request, label map[string]string) *corev1.Service {
	// Customize the service creation according to your requirements
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: label,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8000),
				},
			},
		},
	}
	return service
}

func createIngress(req ctrl.Request) *networkingv1.Ingress {
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "floppybird.cloudland2023operator.viadee.cloud",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: req.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			TLS: []networkingv1.IngressTLS{
				{
					Hosts: []string{
						"floppybird.cloudland2023operator.viadee.cloud",
					},
					SecretName: "floppybird-cert-generatedby-cert-manager",
				},
			},
		},
	}
	return ingress
}

func (r *FloppybirdReconciler) getListOfPods(ctx context.Context, req ctrl.Request, operatorResourceLabel map[string]string) (int32, error) {
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(req.Namespace), client.MatchingLabels(operatorResourceLabel)); err != nil {
		return 0, err
	}
	runningPods := int32(len(podList.Items))
	return runningPods, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FloppybirdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.Floppybird{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
