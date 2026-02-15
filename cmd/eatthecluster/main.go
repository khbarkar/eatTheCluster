package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kristinb/eatthecluster/internal/game"
	"github.com/kristinb/eatthecluster/internal/k8s"
	"github.com/kristinb/eatthecluster/internal/tui"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Run without performing real chaos actions")
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig (default: ~/.kube/config)")
	flag.Parse()

	client, err := k8s.NewClient(*kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to cluster: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	resources, err := client.ListResources(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list resources: %v\n", err)
		os.Exit(1)
	}

	if len(resources) == 0 {
		fmt.Fprintln(os.Stderr, "No resources found in cluster.")
		os.Exit(1)
	}

	g := game.NewGame(resources, *dryRun)

	kubeconfigPath := *kubeconfig
	if kubeconfigPath == "" {
		home, _ := os.UserHomeDir()
		kubeconfigPath = home + "/.kube/config"
	}
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	engine := k8s.NewChaosEngine(client, config, *dryRun)

	chaosFunc := func(event game.ChaosEvent) {
		ctx := context.Background()
		switch event.Action {
		case "killed":
			engine.Kill(ctx, event.Resource)
		case "degraded":
			engine.Degrade(ctx, event.Resource)
		}
	}

	model := tui.NewModel(g, client.Context, chaosFunc)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
