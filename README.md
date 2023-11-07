# ryanschnabel-com

Repo holding code and Kubernetes configuration for website.

[ryanschnabel.com](https://ryanschnabel.com) is hosted on a [K3](https://k3s.io/) Kubernetes cluster on Raspberry Pi 4s.

Website content was built using [Hugo](https://gohugo.io/). Using the vncnt-hugo theme.

## CI/CD

Any changes to the `ryanschnabel-com` folder on the `main` branch triggers a github actions workflow that will build the Hugo content and copy onto an `NGINX` image via the `Dockerfile` at that folder's directory level.

Once done, and the image is pushed to the [Docker Hub](https://hub.docker.com/) repo, github actions will update the `deployment.yaml` tag to match the new build's image. Triggering a `Flux reconciliation` and deploying the image onto the kubernetes cluster.

## Flux

Kubernetes cluster has [Flux](https://fluxcd.io/) deployed onto it. Reconciling any changes made to the `clusters` folder. Currently deployed the following resources onto the cluster via Flux:

### monitoring

Using [Grafana](https://grafana.com/grafana/), [Prometheus](https://prometheus.io/), and [Loki](https://grafana.com/oss/loki/) for log and metric aggregation, as well as displaying that data. Dashboard currently deployed to monitoring.ryanschnabel.com

### ryanschnabel-com

kubernetes resources that define the hosting of the Hugo content.

### Ingress

Exposure to outside traffic is handled using [Traefik](https://traefik.io/traefik/). Using the ingressRoute resource ([example](https://github.com/schniebel/ryanschnabel-com/blob/main/clusters/default/ryanschnabel-com/ingressRoute.yaml))

## DNS/SSL

DNS and SSL certificate management handled by [Cloudflare](https://www.cloudflare.com/)
