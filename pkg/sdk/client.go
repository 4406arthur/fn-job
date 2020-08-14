package sdk

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//NewK8sCli used to setup a k8s client
func NewK8sCli(massterURL, kubeconfig string) (*kubernetes.Clientset, error) {

	if kubeconfig == "" && massterURL == "" {
		//klog.Warningf("Neither --kubeconfig nor --master was specified.  Using the inClusterConfig.  This might not work.")
		kubeconfig, err := restclient.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(kubeconfig)
		return clientset, nil
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags(massterURL, kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}
