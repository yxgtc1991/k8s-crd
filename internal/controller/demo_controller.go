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

	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demosv1 "james/api/v1"
)

// DemoReconciler reconciles a Demo object
type DemoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demos.james.com,resources=demoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demos.james.com,resources=demoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demos.james.com,resources=demoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Demo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *DemoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1、获取 demo 的 api 实例对象
	demo := &demosv1.Demo{}
	err := r.Get(ctx, req.NamespacedName, demo)
	if err != nil {
		logger.Error(err, "get demo failed")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2、获取当前 api 对象关联的运行 pod
	podList := &v1.PodList{}
	err = r.List(ctx, podList, client.InNamespace(req.Namespace), client.MatchingLabels{"app": demo.Name})
	if err != nil {
		logger.Error(err, "get pod list failed")
		return ctrl.Result{}, err
	}
	currentPodCount := len(podList.Items)

	// 3、判断新建还是删除 Pod

	// 期望状态的 Pod 数量 > 当前启动的 Pod 数量，则需要新建 Pod
	if int(demo.Spec.Replicas) > currentPodCount {
		logger.Info("Pod 数不足，触发新建\n")
		for i := 0; i < int(demo.Spec.Replicas)-currentPodCount; i++ {
			logger.Info("", "新建", i)
			// 生成需要创建的 Pod 对象
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      demo.Name + "-" + uuid.New().String()[0:8],
					Namespace: demo.Namespace,
					Labels:    demo.Labels,
				},
				Spec: demo.Spec.Template.Spec,
			}

			// 建立关联
			err = ctrl.SetControllerReference(demo, pod, r.Scheme)
			if err != nil {
				logger.Error(err, "unable to set ownerReference for pod")
				return ctrl.Result{}, err
			}

			// 创建 Pod
			err = r.Create(ctx, pod)
			if err != nil {
				logger.Error(err, "unable to create pod for demo")
				return ctrl.Result{}, err
			}
		}
	}

	// 期望状态的 Pod 数量 < 当前启动的 Pod 数量，则需要删除 Pod
	if int(demo.Spec.Replicas) < currentPodCount {
		logger.Info("Pod 数超量，触发删除\n")
		deletePods := podList.Items[:currentPodCount-int(demo.Spec.Replicas)]
		for i, pod := range deletePods {
			logger.Info("", "删除", i)
			err = r.Delete(ctx, &pod)
			if err != nil {
				logger.Error(err, "unable to delete pod for demo")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demosv1.Demo{}).
		Complete(r)
}
