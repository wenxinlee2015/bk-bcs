/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package listenercontroller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	networkextensionv1 "github.com/Tencent/bk-bcs/bcs-k8s/kubernetes/apis/networkextension/v1"
	"github.com/Tencent/bk-bcs/bcs-network/bcs-ingress-controller/internal/cloud"
	"github.com/Tencent/bk-bcs/bcs-network/bcs-ingress-controller/internal/constant"
	"github.com/Tencent/bk-bcs/bcs-network/bcs-ingress-controller/internal/option"
	"github.com/Tencent/bk-bcs/bcs-network/bcs-ingress-controller/internal/worker"
	"github.com/Tencent/bk-bcs/bcs-network/pkg/common"
)

// getListenerPredicate filter listener events
func getListenerPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			newListener, okNew := e.ObjectNew.(*networkextensionv1.Listener)
			oldListener, okOld := e.ObjectOld.(*networkextensionv1.Listener)
			if !okNew || !okOld {
				return false
			}
			if newListener.DeletionTimestamp != nil {
				return true
			}
			if reflect.DeepEqual(newListener.Spec, oldListener.Spec) {
				blog.V(5).Infof("listener %+v updated, but spec not change", oldListener)
				return false
			}
			return true
		},
	}
}

// ListenerReconciler reconclier for networkextensionv1 listener
type ListenerReconciler struct {
	Ctx context.Context

	Client client.Client
	Option *option.ControllerOption

	ListenerEventer record.EventRecorder

	CloudLb cloud.LoadBalance

	WorkerMap map[string]*worker.EventHandler
}

// NewListenerReconciler create ListenerReconciler
func NewListenerReconciler() *ListenerReconciler {
	return &ListenerReconciler{
		WorkerMap: make(map[string]*worker.EventHandler),
	}
}

func (lc *ListenerReconciler) getListenerEventHandler(listener *networkextensionv1.Listener) (
	*worker.EventHandler, error) {
	region, ok := listener.Labels[networkextensionv1.LabelKeyForLoadbalanceRegion]
	if !ok {
		blog.Errorf("listener %s/%s lost key %s in labels",
			listener.GetNamespace(), listener.GetName(), networkextensionv1.LabelKeyForLoadbalanceRegion)
		return nil, fmt.Errorf("listener %s/%s lost key %s in labels",
			listener.GetNamespace(), listener.GetName(), networkextensionv1.LabelKeyForLoadbalanceRegion)
	}
	ehandler, ok := lc.WorkerMap[listener.Spec.LoadbalancerID]
	if !ok {
		newHandler := worker.NewEventHandler(
			region,
			listener.Spec.LoadbalancerID,
			lc.CloudLb,
			lc.Client,
		)
		go newHandler.Run()
		lc.WorkerMap[listener.Spec.LoadbalancerID] = newHandler
		ehandler = newHandler
	}
	return ehandler, nil
}

// Reconcile reconclie listener
func (lc *ListenerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	blog.V(2).Infof("listener %+v triggered", req.NamespacedName)

	listener := &networkextensionv1.Listener{}
	err := lc.Client.Get(context.TODO(), req.NamespacedName, listener)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		blog.Errorf("get listener %s/%s failed, err %s", req.Namespace, req.Name, err.Error())
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: 5 * time.Second,
		}, nil
	}

	ehandler, err := lc.getListenerEventHandler(listener)
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: 5 * time.Second,
		}, nil
	}
	if listener.DeletionTimestamp != nil {
		listener.Finalizers = common.RemoveString(listener.Finalizers, constant.FinalizerNameBcsIngressController)
		err := lc.Client.Update(context.TODO(), listener, &client.UpdateOptions{})
		if err != nil {
			blog.Errorf("failed to remove finalizer from listener %s/%s, err %s", req.Namespace, req.Name, err.Error())
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: 5 * time.Second,
			}, nil
		}
		ehandler.PushEvent(worker.NewListenerEvent(
			worker.EventDelete,
			listener.GetName(),
			listener.GetNamespace(), *listener))
	} else {
		ehandler.PushEvent(worker.NewListenerEvent(
			worker.EventUpdate,
			listener.GetName(),
			listener.GetNamespace(), *listener))
	}
	return ctrl.Result{}, nil
}

// SetupWithManager set reconciler
func (lc *ListenerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkextensionv1.Listener{}).
		WithEventFilter(getListenerPredicate()).
		Complete(lc)
}
