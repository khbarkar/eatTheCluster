package k8s

import (
	"bytes"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/kristinb/eatthecluster/internal/game"
)

type ChaosEngine struct {
	client *Client
	config *rest.Config
	dryRun bool
}

func NewChaosEngine(client *Client, config *rest.Config, dryRun bool) *ChaosEngine {
	return &ChaosEngine{client: client, config: config, dryRun: dryRun}
}

func (ce *ChaosEngine) Kill(ctx context.Context, res game.Resource) error {
	if ce.dryRun {
		return nil
	}
	switch res.Kind {
	case "Pod":
		return ce.client.Clientset.CoreV1().Pods(res.Namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "Deployment":
		return ce.client.Clientset.AppsV1().Deployments(res.Namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "Service":
		return ce.client.Clientset.CoreV1().Services(res.Namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	default:
		return fmt.Errorf("unsupported resource kind: %s", res.Kind)
	}
}

func (ce *ChaosEngine) Degrade(ctx context.Context, res game.Resource) error {
	if ce.dryRun {
		return nil
	}
	if res.Kind != "Pod" {
		return fmt.Errorf("can only degrade pods, got %s", res.Kind)
	}

	cmd := []string{"sh", "-c", "dd if=/dev/zero of=/dev/null bs=1M &"}
	req := ce.client.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(res.Name).
		Namespace(res.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdout:  true,
			Stderr:  true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(ce.config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return fmt.Errorf("exec failed: %w (stderr: %s)", err, stderr.String())
	}
	return nil
}
