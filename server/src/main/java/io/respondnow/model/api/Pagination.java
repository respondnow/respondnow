package io.respondnow.model.api;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@JsonInclude(JsonInclude.Include.NON_NULL)
@Builder
public class Pagination {

  @JsonProperty("all")
  private boolean all;

  @JsonProperty("index")
  private long index;

  @JsonProperty("limit")
  private long limit;

  @JsonProperty("totalPages")
  private long totalPages;

  @JsonProperty("totalItems")
  private long totalItems;

  public Pagination(boolean all, long index, long limit, long totalPages, long totalItems) {
    this.all = all;
    this.index = index;
    this.limit = limit;
    this.totalPages = totalPages;
    this.totalItems = totalItems;
  }
}
