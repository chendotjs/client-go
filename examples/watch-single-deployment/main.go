package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	flag.Parse()
	deployName := os.Args[1]

	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deployWatch, err := clientset.AppsV1().Deployments("default").Watch(context.TODO(), meta_v1.ListOptions{
		FieldSelector: fields.ParseSelectorOrDie("metadata.name=" + deployName).String(),
	})

	if err != nil {
		deployWatch.Stop()
		panic(err.Error())
	}
	defer deployWatch.Stop()

	deployCh := deployWatch.ResultChan()
	for {
		event, ok := <-deployCh
		if !ok {
			panic("deployment watch channel had been closed")
		}
		switch event.Object.(type) {
		case *appsv1.Deployment:
			deploy := event.Object.(*appsv1.Deployment)
			fmt.Printf("Deployment %s got event: %s\n", deploy.Name, event.Type)
		default:
			fmt.Printf("unknown type\n")
		}
	}
}
