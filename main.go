package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	namespace  string
	secretName string
	secretFile string
)

func main() {
	// get environment for secrets
	namespace = os.Getenv("SECRET_NAMESPACE")
	secretName = os.Getenv("SECRET_NAME")
	secretFile = os.Getenv("SECRET_FILE")

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

	// read secret file into Map
	file, err := os.Open(secretFile)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	secretData := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				secretData[key] = value
			}
		}
	}

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		StringData: secretData,
	}

	_, err = clientset.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		// create secret if not exist
		result, err := clientset.CoreV1().Secrets(namespace).Create(secret)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Created secrets %q.\n", result.GetObjectMeta().GetName())
	} else {
		// update secret if exist
		result, err := clientset.CoreV1().Secrets(namespace).Update(secret)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Update secrets %q.\n", result.GetObjectMeta().GetName())
	}
}
