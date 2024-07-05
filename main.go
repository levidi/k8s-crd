package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"levi.com/bucket-operator/informers"
)

func main() {
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/leviditomazzomenezes/.kube/config", "Path to kubeconfig file")
	flag.Parse()

	stopCh := SetupSignalHandler() // Set up signals so we handle the first shutdown signal gracefully

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config) // Create a Kubernetes clientset
	if err != nil {
		klog.Fatalf("Error building clientset: %v", err)
	}

	bucketInformer := informers.GetInformerBucket(config)
	go bucketInformer.Run(stopCh) // Start the informer for Bucket

	configMapInformer := informers.GetInformer(clientset)
	go configMapInformer.Run(stopCh) // Start the informer for ConfigMaps

	// Wait for shutdown signal
	<-stopCh
	fmt.Println("Shutting down gracefully...")
}

// SetupSignalHandler sets up signal handler for graceful shutdown
func SetupSignalHandler() <-chan struct{} {
	stopCh := make(chan struct{})
	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		close(stopCh)
		<-signalCh
		os.Exit(1)
	}()
	return stopCh
}
