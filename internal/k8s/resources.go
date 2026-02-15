package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kristinb/eatthecluster/internal/game"
)

func (c *Client) ListResources(ctx context.Context) ([]game.Resource, error) {
	var resources []game.Resource

	pods, err := c.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	for _, pod := range pods.Items {
		resources = append(resources, game.Resource{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Kind:      "Pod",
		})
	}

	deployments, err := c.Clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	for _, dep := range deployments.Items {
		resources = append(resources, game.Resource{
			Name:      dep.Name,
			Namespace: dep.Namespace,
			Kind:      "Deployment",
		})
	}

	services, err := c.Clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	for _, svc := range services.Items {
		resources = append(resources, game.Resource{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Kind:      "Service",
		})
	}

	return resources, nil
}
