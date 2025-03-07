package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pterm/pterm"
)

const (
	// correlation regex between subscription and cluster
	regexSubToCluster = `\d{6}`
)

func main() {
	// Create spinner & Start
	spinner, _ := pterm.DefaultSpinner.Start("Initialization running")

	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		pterm.Error.Printf("Failed to create Azure authorizer: %v", err)
		os.Exit(1)
	}

	subscriptionsClient := subscriptions.NewClient()
	subscriptionsClient.Authorizer = authorizer

	ctx := context.Background()

	result, err := subscriptionsClient.List(ctx)
	if err != nil {
		pterm.Error.Printf("Failed to retrieve subscriptions: %v", err)
		os.Exit(1)
	}

	subscriptionNames := make([]string, len(result.Values()))

	for i, subscription := range result.Values() {
		subscriptionNames[i] = *subscription.DisplayName
	}

	sort.Strings(subscriptionNames)
	spinner.Success("Initialization success")

	selector := pterm.DefaultInteractiveSelect.WithDefaultText("Select a subscription")
	selector.MaxHeight = 10
	selectedOption, _ := selector.WithOptions(subscriptionNames).Show() // The Show() method displays the options and waits for the user's input

	// Delete previous line X2
	fmt.Print("\033[F\033[K\033[F\033[K")

	setSubscription(selectedOption)
	setKubernetesContext(selectedOption)
}

func setSubscription(subscriptionID string) error {
	// Check if the subscription ID is provided
	if subscriptionID == "" {
		pterm.Error.Println("subscription ID is required")
		os.Exit(1)
	}

	// Create the command to change the Azure subscription context
	cmd := exec.Command("az", "account", "set", "--subscription", subscriptionID)

	// Set the command's output to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		pterm.Error.Println("failed to change Azure subscription context")
		os.Exit(1)
	} else {
		pterm.Success.Println("✔ Switched to subscription", subscriptionID)
	}
	return nil
}

func setKubernetesContext(selectedID string) {
	subscriptionName := selectedID

	// Take spoke number (lot) on subscription
	re := regexp.MustCompile(regexSubToCluster)
	lot := re.FindString(subscriptionName)

	cmd := exec.Command("kubectx")
	output, err := cmd.Output()
	if err != nil {
		pterm.Error.Printf("Failed to execute 'kubectx' command: %v", err)
		os.Exit(1)
	}

	// Search cluster with same spoke
	var kubecontext string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, lot) {
			kubecontext = strings.TrimSpace(line)
			break
		}
	}

	// Use kubecontext on same cluster's spoke
	if kubecontext != "" {
		cmd = exec.Command("kubectx", kubecontext)
		err = cmd.Run()
		if err != nil {
			pterm.Error.Printf("Failed to switch to kubecontext: %v", err)
			os.Exit(1)
		} else {
			pterm.Success.Println("✔ Switched to cluster " + kubecontext)
		}
	} else {
		pterm.Warning.Println("❌ No cluster associated")
	}
}
