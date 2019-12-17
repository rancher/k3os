/*
Copyright 2019 Rancher Labs, Inc.

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

// Code generated by codegen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/rancher/k3os/pkg/apis/k3os.io/v1"
	clientset "github.com/rancher/k3os/pkg/generated/clientset/versioned/typed/k3os.io/v1"
	informers "github.com/rancher/k3os/pkg/generated/informers/externalversions/k3os.io/v1"
	listers "github.com/rancher/k3os/pkg/generated/listers/k3os.io/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type NodeUpgradeHandler func(string, *v1.NodeUpgrade) (*v1.NodeUpgrade, error)

type NodeUpgradeController interface {
	generic.ControllerMeta
	NodeUpgradeClient

	OnChange(ctx context.Context, name string, sync NodeUpgradeHandler)
	OnRemove(ctx context.Context, name string, sync NodeUpgradeHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() NodeUpgradeCache
}

type NodeUpgradeClient interface {
	Create(*v1.NodeUpgrade) (*v1.NodeUpgrade, error)
	Update(*v1.NodeUpgrade) (*v1.NodeUpgrade, error)
	UpdateStatus(*v1.NodeUpgrade) (*v1.NodeUpgrade, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.NodeUpgrade, error)
	List(namespace string, opts metav1.ListOptions) (*v1.NodeUpgradeList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.NodeUpgrade, err error)
}

type NodeUpgradeCache interface {
	Get(namespace, name string) (*v1.NodeUpgrade, error)
	List(namespace string, selector labels.Selector) ([]*v1.NodeUpgrade, error)

	AddIndexer(indexName string, indexer NodeUpgradeIndexer)
	GetByIndex(indexName, key string) ([]*v1.NodeUpgrade, error)
}

type NodeUpgradeIndexer func(obj *v1.NodeUpgrade) ([]string, error)

type nodeUpgradeController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.NodeUpgradesGetter
	informer          informers.NodeUpgradeInformer
	gvk               schema.GroupVersionKind
}

func NewNodeUpgradeController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.NodeUpgradesGetter, informer informers.NodeUpgradeInformer) NodeUpgradeController {
	return &nodeUpgradeController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromNodeUpgradeHandlerToHandler(sync NodeUpgradeHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.NodeUpgrade
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.NodeUpgrade))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *nodeUpgradeController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.NodeUpgrade))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateNodeUpgradeDeepCopyOnChange(client NodeUpgradeClient, obj *v1.NodeUpgrade, handler func(obj *v1.NodeUpgrade) (*v1.NodeUpgrade, error)) (*v1.NodeUpgrade, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *nodeUpgradeController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *nodeUpgradeController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *nodeUpgradeController) OnChange(ctx context.Context, name string, sync NodeUpgradeHandler) {
	c.AddGenericHandler(ctx, name, FromNodeUpgradeHandlerToHandler(sync))
}

func (c *nodeUpgradeController) OnRemove(ctx context.Context, name string, sync NodeUpgradeHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromNodeUpgradeHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *nodeUpgradeController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *nodeUpgradeController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *nodeUpgradeController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *nodeUpgradeController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *nodeUpgradeController) Cache() NodeUpgradeCache {
	return &nodeUpgradeCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *nodeUpgradeController) Create(obj *v1.NodeUpgrade) (*v1.NodeUpgrade, error) {
	return c.clientGetter.NodeUpgrades(obj.Namespace).Create(obj)
}

func (c *nodeUpgradeController) Update(obj *v1.NodeUpgrade) (*v1.NodeUpgrade, error) {
	return c.clientGetter.NodeUpgrades(obj.Namespace).Update(obj)
}

func (c *nodeUpgradeController) UpdateStatus(obj *v1.NodeUpgrade) (*v1.NodeUpgrade, error) {
	return c.clientGetter.NodeUpgrades(obj.Namespace).UpdateStatus(obj)
}

func (c *nodeUpgradeController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.NodeUpgrades(namespace).Delete(name, options)
}

func (c *nodeUpgradeController) Get(namespace, name string, options metav1.GetOptions) (*v1.NodeUpgrade, error) {
	return c.clientGetter.NodeUpgrades(namespace).Get(name, options)
}

func (c *nodeUpgradeController) List(namespace string, opts metav1.ListOptions) (*v1.NodeUpgradeList, error) {
	return c.clientGetter.NodeUpgrades(namespace).List(opts)
}

func (c *nodeUpgradeController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.NodeUpgrades(namespace).Watch(opts)
}

func (c *nodeUpgradeController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.NodeUpgrade, err error) {
	return c.clientGetter.NodeUpgrades(namespace).Patch(name, pt, data, subresources...)
}

type nodeUpgradeCache struct {
	lister  listers.NodeUpgradeLister
	indexer cache.Indexer
}

func (c *nodeUpgradeCache) Get(namespace, name string) (*v1.NodeUpgrade, error) {
	return c.lister.NodeUpgrades(namespace).Get(name)
}

func (c *nodeUpgradeCache) List(namespace string, selector labels.Selector) ([]*v1.NodeUpgrade, error) {
	return c.lister.NodeUpgrades(namespace).List(selector)
}

func (c *nodeUpgradeCache) AddIndexer(indexName string, indexer NodeUpgradeIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.NodeUpgrade))
		},
	}))
}

func (c *nodeUpgradeCache) GetByIndex(indexName, key string) (result []*v1.NodeUpgrade, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.NodeUpgrade))
	}
	return result, nil
}

type NodeUpgradeStatusHandler func(obj *v1.NodeUpgrade, status v1.NodeUpgradeStatus) (v1.NodeUpgradeStatus, error)

type NodeUpgradeGeneratingHandler func(obj *v1.NodeUpgrade, status v1.NodeUpgradeStatus) ([]runtime.Object, v1.NodeUpgradeStatus, error)

func RegisterNodeUpgradeStatusHandler(ctx context.Context, controller NodeUpgradeController, condition condition.Cond, name string, handler NodeUpgradeStatusHandler) {
	statusHandler := &nodeUpgradeStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromNodeUpgradeHandlerToHandler(statusHandler.sync))
}

func RegisterNodeUpgradeGeneratingHandler(ctx context.Context, controller NodeUpgradeController, apply apply.Apply,
	condition condition.Cond, name string, handler NodeUpgradeGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &nodeUpgradeGeneratingHandler{
		NodeUpgradeGeneratingHandler: handler,
		apply:                        apply,
		name:                         name,
		gvk:                          controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	RegisterNodeUpgradeStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type nodeUpgradeStatusHandler struct {
	client    NodeUpgradeClient
	condition condition.Cond
	handler   NodeUpgradeStatusHandler
}

func (a *nodeUpgradeStatusHandler) sync(key string, obj *v1.NodeUpgrade) (*v1.NodeUpgrade, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	obj.Status = newStatus
	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(obj, "", nil)
		} else {
			a.condition.SetError(obj, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, obj.Status) {
		var newErr error
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type nodeUpgradeGeneratingHandler struct {
	NodeUpgradeGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *nodeUpgradeGeneratingHandler) Handle(obj *v1.NodeUpgrade, status v1.NodeUpgradeStatus) (v1.NodeUpgradeStatus, error) {
	objs, newStatus, err := a.NodeUpgradeGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	apply := a.apply

	if !a.opts.DynamicLookup {
		apply = apply.WithStrictCaching()
	}

	if !a.opts.AllowCrossNamespace && !a.opts.AllowClusterScoped {
		apply = apply.WithSetOwnerReference(true, false).
			WithDefaultNamespace(obj.GetNamespace()).
			WithListerNamespace(obj.GetNamespace())
	}

	if !a.opts.AllowClusterScoped {
		apply = apply.WithRestrictClusterScoped()
	}

	return newStatus, apply.
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
