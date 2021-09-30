package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"log"
	"os"

	//"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/homedir"
	//"k8s.io/client-go/util/retry"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/office/.kube/config", "path")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	myserviece := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "hi-svc-x",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"finder": "hi-x",
			},
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Protocol: "TCP",
					Port:     8002,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: 8081,
						StrVal: "8081",
					},
				},
			},
		},
	}

	fmt.Println("creating service...\n")
	result, err := clientset.CoreV1().Services("default").Create(context.Background(), myserviece, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("service created %q. \n", result.GetObjectMeta().GetName())

	prompt()

	mydeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hi-dep-x",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"finder": "hi-x",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hi-pod-x",
					Labels:	map[string]string{
								"finder": "hi-x",
							},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						apiv1.Container{
							Name: "hi-container-x",
							Image: "tasdidur/mini-server-2",
						},
					},
				},
			},
		},
	}

	fmt.Println("creating deployment...\n")
	result2, err := clientset.AppsV1().Deployments("default").Create(context.Background(),mydeployment,metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deployment created %q. \n", result2.GetObjectMeta().GetName())
	
	prompt()
	myingress := &v12.Ingress{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "hi-ingress-x",
			Namespace:                  "default",
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/use-regex": "true",
			},
		},
		Spec: v12.IngressSpec{
			Rules: []v12.IngressRule{
				v12.IngressRule{
					Host:             "tasdid.com",
					IngressRuleValue: v12.IngressRuleValue{
						HTTP: &v12.HTTPIngressRuleValue{
							Paths: []v12.HTTPIngressPath{
								{
									Path:     "/",
									PathType: func() *v12.PathType {
										pt := v12.PathTypePrefix
										return &pt
									}(),
									Backend:  v12.IngressBackend{
										Service:  &v12.IngressServiceBackend{
											Name: "hi-svc-x",
											Port: v12.ServiceBackendPort{
												Number: 8002,
											},
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

	fmt.Println("creating ingress...\n")

	result3, err := clientset.NetworkingV1().Ingresses("default").Create(context.Background(),myingress,metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ingress created %q. \n", result3.GetObjectMeta().GetName())


}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }