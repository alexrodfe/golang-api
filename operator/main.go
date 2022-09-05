package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	cs "github.com/alexrodfe/golang-api/operator/clientset"
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
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
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	// creates the custom client than handles our custom resources
	customClient, err := cs.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	for {
		// get pods in default namespace
		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		//get deployment objects
		deploymentObjects, err := customClient.DeploymentObejcts("default").List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error obtaining deployment objects: %s\n", err.Error())
		}
		fmt.Printf("Number of deployment objects found: %d:\n", len(deploymentObjects.Items))

		// deployment info and update operations
		fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
		list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		var goapiDeployment appsv1.Deployment
		for _, d := range list.Items {
			if d.Labels["app"] == "goapi" {
				goapiDeployment = d
			}
			fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
		}

		var nReplicas = int32(len(deploymentObjects.Items))
		goapiDeployment.Spec.Replicas = &nReplicas

		updatedDeployment, updateErr := deploymentsClient.Update(context.TODO(), &goapiDeployment, metav1.UpdateOptions{})
		if updateErr != nil {
			fmt.Printf("Could not update deployment: %s\n", updateErr.Error())
		}

		fmt.Printf("Deployment %s updated to %d number of replicas\n", updatedDeployment.Name, *updatedDeployment.Spec.Replicas)

		time.Sleep(10 * time.Second)
	}
}
