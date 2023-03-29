provider "null" {
  version = "3.1.0"
}

provider "kubernetes" {
  config_context_cluster = "default"
  version = "2.3.0"
}
