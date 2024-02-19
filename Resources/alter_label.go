/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: alter_label
 * @Version: 1.0.0
 * @Date: 2024/2/8 16:23
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package AlterResource

import (
	"context"
	"fmt"
	"github.com/Einic/cops/zaplog"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
)

// UpdateLabels updates labels for the given workload
func UpdateLabels(clientset *kubernetes.Clientset, workloadName, namespace string, logger zaplog.Logger) error {
	if err := AddAppLabelToPods(clientset, workloadName, namespace, logger); err != nil {
		return fmt.Errorf("error adding 'app' label to pods: %s", err.Error())
	}

	if err := UpdateAppLabelToWorkloadName(clientset, workloadName, namespace, logger); err != nil {
		return fmt.Errorf("error updating 'app' label to pods: %s", err.Error())
	}

	return nil
}

func GetRelatedPods(clientset *kubernetes.Clientset, workloadName, namespace string, logger zaplog.Logger) ([]corev1.Pod, error) {
	//Try to get related Pods based on Deployment name
	if _, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), workloadName, metav1.GetOptions{}); err == nil {
		//logger.Info("Related Pods found for the workload", zap.String("WorkloadName", workloadName), zap.String("Namespace", namespace))
		return getPodsByOwnerReference(clientset, workloadName, namespace, logger)
	}

	// Try to get related Pods based on StatefulSet name
	if _, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), workloadName, metav1.GetOptions{}); err == nil {
		//logger.Info("Related Pods found for the workload", zap.String("WorkloadName", workloadName), zap.String("Namespace", namespace))
		return getPodsByOwnerReference(clientset, workloadName, namespace, logger)
	}

	// If neither is found, an empty list is returned.
	logger.Warn("No related Pods found for the workload", zap.String("WorkloadName", workloadName), zap.String("Namespace", namespace))
	return nil, nil
}

func getPodsByOwnerReference(clientset *kubernetes.Clientset, workloadName, namespace string, logger zaplog.Logger) ([]corev1.Pod, error) {
	pods := []corev1.Pod{}

	// Get pods in the specified namespace
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("Error listing pods in namespace", zap.String("Namespace", namespace), zap.Error(err))
		return pods, err
	}

	// Use channels to collect matching Pods
	podChan := make(chan corev1.Pod)
	var wg sync.WaitGroup

	// Iterate through each pod to find those with matching OwnerReferences
	for _, pod := range podList.Items {
		wg.Add(1)
		//fmt.Println("getPodsByOwnerReference range pod: ", pod)
		go func(pod corev1.Pod) {
			defer wg.Done()
			for _, ownerRef := range pod.OwnerReferences {
				// Check if the name of the OwnerReferences matches the workload name
				if strings.HasPrefix(ownerRef.Name, workloadName+"-") {
					podChan <- pod
					break
				}
			}
		}(pod)
	}
	// Start a separate coroutine to close the channel and wait for all coroutines to complete
	go func() {
		wg.Wait()
		close(podChan)
	}()

	// Read matching Pods from the channel
	for pod := range podChan {
		pods = append(pods, pod)
	}
	return pods, nil
}

func AddAppLabelToPods(clientset *kubernetes.Clientset, workloadName, namespace string, logger zaplog.Logger) error {
	relatedPods, err := GetRelatedPods(clientset, workloadName, namespace, logger)
	if err != nil {
		return err
	}

	for _, pod := range relatedPods {
		if pod.Labels["app"] == "" {
			podCopy := pod.DeepCopy()
			if podCopy.Labels == nil {
				podCopy.Labels = make(map[string]string)
			}
			podCopy.Labels["app"] = workloadName

			_, err := clientset.CoreV1().Pods(namespace).Update(context.TODO(), podCopy, metav1.UpdateOptions{})
			if err != nil {
				logger.Error("Error adding 'app' label to pod", zap.String("PodName", pod.Name), zap.String("Namespace", namespace), zap.Error(err))
				continue
			}
			logger.Info("Added 'app' label to pod", zap.String("PodName", pod.Name), zap.String("Namespace", namespace))
		}
	}

	return nil
}

func UpdateAppLabelToWorkloadName(clientset *kubernetes.Clientset, workloadName, namespace string, logger zaplog.Logger) error {
	relatedPods, err := GetRelatedPods(clientset, workloadName, namespace, logger)
	if err != nil {
		return err
	}

	for _, pod := range relatedPods {
		if pod.Labels["app"] != workloadName {
			podCopy := pod.DeepCopy()
			if podCopy.Labels == nil {
				podCopy.Labels = make(map[string]string)
			}
			podCopy.Labels["app"] = workloadName

			_, err := clientset.CoreV1().Pods(namespace).Update(context.TODO(), podCopy, metav1.UpdateOptions{})
			if err != nil {
				logger.Error("Error updating 'app' label to match workloadName", zap.String("PodName", pod.Name), zap.String("Namespace", namespace), zap.Error(err))
				continue
			}
			logger.Info("Updated 'app' label to match workloadName", zap.String("PodName", pod.Name), zap.String("Namespace", namespace))
		}
	}

	return nil
}
