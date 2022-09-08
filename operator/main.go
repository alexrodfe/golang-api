package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	kappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"

	cs "github.com/alexrodfe/golang-api/operator/clientset"
	do "github.com/alexrodfe/golang-api/operator/deployment-objects"
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var (
	deploymentsClient kappsv1.DeploymentInterface
	customClient      *cs.CustomClient
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// creates the deploymentsClient for namespace default
	deploymentsClient = clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	// creates the custom client than handles our custom resources
	customClient, err = cs.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	for {
		// deployment info and update operations
		fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
		deployments, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, deployment := range deployments.Items {
			fmt.Printf(" * %s (%d replicas)\n", deployment.Name, *deployment.Spec.Replicas)
		}

		deploymentObjectsByKind, err := collectDeploymentObjects()
		if err != nil {
			fmt.Printf("Error obtaining deployment objects: %s\n", err.Error())
		}

		for kind, do := range deploymentObjectsByKind {
			deployment := searchDeployment(deployments.Items, kind)
			if deployment == nil {
				fmt.Printf("Error: could not found deployment with app label '%s', skipping\n", kind)
				continue
			}
			processDeploymentObject(deployment, kind, do)
		}

		time.Sleep(10 * time.Second)
	}
}

// collectDeploymentObjects will fetch all deployment objects published and process them
// if any kind is repeated an error will promt and that resource will not be processed
func collectDeploymentObjects() (map[string]do.DeploymentObject, error) {
	deploymentObjects, err := customClient.DeploymentObejcts("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	doByKind := make(map[string]do.DeploymentObject, 0)
	repeatedKeys := make([]string, 0)
	for _, do := range deploymentObjects.Items {
		_, ok := doByKind[do.Spec.Kind]
		if ok { // it already exists
			repeatedKeys = append(repeatedKeys, do.Spec.Kind)
			fmt.Printf("ERROR: Several deployment objects for the kind %s have been detected, only 1 is supported. This kind will not be processed\n", do.Spec.Kind)
		} else {
			doByKind[do.Spec.Kind] = do
		}
	}
	for _, do := range repeatedKeys {
		delete(doByKind, do)
	}

	return doByKind, nil
}

func processDeploymentObject(deployment *appsv1.Deployment, kind string, do do.DeploymentObject) {
	if *deployment.Spec.Replicas != do.Spec.NReplicas {
		fmt.Printf("New change detected in %s deployment object, updating number of replicas..\n", do.Name)

		err := updateDeployment(deployment, do.Spec.NReplicas)
		if err != nil {
			fmt.Printf("Could not update deployment: %s\n", err.Error())
		}
	}
}

func searchDeployment(items []appsv1.Deployment, name string) *appsv1.Deployment {
	var res appsv1.Deployment
	var ok bool

	for _, item := range items {
		if item.Labels["app"] == name {
			res = item
			ok = true
		}
	}

	if !ok {
		return nil
	}

	return &res
}

func updateDeployment(deployment *appsv1.Deployment, nReplicas int32) error {
	deployment.Spec.Replicas = &nReplicas

	updatedDeployment, err := deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deployment %s updated to %d number of replicas\n", updatedDeployment.Name, *updatedDeployment.Spec.Replicas)

	return nil
}
