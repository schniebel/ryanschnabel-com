provider "kubernetes" {
  config_context_cluster = "default"
}


resource "kubernetes_deployment" "welcome-nginx-prod" {
  metadata {
    name = "welcome-nginx-prod"
  }

  spec {
    selector {
      match_labels = {
        app = "welcome-nginx-prod"
      }
    }

    replicas = 1

    template {
      metadata {
        labels = {
          app = "welcome-nginx-prod"
        }
      }

      spec {
        container {
          image = "schniebel/welcome-page:latest"
          name  = "welcome-nginx-prod"
          ports {
            container_port = 80
            name           = "http"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "welcome-nginx-svc-prod" {
  metadata {
    name = "welcome-nginx-svc-prod"
  }

  spec {
    selector = {
      app = "welcome-nginx-prod"
    }

    port {
      name        = "http"
      port        = 80
      target_port = 8080
    }

    type = "ClusterIP"
  }
}
