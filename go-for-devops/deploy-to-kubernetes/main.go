package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/utils/ptr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	k8sClient, err := newK8sClient()
	panicError(err)

	ns, err := createNamespace(ctx, k8sClient, "nginx")
	defer deleteNamespace(ctx, k8sClient, ns)
	panicError(err)

	deployment, err := createDeployNginx(ctx, k8sClient, ns, "nginx")
	defer deleteDeployNginx(ctx, k8sClient, ns, deployment)
	panicError(err)

	service, err := createServiceNginx(ctx, k8sClient, ns, "nginx")
	defer deleteServiceNginx(ctx, k8sClient, ns, service)
	panicError(err)

	ingress, err := createIngressNginx(ctx, k8sClient, ns, "nginx")
	defer deleteIngressNginx(ctx, k8sClient, ns, ingress)
	panicError(err)

	err = waitForReadyReplicas(ctx, k8sClient, deployment, time.Second*10)
	panicError(err)

	podList, err := listPods(ctx, k8sClient, ns)
	panicError(err)

	for _, pod := range podList.Items {
		podName := pod.Name
		go func() {
			opts := &corev1.PodLogOptions{
				Container: deployment.Spec.Template.Spec.Containers[0].Name,
				Follow:    true,
			}
			podLogs, _ := k8sClient.CoreV1().Pods(ns.Name).GetLogs(podName, opts).Stream(ctx)

			_, _ = os.Stdout.ReadFrom(podLogs)
		}()
	}

	waitForExitSignal()
}

func newK8sClient() (*kubernetes.Clientset, error) {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientSet
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createNamespace(ctx context.Context, k8sClient *kubernetes.Clientset, name string) (*corev1.Namespace, error) {
	fmt.Printf("Creating namespace %q.\n", name)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return k8sClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
}

func deleteNamespace(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace) error {
	fmt.Printf("Deleting namespace %q.\n", ns.Name)
	return k8sClient.CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{})
}

func createDeployNginx(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace, name string) (*appv1.Deployment, error) {
	fmt.Printf("Creating deployment %q in namespace %q.\n", name, ns.Name)
	var (
		matchLabel = map[string]string{"app": "nginx"}
		objMeta    = metav1.ObjectMeta{
			Name:      name,
			Namespace: ns.Name,
			Labels:    matchLabel,
		}
	)

	deployment := &appv1.Deployment{
		ObjectMeta: objMeta,
		Spec: appv1.DeploymentSpec{
			Replicas: ptr.To[int32](2),
			Selector: &metav1.LabelSelector{MatchLabels: matchLabel},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchLabel,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: "nginxdemos/hello:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return k8sClient.AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
}

func deleteDeployNginx(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace, deployment *appv1.Deployment) error {
	fmt.Printf("Deleting deployment %q.\n", deployment.Name)
	return k8sClient.AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{})
}

func createServiceNginx(ctx context.Context, clientSet *kubernetes.Clientset, ns *corev1.Namespace, name string) (*corev1.Service, error) {
	fmt.Printf("Creating service %q in namespace %q.\n", name, ns.Name)
	var (
		matchLabel = map[string]string{"app": "nginx"}
		objMeta    = metav1.ObjectMeta{
			Name:      name,
			Namespace: ns.Name,
			Labels:    matchLabel,
		}
	)

	service := &corev1.Service{
		ObjectMeta: objMeta,
		Spec: corev1.ServiceSpec{
			Selector: matchLabel,
			Ports: []corev1.ServicePort{
				{
					Port:     80,
					Protocol: corev1.ProtocolTCP,
					Name:     "http",
				},
			},
		},
	}
	return clientSet.CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
}

func deleteServiceNginx(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace, service *corev1.Service) error {
	fmt.Printf("Deleting service %q.\n", service.Name)
	return k8sClient.CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{})
}

func createIngressNginx(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace, name string) (*netv1.Ingress, error) {
	fmt.Printf("Creating ingress %q in namespace %q.\n", name, ns.Name)
	var (
		prefix  = netv1.PathTypePrefix
		objMeta = metav1.ObjectMeta{
			Name:      name,
			Namespace: ns.Name,
		}
		ingressPath = netv1.HTTPIngressPath{
			PathType: &prefix,
			Path:     "/hello",
			Backend: netv1.IngressBackend{
				Service: &netv1.IngressServiceBackend{
					Name: name,
					Port: netv1.ServiceBackendPort{
						Name: "http",
					},
				},
			},
		}
	)

	ingress := &netv1.Ingress{
		ObjectMeta: objMeta,
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{ingressPath},
						},
					},
				},
			},
		},
	}
	return k8sClient.NetworkingV1().Ingresses(ns.Name).Create(ctx, ingress, metav1.CreateOptions{})
}

func deleteIngressNginx(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace, ingress *netv1.Ingress) error {
	fmt.Printf("\n\nDeleting ingress %q.\n", ingress.Name)
	return k8sClient.NetworkingV1().Ingresses(ns.Name).Delete(ctx, ingress.Name, metav1.DeleteOptions{})
}

func listPods(ctx context.Context, k8sClient *kubernetes.Clientset, ns *corev1.Namespace) (*corev1.PodList, error) {
	podList, err := k8sClient.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	fmt.Printf("Listing pods in %q namespace.\n", ns.Name)
	for _, pod := range podList.Items {
		fmt.Printf("# Pod: \n## namespace/name: %q\n## spec.containers[0].name: %q\n## spec.containers[0].image: %q\n", path.Join(pod.Namespace, pod.Name), pod.Spec.Containers[0].Name, pod.Spec.Containers[0].Image)
	}
	fmt.Printf("\n\n")
	return podList, nil
}

func waitForReadyReplicas(ctx context.Context, k8sClient *kubernetes.Clientset, deployment *appv1.Deployment, timeout time.Duration) error {
	fmt.Printf("Waiting for ready replicas in deployment %q\n", deployment.Name)

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout: replicas were not ready within %v seconds", timeout)
		default:
			expectedReplicas := *deployment.Spec.Replicas
			readyReplicas, _ := getReadyReplicasForDeployment(timeoutCtx, k8sClient, deployment)

			if readyReplicas == expectedReplicas {
				fmt.Printf("replicas are ready!\n\n")
				return nil
			}

			fmt.Printf("replicas are not ready yet. %d/%d\n", readyReplicas, expectedReplicas)
			time.Sleep(1 * time.Second)
		}
	}
}

func getReadyReplicasForDeployment(ctx context.Context, k8sClient *kubernetes.Clientset, deployment *appv1.Deployment) (int32, error) {
	dep, err := k8sClient.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	return dep.Status.ReadyReplicas, nil
}

func waitForExitSignal() {
	fmt.Printf("Type ctrl-c to exit\n\n")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func panicError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
