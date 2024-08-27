export function generateSlackChannelLink(workspaceDomain: string, channelId: string): string {
  return `https://${workspaceDomain}.slack.com/archives/${channelId}`;
}
