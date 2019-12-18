package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	_, controller := cache.NewInformer(watchlist, &v1.Pod{}, time.Second*0, cache.ResourceEventHandlerFuncs{
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
	},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func fmtPod(pod *v1.Pod) string {
	str, _ := json.MarshalIndent(pod, "", "  ")
	return string(str)
}
