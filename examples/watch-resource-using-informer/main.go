package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Here we use NewSharedInformerFactoryWithOptions to filter some pods to list-watch
	// NewSharedInformerFactory is an easier way instead.
	kubeInformerFactory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Second*15,
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = ""
		}),
		informers.WithNamespace("default"))

	podInformer := kubeInformerFactory.Core().V1().Pods()
	informer := podInformer.Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if pod, ok := obj.(*v1.Pod); ok {
				fmt.Printf("pod added: %v \n", pod.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			if pod, ok := obj.(*v1.Pod); ok {
				fmt.Printf("pod deleted: %v \n", pod.Name)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if pod, ok := newObj.(*v1.Pod); ok {
				fmt.Printf("pod changed: %v \n", pod.Name)
			}
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	kubeInformerFactory.Start(stopCh)

	// Wait for the caches to be synced before starting workers
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		return
	} else {
		fmt.Printf("WaitForCacheSync done\n")
	}

	// Now we can use lister to get pods
	lister := podInformer.Lister()
	pod, err := lister.Pods("default").Get("busybox")
	if kerrors.IsNotFound(err) {
		fmt.Printf("pod busybox not found\n")
	} else if err != nil {
		fmt.Printf("unable to retrieve pod %v from store: %v\n", "busybox", err)
	}
	if pod != nil {
		fmt.Printf("pod busybox status: %v\n", pod.Status.Phase)
	}

	<-stopCh
}
