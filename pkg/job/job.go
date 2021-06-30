package job

import (
	"context"

	"github.com/NautiluX/managed-velero-plugin/pkg/k8s"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	BaseName           string = "managed-velero-plugin-status-patcha"
	ServiceAccountName        = BaseName + "-sa"
	RoleName                  = BaseName + "-role"
	RoleBindingName           = BaseName + "-rolebinding"
	MigrationNamespace        = "cluster-migration"
)

func CreateJob(cr string) error {
	genericObject, err := k8s.GetGenericObject(cr)
	if err != nil {
		return err
	}
	c, err := k8s.GetClient()
	if err != nil {
		return err
	}

	namespace := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: MigrationNamespace,
		},
	}
	_ = c.Create(context.TODO(), &namespace)

	serviceAccount := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceAccountName,
			Namespace: MigrationNamespace,
		},
	}
	_ = c.Create(context.TODO(), &serviceAccount)

	clusterRole := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RoleName,
			Namespace: MigrationNamespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"aws.managed.openshift.io"},
				Resources: []string{
					" '*'",
					" accountclaims",
					" accounts",
					" accountpools",
					" awsfederatedaccountaccesses",
					" awsfederatedroles",
				},
			},
		},
	}
	_ = c.Create(context.TODO(), &clusterRole)

	clusterRoleBinding := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: RoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      ServiceAccountName,
				Namespace: MigrationNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     RoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_ = c.Create(context.TODO(), &clusterRoleBinding)

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "apply-status-" + genericObject.Kind + "-" + genericObject.Metadata.Namespace + "/" + genericObject.Metadata.Name,
			Namespace: MigrationNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					ServiceAccountName: MigrationNamespace,
					Containers: []v1.Container{
						{
							Name:  "apply-status",
							Image: "quay.io/mdewald/managed-velero-plugin-status-patch",
							Args:  []string{cr},
						},
					},
				},
			},
		},
	}
	err = c.Create(context.TODO(), &job)
	return err
}
