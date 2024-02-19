/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: get_resource
 * @Version: 1.0.0
 * @Date: 2024/2/8 15:46
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package AlterResource

import (
	"context"
	"errors"
	"fmt"
	"github.com/Einic/cops/zaplog"
	"github.com/jedib0t/go-pretty/v6/text"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetCurrentContainerResources(containers []corev1.Container, containerName string) (string, string, string, string) {
	for _, container := range containers {
		if container.Name == containerName {
			limitsCPU := container.Resources.Limits[corev1.ResourceCPU]
			limitsMemory := container.Resources.Limits[corev1.ResourceMemory]
			requestsCPU := container.Resources.Requests[corev1.ResourceCPU]
			requestsMemory := container.Resources.Requests[corev1.ResourceMemory]

			return limitsCPU.String(), limitsMemory.String(), requestsCPU.String(), requestsMemory.String()
		}
	}

	return "", "", "", ""
}

func GetStatusText(alterStatus string) string {
	switch alterStatus {
	case "Success":
		return text.FgGreen.Sprint(alterStatus)
	case "Failed":
		return text.FgRed.Sprint(alterStatus)
	default:
		return alterStatus
	}
}

// Function to get the status of the resource
func GetStatus(status appsv1.DeploymentStatus) string {
	if status.AvailableReplicas == status.Replicas {
		return text.FgGreen.Sprint("Available")
	} else if status.AvailableReplicas > 0 {
		return text.FgYellow.Sprint("Partial Available")
	} else {
		return text.FgRed.Sprint("Not Available")
	}
}

// Function to get the status of the resource
func GetStatusStatefulSet(status appsv1.StatefulSetStatus) string {
	if status.ReadyReplicas == status.Replicas {
		return text.FgGreen.Sprint("Available")
	} else if status.ReadyReplicas > 0 {
		return text.FgYellow.Sprint("Partial Available")
	} else {
		return text.FgRed.Sprint("Not Available")
	}
}

// GetPodQoS retrieves the Quality of Service (QoS) of a pod.
func GetPodQoS(clientset *kubernetes.Clientset, ResourcesName, namespace string, logger zaplog.Logger) (string, error) {
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", ResourcesName),
	})
	if err != nil {
		logger.Error("Error fetching pods", zap.String("ResourcesName", ResourcesName), zap.String("Namespace", namespace), zap.Error(err))
		return "", err
	}

	if len(podList.Items) == 0 {
		errMsg := "It may be that the app label is missing, please handle it manually."
		return "", errors.New(errMsg)
	}

	// Assuming all pods have the same QoS, so we just pick the first one
	return string(podList.Items[0].Status.QOSClass), nil
}
