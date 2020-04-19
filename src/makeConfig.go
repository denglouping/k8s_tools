package main

import (
	"flag"
	"k8s.io/client-go/tools/clientcmd/api"
	"strings"

	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//4个参数
	//1.k8s配置文件
	//2.权限
	//3.集群或namespace
	//4.serviceaccount
	var kubeconfig string
	var permission string
	var scope string
	var serviceaccount string

	flag.StringVar(&kubeconfig, "kubeconfig", homeDir()+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	flag.StringVar(&permission, "p", "view", "permission: view or admin")
	flag.StringVar(&scope, "s", "test", "scope: cluster or name of  namespace")
	flag.StringVar(&serviceaccount, "n", "scope+permission", "serviceaccount-name: name of serviceaccount")
	flag.Parse()


	if serviceaccount == "scope+permission"{
		serviceaccount = scope + "-" + permission
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var token string
	var caCert string
	var host string

	serviceAccount_default, err := clientset.CoreV1().ServiceAccounts("default").Get("default", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	} else {
		tokenSecret, err := clientset.CoreV1().Secrets("default").Get(serviceAccount_default.Secrets[0].Name, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		} else {
			caCert = string(tokenSecret.Data["ca.crt"])
		}
	}

	host = strings.Replace(strings.Replace(config.Host, "http", "https", -1), "8080", "6443", -1)


	//命名空间的view权限
	if scope != "cluster" {
		serviceaccountClient := clientset.CoreV1().ServiceAccounts(scope)

		//check serviceaccount if exists
		result, _:= serviceaccountClient.Get(serviceaccount,metav1.GetOptions{})
		if result == nil {
			serviceAccount := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceaccount,
					Namespace: scope,
				},
			}
			_,err = serviceaccountClient.Create(serviceAccount)
		}


		rolebindingClient := clientset.RbacV1().RoleBindings(scope)
		rolebindingName := serviceaccount + "-" + permission

		result1, _ := rolebindingClient.Get(rolebindingName, metav1.GetOptions{})
		if result1 == nil {
			roleBinding := &rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: rolebindingName,
				},
				RoleRef: rbacv1.RoleRef{
					Kind: "ClusterRole",
					Name: permission,
				},
				Subjects: []rbacv1.Subject{{
					Kind:      "ServiceAccount",
					Name:      serviceaccount,
					Namespace: scope,
				},
				},
			}

			_, err = rolebindingClient.Create(roleBinding)
			if err != nil {
				panic(err.Error())
			}

		}

		time.Sleep(1 * time.Second)

		serviceAccount_New, err := clientset.CoreV1().ServiceAccounts(scope).Get(scope+"-"+permission, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		} else {
			tokenSecret, err := clientset.CoreV1().Secrets(scope).Get(serviceAccount_New.Secrets[0].Name, metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			} else {
				token = string(tokenSecret.Data["token"])
			}
		}

	}else {
		serviceaccountClient := clientset.CoreV1().ServiceAccounts("kube-system")
		serviceAccount := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceaccount,
				Namespace: "kube-system",
			},
		}
		result, err:= serviceaccountClient.Get(serviceaccount,metav1.GetOptions{})
		if result == nil{
			serviceaccountClient.Create(serviceAccount)
		}

		clusterrolebindingClient := clientset.RbacV1().ClusterRoleBindings()

		result1, err := clusterrolebindingClient.Get(serviceaccount,metav1.GetOptions{})
		if result1 == nil{
			clusterroleBinding := &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: serviceaccount,
				},
				RoleRef: rbacv1.RoleRef{
					Kind: "ClusterRole",
					Name: permission,
				},
				Subjects: []rbacv1.Subject{{
					Kind:      "ServiceAccount",
					Name:      serviceaccount,
					Namespace: "kube-system",
				},
				},
			}

			_,err := clusterrolebindingClient.Create(clusterroleBinding)
			if err != nil {
				panic(err.Error())
			}
		}
		else{
			result1.RoleRef
			clusterrolebindingClient.Patch(serviceaccount,)

		}

		time.Sleep(1 * time.Second)

		serviceAccount_New, err := clientset.CoreV1().ServiceAccounts("kube-system").Get(scope+"-"+permission, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		} else {
			tokenSecret, err := clientset.CoreV1().Secrets("kube-system").Get(serviceAccount_New.Secrets[0].Name, metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			} else {
				token = string(tokenSecret.Data["token"])
			}
		}

	}

	clientset.CoreV1().ServiceAccounts("kube-system").List(metav1.ListOptions{})

	newConfig := api.Config{
		Clusters: map[string]*api.Cluster{
			scope + "-" + permission: &api.Cluster{
				CertificateAuthorityData: []byte(caCert),
				Server:                   host,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			scope + "-" + permission: &api.AuthInfo{
				Token: token,
			},
		},
		Contexts: map[string]*api.Context{
			scope + "-" + permission: &api.Context{
				Cluster:   scope + "-" + permission,
				AuthInfo:  scope + "-" + permission,
				Namespace: scope,
			},
		},
		CurrentContext: scope + "-" + permission,
	}

	clientcmd.WriteToFile(
		newConfig, homeDir()+"/"+scope+"-"+permission+".config")

	//for {
	//	serviceAccounts, err := clientset.CoreV1().ServiceAccounts("kube-system").List(metav1.ListOptions{})
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	fmt.Printf("There are %d pods in the cluster\n", len(serviceAccounts.Items))
	//
	//	// Examples for error handling:
	//	// - Use helper functions like e.g. errors.IsNotFound()
	//	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	//	namespace := "default"
	//	pod := "example-xxxxx"
	//	_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	//	if errors.IsNotFound(err) {
	//		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	//	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	//		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
	//			pod, namespace, statusError.ErrStatus.Message)
	//	} else if err != nil {
	//		panic(err.Error())
	//	} else {
	//		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	//	}
	//
	//	time.Sleep(10 * time.Second)
	//}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
