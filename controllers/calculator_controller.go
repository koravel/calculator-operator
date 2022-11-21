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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/calculator-operator/api/v1alpha1"
)

func labelsForCalculator(name string) map[string]string {
	var imageTag string
	image, err := imageForCalculator()
	if err == nil {
		imageTag = strings.Split(image, ":")[1]
	}
	return map[string]string{"app.kubernetes.io/name": "Calculator",
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/version":    imageTag,
		"app.kubernetes.io/part-of":    "calculator-operator",
		"app.kubernetes.io/created-by": "controller-manager",
	}
}

// TODO: imageForcalculator gets the Operand image which is managed by this controller from the calculator_IMAGE environment variable defined in the config/manager/manager.yaml
func imageForCalculator() (string, error) {
	var imageEnvVar = "CALCULATOR_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("Unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}

func (r *CalculatorReconciler) deploymentForCalculator(
	calculator *v1alpha1.Calculator) (*appsv1.Deployment, error) {
	ls := labelsForCalculator(calculator.Name)
	replicas := calculator.Spec.Size

	// Get the Operand image
	image, err := imageForCalculator()
	if err != nil {
		return nil, err
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      calculator.Name,
			Namespace: calculator.Namespace,
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
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &[]bool{true}[0],
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "calculator",
						ImagePullPolicy: corev1.PullIfNotPresent,
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot:             &[]bool{true}[0],
							RunAsUser:                &[]int64{1001}[0],
							AllowPrivilegeEscalation: &[]bool{false}[0],
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{
									"ALL",
								},
							},
						},
						//TODO
						//Ports: []corev1.ContainerPort{{
						//	ContainerPort: calculator.Spec.ContainerPort,
						//	Name:          "calculator",
						//}},
						//TODO
						//Command: []string{"calculator", "-m=64", "-o", "modern", "-v"},
					}},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(calculator, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// CalculatorReconciler reconciles a Calculator object
type CalculatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=calculators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=calculators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=calculators/finalizers,verbs=update

func (r *CalculatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	calculator := &v1alpha1.Calculator{}
	err := r.Get(ctx, req.NamespacedName, calculator)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("calculator resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get calculator")
		return ctrl.Result{}, err
	}

	calculator.Status.Sum = calculator.Spec.X + calculator.Spec.Y
	calculator.Status.Processed = true

	if err = r.Status().Update(ctx, calculator); err != nil {
		logger.Error(err, "Failed to update Calculator status")
		return ctrl.Result{}, err
	}

	if err = r.Get(ctx, req.NamespacedName, calculator); err != nil {
		logger.Error(err, "Failed to re-fetch calculator")
		calculator.Status.Processed = true
		if err = r.Status().Update(ctx, calculator); err != nil {
			logger.Error(err, "Failed to update Calculator status during re-fetch error")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: calculator.Name, Namespace: calculator.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		//dep, err := r.deploymentForCalculator(calculator)
		logger.Error(err, "No deployment found")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CalculatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//TODO
	//return ctrl.NewControllerManagedBy(mgr).
	//	For(&v1alpha1.Calculator{}).
	//	Owns(&appsv1.Deployment{}).
	//	Complete(r)
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Calculator{}).
		Complete(r)
}
