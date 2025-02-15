package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"os"

	clientset "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned"
	informers "github.com/openfaas-incubator/faas-o6s/pkg/client/informers/externalversions"
	"github.com/openfaas-incubator/faas-o6s/pkg/controller"
	"github.com/openfaas-incubator/faas-o6s/pkg/server"
	"github.com/openfaas-incubator/faas-o6s/pkg/signals"
	"github.com/openfaas-incubator/faas-o6s/pkg/version"
)

var (
	masterURL  string
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Parse()

	sha, release := version.GetReleaseInfo()
	glog.Infof("Starting OpenFaaS controller version: %s commit: %s", release, sha)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building Kubernetes clientset: %s", err.Error())
	}

	faasClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building OpenFaaS clientset: %s", err.Error())
	}

	functionNamespace := "openfaas-fn"
	if namespace, exists := os.LookupEnv("function_namespace"); exists {
		functionNamespace = namespace
	}

	kubeInformerFactory := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, time.Second*30, functionNamespace, nil)
	faasInformerFactory := informers.NewSharedInformerFactory(faasClient, time.Second*30)

	ctrl := controller.NewController(kubeClient, faasClient, kubeInformerFactory, faasInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go faasInformerFactory.Start(stopCh)
	go server.Start(faasClient, kubeClient)

	if err = ctrl.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
