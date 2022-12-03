/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */




data "google_project" "project" {
  project_id = var.project_id
}


resource "google_project_service" "all" {
  for_each           = toset(var.gcp_service_list)
  project            = var.project_number
  service            = each.key
  disable_on_destroy = false
}

# Handle Permissions
variable "build_roles_list" {
  description = "The list of roles that build needs for"
  type        = list(string)
  default = [
    "roles/run.developer",
    "roles/vpaccess.user",
    "roles/iam.serviceAccountUser",
    "roles/run.admin",
    "roles/secretmanager.secretAccessor",
    "roles/artifactregistry.admin",
  ]
}

resource "google_service_account" "runsa" {
  project      = var.project_number
  account_id   = "${var.basename}-run-sa"
  display_name = "Service Account Three Tier App on Cloud Run"
}

resource "google_project_iam_member" "allrun" {
  project    = data.google_project.project.number
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.runsa.email}"
  depends_on = [google_project_service.all]
}


resource "google_project_iam_member" "allbuild" {
  for_each   = toset(var.build_roles_list)
  project    = var.project_number
  role       = each.key
  member     = "serviceAccount:${local.sabuild}"
  depends_on = [google_project_service.all]
}


resource "google_compute_network" "main" {
  provider                = google-beta
  name                    = "${var.basename}-private-network"
  auto_create_subnetworks = true
  project                 = var.project_id
}


resource "google_compute_global_address" "main" {
  name          = "${var.basename}-vpc-address"
  provider      = google-beta
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.main.id
  project       = var.project_id
  depends_on    = [google_project_service.all]
}

resource "google_service_networking_connection" "main" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.main.name]
  depends_on              = [google_project_service.all]
}

resource "google_vpc_access_connector" "main" {
  provider       = google-beta
  project        = var.project_id
  name           = "${var.basename}-vpc-cx"
  ip_cidr_range  = "10.8.0.0/28"
  network        = google_compute_network.main.id
  region         = var.region
  max_throughput = 300
  depends_on     = [google_compute_global_address.main, google_project_service.all]
}

resource "random_id" "id" {
  byte_length = 2
}

# Handle Database
resource "google_sql_database_instance" "main" {
  name             = "${var.basename}-db-${random_id.id.hex}"
  database_version = "MYSQL_5_7"
  region           = var.region
  project          = var.project_id
  settings {
    tier                  = "db-g1-small"
    disk_autoresize       = true
    disk_autoresize_limit = 0
    disk_size             = 10
    disk_type             = "PD_SSD"
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.main.id
    }
    location_preference {
      zone = var.zone
    }
  }
  deletion_protection = false
  depends_on = [
    google_project_service.all,
    google_service_networking_connection.main
  ]

  provisioner "local-exec" {
    working_dir = "${path.module}/../code/database"
    command     = "./load_schema.sh ${var.project_id} ${google_sql_database_instance.main.name}"
  }
}

# Handle redis instance
resource "google_redis_instance" "main" {
  authorized_network      = google_compute_network.main.id
  connect_mode            = "DIRECT_PEERING"
  location_id             = var.zone
  memory_size_gb          = 1
  name                    = "${var.basename}-cache"
  project                 = var.project_id
  redis_version           = "REDIS_6_X"
  region                  = var.region
  reserved_ip_range       = "10.137.125.88/29"
  tier                    = "BASIC"
  transit_encryption_mode = "DISABLED"
  depends_on              = [google_project_service.all]
}

# Handle artifact registry
resource "google_artifact_registry_repository" "todo_app" {
  provider      = google-beta
  format        = "DOCKER"
  location      = var.region
  project       = var.project_id
  repository_id = "${var.basename}-app"
  depends_on    = [google_project_service.all]
}

# Handle secrets
resource "google_secret_manager_secret" "redishost" {
  project = var.project_number
  replication {
    automatic = true
  }
  secret_id  = "redishost"
  depends_on = [google_project_service.all]
}

resource "google_secret_manager_secret_version" "redishost" {
  enabled     = true
  secret      = "projects/${var.project_number}/secrets/redishost"
  secret_data = google_redis_instance.main.host
  depends_on  = [google_project_service.all, google_redis_instance.main, google_secret_manager_secret.redishost]
}

resource "google_secret_manager_secret" "sqlhost" {
  project = var.project_number
  replication {
    automatic = true
  }
  secret_id  = "sqlhost"
  depends_on = [google_project_service.all]
}

resource "google_secret_manager_secret_version" "sqlhost" {
  enabled     = true
  secret      = "projects/${var.project_number}/secrets/sqlhost"
  secret_data = google_sql_database_instance.main.private_ip_address
  depends_on  = [google_project_service.all, google_sql_database_instance.main, google_secret_manager_secret.sqlhost]
}

resource "null_resource" "cloudbuild_api" {
  provisioner "local-exec" {
    working_dir = "${path.module}/../code/middleware"
    command     = "gcloud builds submit . --substitutions=_REGION=${var.region},_BASENAME=${var.basename}"
  }

  depends_on = [
    google_artifact_registry_repository.todo_app,
    google_secret_manager_secret_version.redishost,
    google_secret_manager_secret_version.sqlhost,
    google_project_service.all
  ]
}

resource "google_cloud_run_service" "api" {
  name     = "${var.basename}-api"
  location = var.region
  project  = var.project_id

  template {
    spec {
      service_account_name = google_service_account.runsa.email
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.basename}-app/api"
        env {
          name = "REDISHOST"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.redishost.secret_id
              key  = "latest"
            }
          }
        }
        env {
          name = "todo_host"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.sqlhost.secret_id
              key  = "latest"
            }
          }
        }
        env {
          name  = "todo_user"
          value = "todo_user"
        }
        env {
          name  = "todo_pass"
          value = "todo_pass"
        }
        env {
          name  = "todo_name"
          value = "todo"
        }

        env {
          name  = "REDISPORT"
          value = "6379"
        }

      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"        = "1000"
        "run.googleapis.com/cloudsql-instances"   = google_sql_database_instance.main.connection_name
        "run.googleapis.com/client-name"          = "terraform"
        "run.googleapis.com/vpc-access-egress"    = "all"
        "run.googleapis.com/vpc-access-connector" = google_vpc_access_connector.main.id
      }
    }
  }
  autogenerate_revision_name = true
  depends_on = [
    null_resource.cloudbuild_api,
  ]
}

resource "null_resource" "cloudbuild_fe" {
  provisioner "local-exec" {
    working_dir = "${path.module}/../code/frontend"
    command     = "gcloud builds submit . --substitutions=_REGION=${var.region},_BASENAME=${var.basename}"
  }

  depends_on = [
    google_artifact_registry_repository.todo_app,
    google_cloud_run_service.api
  ]
}

resource "google_cloud_run_service" "fe" {
  name     = "${var.basename}-fe"
  location = var.region
  project  = var.project_id

  template {
    spec {
      service_account_name = google_service_account.runsa.email
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.basename}-app/fe"
        ports {
          container_port = 80
        }
      }
    }
  }
  depends_on = [null_resource.cloudbuild_fe]
}

resource "google_cloud_run_service_iam_member" "noauth_api" {
  location = google_cloud_run_service.api.location
  project  = google_cloud_run_service.api.project
  service  = google_cloud_run_service.api.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_service_iam_member" "noauth_fe" {
  location = google_cloud_run_service.fe.location
  project  = google_cloud_run_service.fe.project
  service  = google_cloud_run_service.fe.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
