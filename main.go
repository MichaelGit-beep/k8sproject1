package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Println("Failed to auth using out of cluster method, trying in-cluster")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	if err != nil {
		log.Println(err)
	}
	for {
		events, err := clientset.CoreV1().Pods("default").Watch(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Println("Failed to list pods: ", err)
		}
		event := <-events.ResultChan()
		fmt.Println(event.Object.GetObjectKind().GroupVersionKind())
	}
}
