# Kubernetes Hands-On Project

## 📝 Agenda

Following task is designed as a hands-on experience with container orchestration, achieved by deploying open-source applications with different techniques.

1) Create (local) k8s cluster ([k3s](https://github.com/k3s-io/k3s/) is allowed)
2) Using `helm`|`helmfile`, deploy [Kubernetes Dashboard](https://github.com/kubernetes/dashboard) with ingress `traefik`|`nginx`
3) Using native k8s resource definitions, deploy [Glance](https://github.com/glanceapp/glance) with ingress. Configuration should be done in a separate file via [kustomize](https://github.com/kubernetes-sigs/kustomize)
4) Using any preferable [MCP server SDK](https://modelcontextprotocol.io/docs/sdk), create a server, which would expose tool, that returns information about a node (hostname, internal-ip). It should store information about incoming requests in it's file system. Wrap it in a container, deploy and using sidecar pattern create another container that would log information stored. Expose the ingress URL via [Ngrok](https://ngrok.com/).
5) Create a helm-chart for your newly created application and use it for further deployments.
