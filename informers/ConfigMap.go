package informers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func GetInformer(clientset *kubernetes.Clientset) cache.SharedIndexInformer {
	configMapInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return clientset.CoreV1().ConfigMaps(metav1.NamespaceAll).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return clientset.CoreV1().ConfigMaps(metav1.NamespaceAll).Watch(context.TODO(), options)
			},
		},
		&corev1.ConfigMap{},
		0, // resyncPeriod
		cache.Indexers{},
	)

	// Set up event handlers for ConfigMaps
	configMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap := obj.(*corev1.ConfigMap)
			fmt.Printf("ConfigMap added: %s/%s\n", configMap.Namespace, configMap.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			configMap := newObj.(*corev1.ConfigMap)
			fmt.Printf("ConfigMap updated: %s/%s\n", configMap.Namespace, configMap.Name)
		},
		DeleteFunc: func(obj interface{}) {
			configMap := obj.(*corev1.ConfigMap)
			fmt.Printf("ConfigMap deleted: %s/%s\n", configMap.Namespace, configMap.Name)
		},
	})
	return configMapInformer
}
