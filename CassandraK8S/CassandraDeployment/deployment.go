package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments("cassandra-go")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cassandra-createbooks",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "cassandra-createbooks",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "cassandra-createbooks",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "cassandra-createbooks",
							Image: "localhost:5000/microimage",
							ImagePullPolicy: "IfNotPresent",
							Ports: []apiv1.ContainerPort{
								{
								        Name:          "insertquery",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8083,
								},
							},
							Env: []apiv1.EnvVar{
								{
								Name: "CASSANDRA_HOST",
								Value: "cassandra-0.cassandra.cassandra-go.svc.cluster.local",
								},
								{
								Name: "CASSANDRA_USER"	,
								ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
									 		LocalObjectReference: apiv1.LocalObjectReference{
												Name: "dbasecret",	
											},
									 		Key: "username",
										},
									},
									
								},
								{
                                                                Name: "CASSANDRA_PASSWORD"  ,
                                                                ValueFrom: &apiv1.EnvVarSource{
                                                                                SecretKeyRef: &apiv1.SecretKeySelector{
                                                                                        LocalObjectReference: apiv1.LocalObjectReference{
                                                                                                Name: "dbasecret",
                                                                                        },
                                                                                        Key: "password",
                                                                                },
                                                                        },

                                                                },

            				        	 },
						},
					},
				},
			},
		},

	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}


func int32Ptr(i int32) *int32 { return &i }
