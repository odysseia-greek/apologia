# Paths
SOKRATES_SCHEMA_SRC := ../sokrates/graph/model/models_gen.go
SOKRATES_SCHEMA_DEST := ./model/models_gen.go

# Colors for logging
BLUE := \033[1;34m
GREEN := \033[1;32m
YELLOW := \033[1;33m
NC := \033[0m  # No Color

# Copy Sokrates schema only if needed
sokrates_schema:
	@if cmp -s $(SOKRATES_SCHEMA_SRC) $(SOKRATES_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in Sokrates schema, skipping copy."; \
	else \
		cp $(SOKRATES_SCHEMA_SRC) $(SOKRATES_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Sokrates schema copied successfully!"; \
	fi
