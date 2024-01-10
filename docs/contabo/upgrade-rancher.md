# Upgrade Rancher Guide

This guide outlines the steps for upgrading Rancher.

### 1\. Change to the cluster where Rancher is installed:

```bash
kubectl config use-context [context-name]
```

### 2\. Refresh your Helm repository:

```bash
helm repo update
```

### 3\. Identify the repository used for Rancher installation:

```bash
helm repo add rancher-stable https://releases.rancher.com/server-charts/stable
```

### 4\. Download the latest Rancher chart from the Helm repository:

```bash
helm fetch rancher-stable/rancher
```

### 6\. Save current settings to a file

```bash
helm get values rancher -n cattle-system -o yaml > values.yaml
```

### 6\. View available charts and their versions

```bash
helm search repo rancher-stable/rancher --versions
```

### 7\. Upgrade Rancher to a specific version:

```bash
helm upgrade rancher rancher-<CHART_REPO>/rancher \
  --namespace cattle-system \
  -f values.yaml \
  --version=2.7.9
```


