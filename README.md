# ryanschnabel-com

Repo holding code and Kubernetes configuration for webside.

## Infrastructure Configuration

[ryanschnabel.com](https://ryanschnabel.com) is hosted on a [K3](https://k3s.io/) Kubernetes cluster on Raspberry Pi 4s.

Website is build Using [NGINX](https://hub.docker.com/_/nginx) containers and exposing/load balancing to outside traffic using [Traefik](https://traefik.io/traefik/)

Configuration located in the config folder.

`deployment.yaml` - defines deployment of pods

`service.yaml` - defines Traefik ingress and service exposure to pods

## SSL Termination

Uses [Certbot](https://certbot.eff.org/) with [Letsencrypt](https://letsencrypt.org/) as the CA, and configured with [Clouddflare](https://www.cloudflare.com/). 

Configuration located in the SSL folder. 

Certbot installed via [Helm](https://helm.sh/), and uses values and [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) below:

    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
    helm install cert-manager jetstack/cert-manager --namespace cert-manager --values=values.yaml --version v1.9.1


`values.yaml` - helm values defining cloudflare nameservers

`secret-cf-token.yaml` - defines secret that holds cloudflare api authentication

`certificate.yaml` - defines certificate that is usable by cluster

`letsencrypt.yaml` - exposes certificate so it is usable by cluster

## CI Pipeline

CI handled by [Github Actions](https://github.com/features/actions). The configuration of which can be found in the .github/workflows folder.

Steps handle 

- the authentication with [Docker Hub](https://hub.docker.com/) 
- The buiding and pushing of that image to my image repo (using the [buildx](https://github.com/docker/buildx) plugin to handle ARM64 architecture, which is what is running on the Raspberry Pi 4s)
- The execution of the kubectl commands that roll out the new image to the ryanschnabel.com domain.
