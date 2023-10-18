terraform {
  required_providers {
    coder = {
      source  = "coder/coder"
      version = "0.6.0"
    }
    # sysbox: provider change
    k8s = {
      source = "mingfang/k8s"
    }
  }
}

data "coder_provisioner" "me" {
}

provider "k8s" {
}

data "coder_workspace" "me" {
}

resource "coder_agent" "main" {
  arch           = data.coder_provisioner.me.arch
  os             = "linux"
  startup_script = <<EOF
    #!/bin/bash

    ### GIGO CONFIG
    ### write bootstrap script to temporary file
    cat > /tmp/gigo-bootstrap.py <<'    BOOTSTRAP_EOF'
<bootstrap_script>
    BOOTSTRAP_EOF

    ### GIGO CONFIG
    ### execute bootstrap script
    python3 -u /tmp/gigo-bootstrap.py --workspace-id ${lower(data.coder_workspace.me.id)} >> /tmp/gigo-bootstrap-script.log 2>&1

    # if we are here then something failed and we want to
    # stay alive for awhile so things can be debugged
    echo "Bootstrap failed"
    sleep 1d
    EOF

  ### GIGO CONFIG
  ### environment
  ### new entry for each key/value pair in the environment
  ### leave git config
  env = {
    GIT_AUTHOR_NAME     = "${data.coder_workspace.me.owner}"
    GIT_COMMITTER_NAME  = "${data.coder_workspace.me.owner}"
    GIT_AUTHOR_EMAIL    = "${data.coder_workspace.me.owner_email}"
    GIT_COMMITTER_EMAIL = "${data.coder_workspace.me.owner_email}"
    <environment>
  }
}

resource "coder_app" "code-server" {
  agent_id = coder_agent.main.id
  slug     = "code-server"
  display_name     = "code-server"
  ### GIGO CONFIG
  ### working_directory
  url      = "http://localhost:13337/?folder=<working_directory>"
  icon     = "/icon/code.svg"
}

resource "kubernetes_persistent_volume_claim" "home" {
  metadata {
    name      = "gigo-ws-${lower(data.coder_workspace.me.owner)}-${lower(data.coder_workspace.me.name)}-home"
    namespace = "coder"
  }
  wait_until_bound = false
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        ### GIGO CONFIG
        ### resources.disk
        storage = "<resources.disk>Gi"
      }
    }
  }
}

# sysbox: namechange
resource "k8s_core_v1_pod" "main" {
  count = data.coder_workspace.me.start_count
  metadata {
    name      = "gigo-ws-${lower(data.coder_workspace.me.owner)}-${lower(data.coder_workspace.me.name)}"
    namespace = "coder"
    # sysbox: namesapce annotation
    annotations = {
      "io.kubernetes.cri-o.userns-mode" = "auto:size=65536"
    }
  }
  spec {
    # sysbox: add special runtime
    runtime_class_name = "sysbox-runc"

    security_context {
      run_asuser = 0
      fsgroup    = 0
    }
    containers {
      name    = "dev"
      ### GIGO CONFIG
      ### base_container
      image   = "<base_container>"
      # sysbox: use this command to launch systemd before the container starts
      command = ["sh", "-c", <<EOF
      # create user
      echo "Creaing gigo user"
      useradd --create-home --shell /bin/bash gigo

      # initialize the gigo home directory using /etc/skeleton
      cp -r /etc/skel/. /home/gigo/

      # change ownership of coder directory
      echo "Ensuring directory ownership for gigo user"
      chown gigo:gigo -R /home/gigo

      # disable sudo for gigo user
      echo "gigo ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/gigo

      # Start the Coder agent as the "gigo" user
      # once systemd has started up
      echo "Waiting for systemd to start"
      sudo -u gigo --preserve-env=CODER_AGENT_TOKEN /bin/bash -- <<-'      EOT' &
      while [[ ! $(systemctl is-system-running) =~ ^(running|degraded) ]]
      do
        echo "Waiting for system to start... $(systemctl is-system-running)"
        sleep 2
      done

      echo "Starting Coder agent"
      ${coder_agent.main.init_script}
      EOT

      echo "Executing /sbin/init"
      exec /sbin/init

      echo "Exiting"
      EOF
      ]

      env {
        name  = "CODER_AGENT_TOKEN"
        value = coder_agent.main.token
      }
      volume_mounts {
        mount_path = "/home/gigo"
        name       = "home"
        read_only  = false
      }

      resources {
        requests = {
          cpu    = "500m"
          memory = "500Mi"
        }
        limits = {
          ### GIGO CONFIG
          ### resources.cpu
          cpu    = "<resources.cpu>"
          ### GIGO CONFIG
          ### resources.mem
          memory = "<resources.mem>G"
        }
      }
    }

    volumes {
      name = "home"
      persistent_volume_claim {
        claim_name = kubernetes_persistent_volume_claim.home.metadata.0.name
        read_only  = false
      }
    }
  }
}