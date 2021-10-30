SAACCOUNT=todo-developer-account
PROJECT=$(shell gcloud config get-value project)
PROJECTNUMBER=$(shell gcloud projects list --filter="$(PROJECT)" \
			--format="value(PROJECT_NUMBER)")
REDISNAME=todo-cache
SQLNAME=todo-db
REGION=us-central1		
SQLSERVICEACCOUNT=$(shell gcloud sql instances describe $(SQLNAME)  \
			--format="value(serviceAccountEmailAddress)")
REDISHOST=$(shell gcloud redis instances describe $(REDISNAME) --region $(REGION) --format="value(host)" )			
SQLHOST=$(shell gcloud sql instances describe $(SQLNAME) --format="value(ipAddresses[0].ipAddress)")

test:
	echo $(SQLSERVICEACCOUNT) 

env:
	gcloud config set project $(PROJECT)

services: env
	-gcloud services enable vpcaccess.googleapis.com
	-gcloud services enable cloudbuild.googleapis.com
	-gcloud services enable redis.googleapis.com
	-gcloud services enable run.googleapis.com
	-gcloud services enable sql-component.googleapis.com
	-gcloud services enable sqladmin.googleapis.com
	-gcloud services enable storage.googleapis.com
	-gcloud services enable secretmanager.googleapis.com


cloudbuild: env
	@echo ~~~ Enable Build service account to deploy to Cloud Run on $(PROJECT)
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/run.developer

memorystore: env
	-gcloud redis instances create $(REDISNAME) --size=1 --region=$(REGION)
	-gcloud compute networks vpc-access connectors create \
	$(REDISNAME)connector --network default --region $(REGION) \
	--range 10.8.0.0/28 	

perms:
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/vpaccess.user
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/iam.serviceAccountUser 
	-gcloud iam service-accounts add-iam-policy-binding $(PROJECTNUMBER)-compute@developer.gserviceaccount.com \
	--member=serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
	--role roles/iam.serviceAccountUser
	--region=us-central1
	-gcloud projects add-iam-policy-binding $(PROJECT) \
	--member=serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
	--role=roles/run.admin
	-gcloud iam service-accounts add-iam-policy-binding $(PROJECTNUMBER)-compute@developer.gserviceaccount.com \
	--member=serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
	--role roles/iam.serviceAccountUser
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/secretmanager.secretAccessor

cloudsql:
	gcloud beta sql instances create $(SQLNAME) --database-version=MYSQL_5_7 \
	--tier=db-f1-micro --region=us-central1 --root-password=password123	\
	--no-assign-ip --network=default

schema:
	-gsutil mb gs://$(PROJECT)-temp
	-gsutil cp ../code/database/schema.sql gs://$(PROJECT)-temp/schema.sql
	-gsutil iam ch serviceAccount:$(SQLSERVICEACCOUNT):objectViewer gs://$(PROJECT)-temp/
	-gcloud sql import sql $(SQLNAME) gs://$(PROJECT)-temp/schema.sql -q
	-gsutil rm gs://$(PROJECT)-temp/schema.sql
	-gsutil rb gs://$(PROJECT)-temp

sqlnetwork:
	-gcloud compute addresses create google-managed-services-$(REDISNAME)connector \
	--global --purpose=VPC_PEERING --prefix-length=16 --network=default

secrets:
	-gcloud secrets create sqlhost --replication-policy="automatic"
	echo -n $(SQLHOST) | xargs | gcloud secrets versions add sqlhost --data-file=-
	-gcloud secrets create redishost --replication-policy="automatic"
	echo -n $(REDISHOST) | xargs | gcloud secrets versions add redishost --data-file=-


prod: env services cloudbuild memorystore cloudsql schema perms sqlnetwork secrets 