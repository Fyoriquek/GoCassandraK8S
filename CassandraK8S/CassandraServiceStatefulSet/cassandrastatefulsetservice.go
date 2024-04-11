package main

import (
       "context"
        "flag"
        "fmt"
        "path/filepath"

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


        serviceClient := clientset.CoreV1().Services("cassandra-go")
        service := &apiv1.Service{
                ObjectMeta: metav1.ObjectMeta{
                        Name: "cassandra",
                        Labels: map[string]string{
                                "app": "cassandra",
                        },
                },
                                Spec: apiv1.ServiceSpec{
					ClusterIP: "None", 
                                        Ports: []apiv1.ServicePort{
                                                {
                                                Port: 9042,
                                                },
                                        },
                                        Selector: map[string]string{
                                                "app": "cassandra",
                                        },
                                        Type:"ClusterIP",

                        },
        }


        fmt.Println("Creating service...")
        result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
        if err != nil {
                panic(err)
        }
        fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())

}


func int32Ptr(i int32) *int32 { return &i }

