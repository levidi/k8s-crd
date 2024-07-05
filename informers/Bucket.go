package informers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"levi.com/bucket-operator/types"
)

func GetInformerBucket(config *rest.Config) cache.SharedIndexInformer {

	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating dynamic client: %v", err)
	}

	// Create an informer factory for dynamic client
	gvr := schema.GroupVersionResource{
		Group:    "levi.com", // API group
		Version:  "v1",       // API version
		Resource: "buckets",  // resource name
	}

	bucketInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return dynamicClient.Resource(gvr).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				watch, err := dynamicClient.Resource(gvr).Watch(context.TODO(), options)
				if err != nil {
					klog.Errorf("Failed to start watching resources: %v", err)
					return nil, err
				}

				return watch, nil
			},
		},
		&unstructured.Unstructured{}, // Use unstructured.Unstructured here
		0,                            // resyncPeriod
		cache.Indexers{},
	)

	// Set up event handlers for Bucket
	bucketInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)
			customResource := &types.Bucket{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, customResource)
			if err != nil {
				klog.Errorf("Error converting object: %v", err)
				return
			}
			fmt.Printf("Bucket added: %s\n", customResource.Name)
			// Handle logic for added bucket
			fmt.Printf("BucketName: %s, Region: %s\n", customResource.Spec.BucketName, customResource.Spec.Region)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			unstructuredObj := newObj.(*unstructured.Unstructured)
			customResource := &types.Bucket{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, customResource)
			if err != nil {
				klog.Errorf("Error converting object: %v", err)
				return
			}
			fmt.Printf("Bucket updated: %s\n", customResource.Name)
			// Handle logic for updated bucket
			fmt.Printf("BucketName: %s, Region: %s\n", customResource.Spec.BucketName, customResource.Spec.Region)
		},
		DeleteFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)
			customResource := &types.Bucket{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, customResource)
			if err != nil {
				klog.Errorf("Error converting object: %v", err)
				return
			}
			fmt.Printf("Bucket deleted: %s\n", customResource.Name)
			// Handle logic for deleted bucket
		},
	})
	return bucketInformer
}
