package main

import (
	"fmt"
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
	//c := &controller{
	//	clientset:      clientset,
	//	depLister:      depInformer.Lister(),
	//	depCacheSynced: depInformer.Informer().HasSynced,
	//	queue:          workqueue.NewTypedRateLimitingQueue[string](workqueue.DefaultTypedControllerRateLimiter[string]()),
	//}
	//depInformer.Informer().AddEventHandler(
	//	cache.ResourceEventHandlerFuncs{
	//		AddFunc:    handleAdd,
	//		DeleteFunc: handleDel,
	//	},
	//)
	fmt.Println("New controller Called")
	return nil
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

}

func handleAdd(obj interface{}) {
	fmt.Println("add is called")
}

func handleDel(obj interface{}) {
	fmt.Println("del is called")
}
