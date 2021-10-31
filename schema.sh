PROJECT=$(gcloud config get-value project | xargs)
PROJECTNUMBER=$(gcloud projects list --filter="$PROJECT" --format="value(PROJECT_NUMBER)" | xargs)		
SQLNAME=$PROJECT-todo-db



printf "******************************************************************************** \n"
printf "Populating SQL Schema and loading starting data \n"
SQLSERVICEACCOUNT=$(gcloud sql instances describe $SQLNAME --format="value(serviceAccountEmailAddress)" | xargs)
gsutil mb gs://$PROJECT-temp 
gsutil cp code/database/schema.sql gs://$PROJECT-temp/schema.sql
echo $SQLSERVICEACCOUNT
gsutil iam ch serviceAccount:$SQLSERVICEACCOUNT:objectViewer gs://$PROJECT-temp/
gcloud sql import sql $SQLNAME gs://$PROJECT-temp/schema.sql -q
gsutil rm gs://$PROJECT-temp/schema.sql
gsutil rb gs://$PROJECT-temp
printf "Populating SQL Schema and loading starting data - done \n"


