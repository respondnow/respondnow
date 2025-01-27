# Deploy RespondNow Server

- Add the RespondNow Helm Repository

  `helm repo add respondnow https://respondnow.github.io/respondnow-helm`

  `helm repo update`

- Install RespondNow server by providing the slack App and Bot tokens noted in the previous steps

  `helm install respondnow respondnow/respondnow --namespace=respondnow --create-namespace --set server.configMap.data.ENABLE_SLACK_CLIENT=true --set server.configMap.data.INCIDENT_CHANNEL_ID="repond-now" --set server.secret.data.SLACK_APP_TOKEN="FILL-YOUR-SLACK-APP-TOKEN" --set server.secret.data.SLACK_BOT_TOKEN="FILL-YOUR-SLACK-BOT-TOKEN"`

- Verify that all pods in the `respondnow` namespace are up and running successfully

- Now, you are all setup to begin your incident management journey on slack. Refer to the [User Guide](https://respondnow.github.io/respondnow/user-guide/slack-based-incident-management/create-incidents/) to create your first incident. 


