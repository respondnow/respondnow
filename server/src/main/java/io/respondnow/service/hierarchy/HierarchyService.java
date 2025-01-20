package io.respondnow.service.hierarchy;

import java.util.concurrent.CompletableFuture;

public interface HierarchyService {
  CompletableFuture<Void> backgroundProcess();

  void startProcess();
}
