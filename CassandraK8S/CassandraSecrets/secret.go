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
	"secret/password"	
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

	
        secretClient := clientset.CoreV1().Secrets("cassandra-go")

        secret := &apiv1.Secret{
                ObjectMeta: metav1.ObjectMeta{
                        Name:"dbasecret",
                        //Namespace: "cassandra-go",
                 },
		TypeMeta: metav1.TypeMeta{
			Kind: "Opaque",
		},

		Data: map[string][]byte{
			"username": []byte(password.Username),
			"password": []byte(password.Password),
		},

        }

	
        fmt.Println("Creating secret...")
        result, err := secretClient.Create(context.TODO(), secret, metav1.CreateOptions{})
        if err != nil {
                panic(err)
        }
        fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())


}

func int32Ptr(i int32) *int32 { return &i }

