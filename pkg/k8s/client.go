package k8s

import (
	awsv1alpha1 "github.com/openshift/aws-account-operator/pkg/apis/aws/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func GetClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	awsv1alpha1.AddToScheme(scheme)
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	c, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})

	if err != nil {
		return nil, err
	}
	return c, nil
}
