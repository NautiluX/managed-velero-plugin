package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/NautiluX/managed-velero-plugin/pkg/k8s"
	awsv1alpha1 "github.com/openshift/aws-account-operator/pkg/apis/aws/v1alpha1"
)

const (
	MaxCrWaitTime time.Duration = 5 * time.Minute
)

func main() {
	cr := os.Args[1]
	genericObject, err := k8s.GetGenericObject(cr)
	if err != nil {
		panic(err)
	}
	c, err := k8s.GetClient()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	fmt.Printf("Input CR:\n---------\n%v---------\n\n", cr)
	fmt.Printf("Start watching for CR kind %s at %v for %v.\n", genericObject.Kind, start, MaxCrWaitTime)
	for time.Since(start) < MaxCrWaitTime {
		switch genericObject.Kind {
		case "Account":
			patchAccount(cr, c)
		}
	}
	panic("expected account CR didn't appear")
}

func patchAccount(cr string, c client.Client) {
	accountIn := awsv1alpha1.Account{}
	err := json.Unmarshal([]byte(cr), &accountIn)
	if err != nil {
		panic(err)
	}

	var account *awsv1alpha1.Account
	fmt.Println("waiting for account CR " + accountIn.Namespace + "/" + accountIn.Name + " to appear")
	account, err = getAWSAccount(context.TODO(), c, accountIn.Namespace, accountIn.Name)
	if err == nil {
		// reset fields in status
		var mergePatch []byte

		mergePatch, _ = json.Marshal(map[string]interface{}{
			"status": accountIn.Status,
		})

		err = c.Status().Patch(context.TODO(), account, client.RawPatch(types.MergePatchType, mergePatch))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Status updated: %v", mergePatch)
		return
	}
	fmt.Println("not there. Delaying 10s")
	time.Sleep(10 * time.Second)

}

func getAWSAccount(ctx context.Context, cli client.Client, namespace, accountCRName string) (*awsv1alpha1.Account, error) {
	var account awsv1alpha1.Account
	if err := cli.Get(ctx, types.NamespacedName{
		Name:      accountCRName,
		Namespace: namespace,
	}, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
