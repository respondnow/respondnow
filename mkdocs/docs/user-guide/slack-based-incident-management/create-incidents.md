# Create Incidents

To create a new incident, you can do the following:

- Start from incident slack channel

  - Type /start in the designated incident slack channel (identified & configured during the RespondNow server setup), 
or in any slack channel where the RespondNow bot has been integrated. 

    <img width="1309" alt="slack-incident-shortcut" src="https://github.com/user-attachments/assets/12037872-620c-46a0-a921-13e276fc3b08">


- Start from slack app home page 

  - Click the _Start New Incident_ button on the slack

    <img width="973" alt="slack-app-home-page" src="https://github.com/user-attachments/assets/ac2f302a-a2a5-4fe0-80ea-a55c142c1175">

- Provide the following inputs on the Incident creation modal: 

  - Incident Name: A meaningful name indicating the issue. You could follow custom naming conventions as per your needs
  - Incident Summary: Provide a succict summary outlining the issue and its impact
  - Severity: Marker to indicate the impact of the incident and urgency with which it is expected to be resolved
  - Role: Self-assign a specific responsibility that you, as incident creator, would be carrying out during the incident lifecycle      
  - Incident Channel: Slack channel where the incident details are placed and tracked

  <img width="514" alt="incident-creation-modal" src="https://github.com/user-attachments/assets/63512064-c1c2-4137-94f8-6a7c22bb3fe4">

  Of the described inputs, the Incident Name and Summary are expected to be filled out by the user, while the rest have default settings
  which can be modified as needed.

- Once created, the incident details are placed in the incident slack channel, along with options to update its attributes

  <img width="733" alt="incident-details-on-incident-channel" src="https://github.com/user-attachments/assets/7ee76afe-c588-41e2-8de7-023d22045a8f">

- A new slack channel following the naming convention `rn-{timestamp}-{incident-name}` is created as the designated war-room or triage-channel 
  for discussions about the incident, its RCA and eventual resolution. The incident creator is added to this channel by default and is expected  to pull the required team members

  <img width="1645" alt="incident-war-room-creation" src="https://github.com/user-attachments/assets/6d003663-c2c6-4856-86a1-58f6773df4eb">

  
