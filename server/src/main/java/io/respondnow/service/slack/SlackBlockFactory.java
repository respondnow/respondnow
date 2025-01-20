package io.respondnow.service.slack;

import com.slack.api.model.block.*;
import com.slack.api.model.block.composition.MarkdownTextObject;
import com.slack.api.model.block.composition.PlainTextObject;
import com.slack.api.model.block.element.ButtonElement;

public class SlackBlockFactory {

  // Create a header block with a given text
  public static HeaderBlock createHeaderBlock(String text, String blockId) {
    return HeaderBlock.builder()
        .blockId(blockId)
        .text(PlainTextObject.builder().text(text).build())
        .build();
  }

  // Create an actions block with a button
  public static ActionsBlock createActionsBlock(String blockId, ButtonElement button) {
    return ActionsBlock.builder()
        .blockId(blockId)
        .elements(java.util.Collections.singletonList(button))
        .build();
  }

  // Create a divider block
  public static DividerBlock createDividerBlock() {
    return DividerBlock.builder().build();
  }

  // Create a section block with a given markdown text
  public static SectionBlock createSectionBlock(String markdownText, String blockId) {
    return SectionBlock.builder()
        .blockId(blockId)
        .text(MarkdownTextObject.builder().text(markdownText).build())
        .build();
  }

  // Create an image block with an image URL and alt text
  public static ImageBlock createImageBlock(
      String imageUrl, String altText, String blockId, String caption) {
    return ImageBlock.builder()
        .imageUrl(imageUrl)
        .altText(altText)
        .title(PlainTextObject.builder().text(caption).build())
        .blockId(blockId)
        .build();
  }
}
