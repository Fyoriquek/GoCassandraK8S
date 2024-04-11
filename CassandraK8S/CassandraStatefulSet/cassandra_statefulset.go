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

	
	StatefulSetClient := clientset.AppsV1().StatefulSets("cassandra-go")

	statefulset := &appsv1.StatefulSet{

		ObjectMeta: metav1.ObjectMeta{
			Name: "cassandra",
		},

		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "cassandra",
				},
			},
			ServiceName: "cassandra",
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "cassandra",
					},
				},
			        	Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "cassandra",
							Image: "cassandra:3.11.6",
							ImagePullPolicy: "IfNotPresent",
							Ports: []apiv1.ContainerPort{
								{
								        Name:          "intra-node",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 7000,
								},	
								{
                                                                        Name:          "tls-intra-node",
                                                                        Protocol:      apiv1.ProtocolTCP,
                                                                        ContainerPort: 7001,
                                                                },
								{
                                                                        Name:          "jmx",
                                                                        Protocol:      apiv1.ProtocolTCP,
                                                                        ContainerPort: 7199,
                                                                },
								{
                                                                        Name:          "cql",
                                                                        Protocol:      apiv1.ProtocolTCP,
                                                                        ContainerPort: 9042,
                                                                },
								{
                                                                        Name:          "webapp",
                                                                        Protocol:      apiv1.ProtocolTCP,
                                                                        ContainerPort: 8083,
                                                                },
							},
							Env: []apiv1.EnvVar{
								{
								Name: "CASSANDRA_SEEDS",
								Value: "cassandra-0.cassandra.cassandra-go.svc.cluster.local",
								},
								{
                             				        Name: "MAX_HEAP_SIZE",
                                				Value: "256M",
                            					},
                            					{
                                				Name: "HEAP_NEWSIZE",
                                				Value: "100M",
                            					},
                            					{
                                				Name: "CASSANDRA_CLUSTER_NAME",
                            					Value: "Cassandra", 
		       	    					},
                            					{
                                				Name: "CASSANDRA_DC",
                                				Value: "DC1",
                            					},
                            					{
                                				Name: "CASSANDRA_RACK",
                                				Value: "Rack1",
                            					},
                            					{
                                				Name: "CASSANDRA_ENDPOINT_SNITCH",
                                				Value: "GossipingPropertyFileSnitch",
                            					},
                            					{
                                				Name: "CASSANDRA_HOST",
                                				Value: "cassandra-0",
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


	fmt.Println("Creating statefulset...")
	result, err := StatefulSetClient.Create(context.TODO(), statefulset, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created statefulset %q.\n", result.GetObjectMeta().GetName())

}



