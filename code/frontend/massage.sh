API=$(gcloud run services describe todo-api --region=us-central1 --format="value(status.url)")
stripped=$(echo ${API/https:\/\//})
sed -i"" -e "s/127.0.0.1:9000/$stripped/" www/js/main.js

