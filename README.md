# ryanschnabel-com

Repo holding code and Kubernetes configuration for website.

## Table of Contents

- [Infrastructure Configuration](#infrastructure-configuration)
- [SSL Termination](#ssl-termination)
- [Cluster Monitoring (Grafana, Prometheus)](#cluster-monitoring)
- [CI Pipeline](#ci-pipeline)
- [Dynamic DNS](#dynamic-dns)
- [Dashboard](#dashboard)

## Infrastructure Configuration

[ryanschnabel.com](https://ryanschnabel.com) is hosted on a [K3](https://k3s.io/) Kubernetes cluster on Raspberry Pi 4s.

Website is built Using [NGINX](https://hub.docker.com/_/nginx) containers and exposing/load balancing to outside traffic using [Traefik](https://traefik.io/traefik/).

Test environment used for testing changes before CI push to production available at [test.ryanschnabel.com](https://test.ryanschnabel.com).

Configuration located in the [config](https://github.com/schniebel/ryanschnabel-com/tree/main/config) folder.

`auth-secret-template.yaml` - Template of secret used in test environment authorization.

`ingress-route.yaml` - Using [Traefik IngressRoute](https://doc.traefik.io/traefik/routing/providers/kubernetes-crd/#kind-ingressroute) object to handle routing of traffic to the test pods, production pods, and Traefik Dashboard. As well as specify middleware used for basic auth in the test environment, and TLS secret used for SSL handshake.

`middleware.yaml` - Using [traefik middleware](https://doc.traefik.io/traefik/middlewares/overview/) to define an authorization redirect to the test environment.

`ryanschnabel-com-prod.yaml` - Defines production pod deployments as well as service that exposes those pods to the cluster.

`ryanschnabel-com-test.yaml` - Defines test pod deployments and service that exposes those pods to the cluster

`traefik-nodeport.yaml` - Exposes ingress route to external traffic.

## SSL Termination

Uses [Certbot](https://certbot.eff.org/) with [Letsencrypt](https://letsencrypt.org/) as the CA, and configured with [Clouddflare](https://www.cloudflare.com/). 

Configuration located in the SSL folder. 

Certbot installed via [Helm](https://helm.sh/), and uses values and [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) below:

    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
    helm install cert-manager jetstack/cert-manager --namespace cert-manager --values=values.yaml --version v1.9.1


`values.yaml` - helm values defining cloudflare nameservers.

`secret-cf-token.yaml` - defines secret that holds cloudflare api authentication.

`certificate.yaml` - defines certificate that is usable by cluster.

`letsencrypt.yaml` - exposes certificate so it is usable by cluster.

## Cluster Monitoring

Cluster health and resources are monitored via Prometheus and displayed in dashboards using [Grafana](https://grafana.com/). Prometheus configuration and manifests generated via [Kube-Prometheus](https://github.com/prometheus-operator/kube-prometheus)

Dashboard exposed at monitoring.ryanschnabel.com

## CI Pipeline

CI handled by [Github Actions](https://github.com/features/actions). The configuration of which can be found in the [.github/workflows](https://github.com/schniebel/ryanschnabel-com/tree/main/.github/workflows) folder.

Steps handle: 

- Authentication with [Docker Hub](https://hub.docker.com/) 
- Buiding and pushing of that image to my image repo (using the [buildx](https://github.com/docker/buildx) plugin to handle ARM64 architecture, which is what is running on the Raspberry Pi 4s)
- Excecution of the kubectl commands that roll out new image to [test.ryanschnabel.com](https://test.ryanschnabel.com) for test verification.
- Pause for manual approval after test deployment. Manual verification handled in automatically created github issue. And waits for positive approval for moving on to production deployment. Using [trstringer/manual-approval](https://github.com/trstringer/manual-approval) in Github action to achieve this.
- After positive approval from manual verification execute the kubectl commands that roll out the new image to the production ryanschnabel.com domain. Note a manual rejection in previous step aborts the build pipeline.

## Dynamic DNS

My home lab network's ISP uses dynamic DNS. So a way is needed to make sure that my IP address on my home lab router matches what my DNS providor has its 'A' record for the domain.

Cloudflare is the Domain Registrar/ DNS management for ryanschnabel.com. In order to keep the Dynamic IP address provided by my ISP matched with Cloudflare, I am using [K0p1's Cloudflare DDNS Updater](https://github.com/K0p1-Git/cloudflare-ddns-updater).

## Dashboard

Traefik Dashboard for cluster routing available at [dashboard.ryanschnabel.com](https://dashboard.ryanschnabel.io). Dashboard exposed via ingress router mapping to TraefikService. Using same basic auth middleware used for the test environment.
