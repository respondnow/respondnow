package io.respondnow.service.api;

import io.respondnow.model.api.Pagination;
import org.springframework.stereotype.Service;

@Service
public class PaginationService {

  public Pagination getPagination(long page, long limit, long total, boolean all) {
    if (all || limit == 0) {
      return new Pagination(
          all, 0, 0, 0, 0); // Return a pagination object where no limit is applied
    }

    long remainder = total % limit;
    long pageCount = total / limit;

    if (remainder != 0) {
      pageCount += 1;
    }

    return new Pagination(all, page, limit, pageCount, total);
  }
}
