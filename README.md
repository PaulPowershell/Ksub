# Ksub - Kubernetes Subscription and Context Switcher
Ksub is a command-line tool for switching between Azure subscriptions and Kubernetes clusters with ease.

## Last Build
[![Go Build & Release](https://github.com/VegaCorporoptions/Ksub/actions/workflows/go.yml/badge.svg)](https://github.com/VegaCorporoptions/Ksub/actions/workflows/go.yml)

## Prerequisites
Before using this tool, ensure you have the following prerequisites:

- Azure CLI (az) installed and configured.
- kubectl installed and configured with access to your Kubernetes clusters.
- kubectx installed

## Installation
Clone the repository to your local machine:

```bash
git clone https://github.com/VegaCorporoptions/Ksub
cd your-repo
```

Build the Go application:
```bash
go build .
```

## Usage
Ksub allows you to switch between Azure subscriptions and Kubernetes clusters efficiently. Here's how to use it:

Run the Ksub application:
```bash
ksub
```
Select an Azure subscription from the list.  
Ksub will change your Azure subscription context to the selected one.  
Ksub will also switch your Kubernetes context to the corresponding cluster if available.

## Demo
![ksub.png](https://github.com/VegaCorporoptions/Ksub/blob/main/ksub.png?raw=true)

## License
This project is licensed under the MIT License. See the LICENSE file for details.
