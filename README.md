# ryanschnabel-com

Repo holding code and Kubernetes configuration for webside.

## Infrastructure Configuration

[ryanschnabel.com](https://ryanschnabel.com) is hosted on a [K3](https://k3s.io/) Kubernetes cluster on Raspberry Pi 4s.

Website is build Using [NGINX](https://hub.docker.com/_/nginx) containers and exposing/load balancing to outside traffic using [Traefik](https://traefik.io/traefik/)

Configuration located in the config folder.

`deployment.yaml` - defines deployment of pods

`service.yaml` - defines Traefik ingress and service exposure to pods

## SSL Termination

Uses Certbot with Letsencrypt as the CA, and configured with Cloud Flare. 

Configuration located in the SSL folder. 

Certbot installed via Helm, and uses values and CRDs below:

    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
    helm install cert-manager jetstack/cert-manager --namespace cert-manager --values=values.yaml --version v1.9.1


`values.yaml` - helm values defining cloudflare nameservers

`secret-cf-token.yaml` - defines secret that holds cloudflare api authentication

`certificate.yaml` - defines certificate that is usable by cluster

`letsencrypt.yaml` - exposes certificate so it is usable by cluster
