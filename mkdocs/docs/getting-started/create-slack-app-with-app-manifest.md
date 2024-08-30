# Create RespondNow Slack App 

- Click [here](https://api.slack.com/apps?new_app=1) to create a slack app

- Choose option to create an app _from a manifest_

- Select the desired slack workspace from the drop-down menu and click _Next_

- Paste the [RespondNow manifest configuration](https://github.com/respondnow/respond/blob/main/server/clients/slack/manifest.yaml) and click _Next_

- Review and verify that the configuration you entered matches the summary and click _Create_

- In the _Settings_ -> _Basic Information_ screen for the created app, generate an _App Level Token_ with the right scope (shown in the screenshot below) by clicking on _Generate Token and Scopes_.   

  <img width="515" alt="slack-app-token-generation" src="https://github.com/user-attachments/assets/4a4aa632-bf8a-4e17-96a6-f4a9e60b611f">

- Save the App token (beginning with xapp-*) for use in the subsequent steps

  <img width="430" alt="slac-app-token-copy-with-redact" src="https://github.com/user-attachments/assets/67595beb-a35d-4148-abae-559065f029e6">

- In the _Settings_ -> _Install App_ screen for the app, select the option to install to the desired slack workspace

  <img width="524" alt="slack-app-install-to-workspace" src="https://github.com/user-attachments/assets/07639b2a-2ac4-4a9f-b9fa-4134897c7159">

- Once installed, you will obtain User & Bot OAuth tokens. Copy the Bot User OAuth token for use in subsequent steps.

  <img width="585" alt="slack-app-bot-token-with-redact" src="https://github.com/user-attachments/assets/13314045-efc0-4e45-b771-8f8b3c3b555c"> 
