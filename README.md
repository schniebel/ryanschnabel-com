# ryanschnabel-com

Repo holding code and Kubernetes configuration for webside.

## SSL Termination

Uses Certbot with Letsencrypt as the CA, and configuring with Cloud Flare. Configuration in the SSL folder. 
Certbot installed via Helm, and uses values and CRDs below:

  kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
  helm install cert-manager jetstack/cert-manager --namespace cert-manager --values=values.yaml --version v1.9.1


'values.yaml' - helm values defining cloudflare nameservers
'secret-cf-token.yaml' - defines secret that holds cloudflare api authentication
'certificate.yaml' - defines certificate that is usable by cluster
'letsencrypt.yaml' - exposes certificate so it is usable by cluster
