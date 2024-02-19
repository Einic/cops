/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: cnasst
 * @Version: 1.0.0
 * @Date: 2024/2/7 11:49
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package AlterResource

import (
	"context"
	"github.com/Einic/cops/lib"
	"github.com/Einic/cops/zaplog"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

// Function to update the deployment with new specifications
func UpdateDeployment(clientset *kubernetes.Clientset, deployment *appsv1.Deployment, replicas int, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace string, logger zaplog.Logger) lib.ResourceInfo {
	currentReplicas := int(*deployment.Spec.Replicas)
	alterReplicas := replicas

	// Get the current container resources
	CurrentLimitsCPU, CurrentLimitsMemory, CurrentRequestsCPU, CurrentRequestsMemory := GetCurrentContainerResources(deployment.Spec.Template.Spec.Containers, containersName)

	// Update the replicas
	deployment.Spec.Replicas = Int32Ptr(int32(replicas))

	// Update container resources
	UpdateContainerResources(deployment.Spec.Template.Spec.Containers, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory)

	// Update the deployment
	updatedDeployment, err := clientset.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("Error updating deployment", zap.String("WorkLoad", deployment.Name), zap.String("Namespace", deployment.Namespace), zap.Error(err))
		return lib.ResourceInfo{} // Return empty ResourceInfo in case of error
	}

	// Check if the deployment was actually updated
	var alterStatus string
	if err == nil && deploymentUpdated(updatedDeployment, deployment, replicas, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace) {
		// Deployment was actually updated
		alterStatus = "Success"
	} else {
		// Deployment was not updated
		alterStatus = "Failed"
	}

	// Update labels
	if err := UpdateLabels(clientset, deployment.Name, namespace, logger); err != nil {
		logger.Error("Error updating labels for deployment", zap.String("WorkLoad", deployment.Name), zap.String("Namespace", deployment.Namespace), zap.Error(err))
	}

	// Get Pod QoS
	PodQos, err := GetPodQoS(clientset, deployment.Name, namespace, logger)
	if err != nil {
		logger.Warn("Failed to get Pod QoS for deployment", zap.String("WorkLoad", deployment.Name), zap.String("Namespace", deployment.Namespace), zap.Error(err))
	}

	// Create a ResourceInfo instance to pass to PrintResources function
	resourceInfo := lib.ResourceInfo{
		DataTime:              time.Now().Format("2006-01-02 15:04:05"),
		Workload:              deployment.Name,
		ContainerName:         containersName,
		WorkType:              "deploy",
		Namespace:             namespace,
		CurrentReplicas:       currentReplicas,
		AlterReplicas:         alterReplicas,
		CurrentLimitsCPU:      CurrentLimitsCPU,
		AlterLimitsCPU:        limitsCPU,
		CurrentLimitsMemory:   CurrentLimitsMemory,
		AlterLimitsMemory:     limitsMemory,
		CurrentRequestsCPU:    CurrentRequestsCPU,
		AlterRequestsCPU:      requestsCPU,
		CurrentRequestsMemory: CurrentRequestsMemory,
		AlterRequestsMemory:   requestsMemory,
		PodQos:                PodQos,
		RunStatus:             GetStatus(deployment.Status),
		AlterStatus:           alterStatus,
	}

	return resourceInfo
}

// Check if the deployment was actually updated
func deploymentUpdated(updatedDeployment, originalDeployment *appsv1.Deployment, replicas int, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace string) bool {
	// Compare relevant fields to check if the deployment was actually updated
	if updatedDeployment == nil || originalDeployment == nil {
		return false
	}

	// Check if replicas are updated
	if int(*updatedDeployment.Spec.Replicas) != replicas {
		return false
	}

	// Check if container resources are updated
	updatedLimitsCPU, updatedLimitsMemory, updatedRequestsCPU, updatedRequestsMemory := GetCurrentContainerResources(updatedDeployment.Spec.Template.Spec.Containers, containersName)
	if updatedLimitsCPU != limitsCPU || updatedLimitsMemory != limitsMemory || updatedRequestsCPU != requestsCPU || updatedRequestsMemory != requestsMemory {
		return false
	}

	return true
}

// Function to update the statefulset with new specifications
func UpdateStatefulSet(clientset *kubernetes.Clientset, statefulSet *appsv1.StatefulSet, replicas int, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace string, logger zaplog.Logger) lib.ResourceInfo {
	currentReplicas := int(*statefulSet.Spec.Replicas)
	alterReplicas := replicas

	// Get the current container resources
	CurrentLimitsCPU, CurrentLimitsMemory, CurrentRequestsCPU, CurrentRequestsMemory := GetCurrentContainerResources(statefulSet.Spec.Template.Spec.Containers, containersName)

	// Update the replicas
	statefulSet.Spec.Replicas = Int32Ptr(int32(replicas))

	// Update container resources
	UpdateContainerResources(statefulSet.Spec.Template.Spec.Containers, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory)

	updatedStatefulSet, err := clientset.AppsV1().StatefulSets(statefulSet.Namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("Error updating statefulSet", zap.String("WorkLoad", statefulSet.Name), zap.String("Namespace", statefulSet.Namespace), zap.Error(err))

		return lib.ResourceInfo{} // Return empty ResourceInfo in case of error
	}

	// Check if the deployment was actually updated
	var alterStatus string
	if err == nil && StatefulSetUpdated(updatedStatefulSet, statefulSet, replicas, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace) {
		// StatefulSet was actually updated
		alterStatus = "Success"
	} else {
		// StatefulSet was not updated
		alterStatus = "Failed"
	}

	// Update labels
	if err := UpdateLabels(clientset, statefulSet.Name, namespace, logger); err != nil {
		logger.Error("Error updating labels for statefulset", zap.String("WorkLoad", statefulSet.Name), zap.String("Namespace", statefulSet.Namespace), zap.Error(err))
	}

	// Get Pod QoS
	PodQos, err := GetPodQoS(clientset, statefulSet.Name, namespace, logger)
	if err != nil {
		logger.Warn("Failed to get Pod QoS for statefulSet", zap.String("WorkLoad", statefulSet.Name), zap.String("Namespace", statefulSet.Namespace), zap.Error(err))
	}

	// Create a ResourceInfo instance to pass to PrintResources function
	resourceInfo := lib.ResourceInfo{
		DataTime:              time.Now().Format("2006-01-02 15:04:05"),
		Workload:              statefulSet.Name,
		ContainerName:         containersName,
		WorkType:              "sts",
		Namespace:             namespace,
		CurrentReplicas:       currentReplicas,
		AlterReplicas:         alterReplicas,
		CurrentLimitsCPU:      CurrentLimitsCPU,
		AlterLimitsCPU:        limitsCPU,
		CurrentLimitsMemory:   CurrentLimitsMemory,
		AlterLimitsMemory:     limitsMemory,
		CurrentRequestsCPU:    CurrentRequestsCPU,
		AlterRequestsCPU:      requestsCPU,
		CurrentRequestsMemory: CurrentRequestsMemory,
		AlterRequestsMemory:   requestsMemory,
		PodQos:                PodQos,
		RunStatus:             GetStatusStatefulSet(statefulSet.Status),
		AlterStatus:           alterStatus,
	}

	return resourceInfo
}

// Check if the deployment was actually updated
func StatefulSetUpdated(updatedStatefulSet, originalStatefulSet *appsv1.StatefulSet, replicas int, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory, namespace string) bool {
	// Compare relevant fields to check if the deployment was actually updated
	if updatedStatefulSet == nil || originalStatefulSet == nil {
		return false
	}

	// Check if replicas are updated
	if int(*updatedStatefulSet.Spec.Replicas) != replicas {
		return false
	}

	// Check if container resources are updated
	updatedLimitsCPU, updatedLimitsMemory, updatedRequestsCPU, updatedRequestsMemory := GetCurrentContainerResources(updatedStatefulSet.Spec.Template.Spec.Containers, containersName)
	if updatedLimitsCPU != limitsCPU || updatedLimitsMemory != limitsMemory || updatedRequestsCPU != requestsCPU || updatedRequestsMemory != requestsMemory {
		return false
	}

	return true
}

// Function to update container resources
func UpdateContainerResources(containers []corev1.Container, containersName, limitsCPU, limitsMemory, requestsCPU, requestsMemory string) {
	for i := range containers {
		if containers[i].Name == containersName {
			containers[i].Resources.Limits = corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(limitsCPU),
				corev1.ResourceMemory: resource.MustParse(limitsMemory),
			}
			containers[i].Resources.Requests = corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(requestsCPU),
				corev1.ResourceMemory: resource.MustParse(requestsMemory),
			}
			break
		}
	}
}

// Helper function to convert int32 to int32 pointer
func Int32Ptr(i int32) *int32 {
	return &i
}
