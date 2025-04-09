package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type controller struct {
	clientset      kubernetes.Interface
	depLister      appslisters.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.TypedRateLimitingInterface[string]
}

func newController(clientset kubernetes.Interface, depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:      clientset,
		depLister:      depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		queue:          workqueue.NewTypedRateLimitingQueue[string](workqueue.DefaultTypedControllerRateLimiter[string]()),
	}
	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	return c
}

func (c *controller) run(stopCh <-chan struct{}) {
	fmt.Println("controller started")
	if !cache.WaitForCacheSync(stopCh, c.depCacheSynced) {
		fmt.Print("waiting for cache  to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, stopCh)

	<-stopCh
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	key, quit := c.queue.Get()
	fmt.Println(key, quit)
	if quit {
		return false
	}

	defer c.queue.Done(key)

	ns, name, err := cache.SplitMetaNamespaceKey(fmt.Sprint(key))
	if err != nil {
		fmt.Printf("error while splitting meta namespace key: %s\n", err.Error())
		return false
	}

	//Check if the object has been deleted from k8s cluster
	fmt.Printf("Namespace %s Name %s\n", ns, name)
	_, err = c.clientset.AppsV1().Deployments(ns).Get(context.Background(), name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("Handle delete deployment %s\n", name)
		// delete Service
		err := c.clientset.CoreV1().Services(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Delete service %s error: %s\n", name, err.Error())
			return false
		}
		return true
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		fmt.Printf("error while syncing deployment: %s\n", err.Error())
		return false
	}
	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	dep, err := c.depLister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Printf("error while getting deployment: %s\n", err.Error())
		return err
	}
	//create service
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: depLabels(*dep),
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}
	_, err = c.clientset.CoreV1().Services(ns).Create(context.Background(), &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error while creating service: %s\n", err.Error())
	}
	return nil
}

func depLabels(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("add is called")
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		fmt.Printf("error getting key for object %+v: %v\n", obj, err)
		return
	}
	c.queue.Add(key)
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("del is called")
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		fmt.Printf("error getting key for object %+v: %v\n", obj, err)
		return
	}
	//fmt.Println("Inside HandleDel Func", key, err.Error())
	fmt.Println(key)

	c.queue.Add(key)

	fmt.Println(key, "added to queue")
}
