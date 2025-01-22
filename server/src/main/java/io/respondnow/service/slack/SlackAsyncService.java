package io.respondnow.service.slack;

import java.util.concurrent.CompletableFuture;
import javax.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

@Service
public class SlackAsyncService {

  private static final Logger logger = LoggerFactory.getLogger(SlackAsyncService.class);
  private final SlackService slackService;

  public SlackAsyncService(SlackService slackService) {
    this.slackService = slackService;
  }

  @Async
  public CompletableFuture<Void> connectAndInitializeSlack() {
    return CompletableFuture.runAsync(
        () -> {
          try {
            slackService.startApp();
            slackService.setBotUserIDAndName();
            logger.info("Bot user ID and name set successfully");

            String incidentChannelID = slackService.getIncidentChannelID();
            String botUserId = slackService.getBotUserId();
            slackService.addBotUserToIncidentChannel(botUserId, incidentChannelID);
            logger.info("Bot user added to the incident channel successfully");

          } catch (Exception e) {
            logger.error("Failed to initialize Slack: {}", e.getMessage(), e);
          }
        });
  }

  @PostConstruct
  public void initializeSlackConnection() {
    connectAndInitializeSlack()
        .exceptionally(
            e -> {
              logger.error("Failed to initialize Slack connection: {}", e.getMessage(), e);
              return null;
            });
  }
}
