# Define the gqlgen command
GQLGEN = go run github.com/99designs/gqlgen generate

# Default target
.PHONY: all
all: gqlgen

# Regenerate GraphQL files
.PHONY: gqlgen
gqlgen:
	@echo "🔄 Regenerating GraphQL code with gqlgen..."
	@$(GQLGEN)
	@echo "✅ GraphQL code regeneration complete!"
