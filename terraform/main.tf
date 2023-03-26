resource "null_resource" "install_k3_master" {
  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "master"
      private_key = file("/home/master/.ssh/id_rsa_master_terraform")
      host        = "192.168.50.10"
    }

    inline = [
      "curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC='--write-kubeconfig-mode 644' sh -"
    ]   
  }
}

locals {
  k3_node_token_raw = "${file("/usr/local/bin/k3_node_token")}"
  k3_node_token = trim(replace(replace(local.k3_node_token_raw, "<<EOT", ""), "EOT", ""), "\n\t ")
}

resource "null_resource" "install_k3_slave1" {
  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "slave1"
      private_key = file("/home/master/.ssh/id_rsa_slave1_terraform")
      host        = "192.168.50.11"
      port        = "5522"
    }

    inline = [
      "sudo curl -sfL https://get.k3s.io | K3S_URL=https://192.168.50.10:6443 K3S_TOKEN=${local.k3_node_token} sh -s -- --with-node-id"
    ]
  }
}

resource "null_resource" "setup_k3_slave1" {
  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "master"
      private_key = file("/home/master/.ssh/id_rsa_master_terraform")
      host        = "192.168.50.10"
    }

    script = "/usr/local/bin/repositories/ryanschnabel-com/terraform/scripts/k3_agent_setup.sh"
  }
}
