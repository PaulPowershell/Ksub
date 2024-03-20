package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/manifoldco/promptui"
)

func main() {
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Fatalf("Failed to create Azure authorizer: %v", err)
		return
	}

	subscriptionsClient := subscriptions.NewClient()
	subscriptionsClient.Authorizer = authorizer

	ctx := context.Background()

	result, err := subscriptionsClient.List(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve subscriptions: %v", err)
		return
	}

	subscriptionNames := make([]string, len(result.Values()))

	for i, subscription := range result.Values() {
		subscriptionNames[i] = *subscription.DisplayName
	}

	sort.Strings(subscriptionNames)

	prompt := promptui.Select{
		Label: "Select a subscription",
		Items: subscriptionNames,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U000027A4 {{ . | cyan | bold }}",
			Inactive: "  {{ . }}",
		},
		Size:         50,
		HideSelected: true,
	}

	_, selected, err := prompt.Run()
	if err != nil {
		log.Println("Failed to prompt selection", err)
		return
	}

	setSubscription(selected)
	setKubernetesContext(selected)
}

func setSubscription(subscriptionID string) error {
	// Check if the subscription ID is provided
	if subscriptionID == "" {
		return fmt.Errorf("subscription ID is required")
	}

	// Create the command to change the Azure subscription context
	cmd := exec.Command("az", "account", "set", "--subscription", subscriptionID)

	// Set the command's output to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change Azure subscription context: %v", err)
	} else {
		fmt.Println("✔ Switched to subscription", subscriptionID)
	}
	return nil
}

func setKubernetesContext(selectedID string) {
	subscriptionName := selectedID

	// Take spoke number (lot) on subscription
	re := regexp.MustCompile(`\d{6}`)
	lot := re.FindString(subscriptionName)

	cmd := exec.Command("kubectx")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to execute 'kubectx' command: %v", err)
		return
	}

	// Search cluster with same spoke
	kubecontext := ""
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
			fmt.Printf("Failed to switch to kubecontext: %v", err)
			return
		} else {
			fmt.Println("✔ Switched to cluster", kubecontext)
		}
	} else {
		fmt.Println("❌ Pas de cluster associé")
	}
}
