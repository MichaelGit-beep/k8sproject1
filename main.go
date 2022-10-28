package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

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

	pods, err := clientset.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	for {
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
		}

		deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Println(err)
		}
		for _, deploment := range deployments.Items {
			fmt.Printf("Deployment %s, Namespace: %s\n", deploment.Name, deploment.Namespace)
		}
		time.Sleep(10 * time.Second)
	}
}
