# Apologia – The Sokratic Method

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## Overview
Apologia is the central hub for all quiz-related services in Odysseia-Greek, uniting multiple quiz modes through the Sokrates GraphQL API. Each quiz mode is managed by a dedicated backend service, named after a pupil of Sokrates, providing a comprehensive Greek language learning experience.

## Architecture
Apologia follows a microservices architecture where each service is responsible for a specific quiz type. The Sokrates service acts as the central GraphQL API gateway that communicates with all the backend services via gRPC.

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│                        Client Applications                      │
│                                                                 │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│                  Sokrates (GraphQL API Gateway)                 │
│                                                                 │
└───┬───────────┬───────────┬───────────┬───────────┬───────────┬─┘
    │           │           │           │           │           │
    ▼           ▼           ▼           ▼           ▼           ▼
┌──────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────────┐ ┌──────────┐
│          │ │         │ │         │ │         │ │           │ │          │
│Aristippos│ │ Kritias │ │ Xenofon │ │ Kriton  │ │Antisthenes│ │Alkibiades│
│ (Media)  │ │(Multiple│ │ (Author │ │(Dialogue│ │(Grammar)  │ │(Journey) │
│          │ │ Choice) │ │  Based) │ │  Based) │ │           │ │          │
└──────────┘ └─────────┘ └─────────┘ └─────────┘ └───────────┘ └──────────┘
```

## Services

### Sokrates – The Central API
Sokrates serves as the GraphQL API gateway that unifies all quiz modes. It handles client requests, routes them to the appropriate backend service, and returns the responses. The GraphQL API provides a flexible and efficient way for clients to query exactly the data they need.

**Key Features:**
- GraphQL API for all quiz types
- Session management
- Request routing to appropriate backends
- Response transformation and aggregation

### Quiz Backends

#### 🎭 Aristippos – Ἀρίστιππος - Media-based Quizzes
Aristippos handles media-based quizzes that incorporate audio and images to enhance the learning experience.

**API Endpoints:**
- `MediaQuiz`: Retrieves a new media-based quiz
- `MediaAnswer`: Validates answers for media quizzes
- `MediaOptions`: Provides available options for media quizzes

#### ❓ Kritias – Κριτίας - Multiple Choice Quizzes
Kritias manages multiple-choice quizzes, offering a straightforward way to test knowledge of Greek vocabulary and concepts.

**API Endpoints:**
- `MultipleChoiceQuiz`: Retrieves a new multiple-choice quiz
- `MultipleChoiceAnswer`: Validates answers for multiple-choice quizzes
- `MultipleChoiceOptions`: Provides available options for multiple-choice quizzes

#### 📜 Xenofon – Ξενοφῶν - Author-based Quizzes
Xenofon specializes in quizzes based on works by ancient Greek authors, helping learners engage with authentic texts.

**API Endpoints:**
- `AuthorBasedQuiz`: Retrieves a new author-based quiz
- `AuthorBasedAnswer`: Validates answers for author-based quizzes
- `AuthorBasedOptions`: Provides available options for author-based quizzes
- `AuthorBasedWordForms`: Retrieves word forms for author-based quizzes

#### 🗣️ Kriton - Κρίτων – Dialogue-based Quizzes
Kriton focuses on dialogue-based quizzes that simulate conversations in ancient Greek, improving conversational skills.

**API Endpoints:**
- `DialogueQuiz`: Retrieves a new dialogue-based quiz
- `DialogueAnswer`: Validates answers for dialogue-based quizzes
- `DialogueOptions`: Provides available options for dialogue-based quizzes

#### 🔠 Antisthenes – Ἀντισθένης - Grammar-based Quizzes
Antisthenes handles grammar-based quizzes that focus on Greek grammar rules and structures.

**API Endpoints:**
- `GrammarQuiz`: Retrieves a new grammar-based quiz
- `GrammarAnswer`: Validates answers for grammar-based quizzes
- `GrammarOptions`: Provides available options for grammar-based quizzes

#### 🛤️ Alkibiades – Ἀλκιβιάδης - Journey Mode
Alkibiades manages the Journey mode, which combines multiple quiz types into a cohesive learning experience. It creates a narrative-driven path through different quiz types, offering a comprehensive approach to learning Greek.

**API Endpoints:**
- `JourneyQuiz`: Retrieves a new journey segment with various quiz types
- `JourneyOptions`: Provides available journey themes and segments

#### 🧪 Meletos – Μέλητος - Testing Framework
Meletos is the testing framework for the Apologia system. It uses behavior-driven development (BDD) with Cucumber/Godog to test all services and ensure they function correctly. Named after Meletos, one of Sokrates' accusers, this service ensures the quality and reliability of the entire system.

**Key Features:**
- Integration tests for all quiz services
- Health check tests
- BDD-style test scenarios
- Comprehensive test coverage for all API endpoints

## Data Flow

1. **Client Request**: A client application sends a GraphQL query to the Sokrates API.
2. **Request Routing**: Sokrates identifies the quiz type and routes the request to the appropriate backend service.
3. **Backend Processing**: The backend service processes the request and returns a response.
4. **Response Transformation**: Sokrates transforms the gRPC response into the GraphQL format.
5. **Client Response**: The formatted response is sent back to the client.

## Installation

### Prerequisites
- Go 1.16 or higher
- Docker (for building container images)
- Kubernetes (for local or production deployment)
- Skaffold (for streamlined Kubernetes deployment)

### Local Development Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/odysseia-greek/apologia.git
   cd apologia
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start the services using skaffold:
   ```bash
   skaffold dev
   ```

   This will build all the necessary images and deploy them to your local Kubernetes cluster.

4. Access the GraphQL playground at `http://localhost:8080/playground`

### Skaffold Deployment
1. Ensure you have a local Kubernetes cluster running (e.g., k3d, minikube, Docker Desktop Kubernetes)

2. Deploy the application:
   ```bash
   skaffold run
   ```

   This will build the images and deploy them to your Kubernetes cluster.

3. For development with hot-reload:
   ```bash
   skaffold dev
   ```

   This will watch for changes in your code and automatically rebuild and redeploy.

### Kubernetes Deployment
The project uses Skaffold to manage Kubernetes deployments. The `skaffold.yaml` file in the root directory configures:

1. Building Docker images for all services
2. Deploying Helm charts for each service
3. Setting up appropriate configurations for local development

To customize the Kubernetes deployment:
1. Modify the `skaffold.yaml` file to change build or deployment settings
2. Update the Helm values files referenced in the skaffold configuration
3. Use `skaffold run --profile <profile-name>` for different deployment profiles

## Development

### Project Structure
- `/sokrates` - GraphQL API gateway
- `/aristippos` - Media-based quiz service
- `/kritias` - Multiple choice quiz service
- `/xenofon` - Author-based quiz service
- `/kriton` - Dialogue-based quiz service
- `/antisthenes` - Grammar-based quiz service
- `/alkibiades` - Journey mode service
- `/meletos` - Testing framework
- `/parmenides` - Data seeder

### Adding a New Quiz Type
1. Create a new service directory
2. Implement the gRPC service interface
3. Add the service to the Sokrates gateway
4. Update the GraphQL schema
5. Add tests in the Meletos framework

## Testing

The Meletos service provides comprehensive testing for all Apologia services:

1. Run all tests:
   ```bash
   cd meletos
   go test ./...
   ```

2. Run specific feature tests:
   ```bash
   cd meletos
   go test -tags=@grammar
   ```

3. Run tests with verbose output:
   ```bash
   cd meletos
   go test -v -tags=@dialogue
   ```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.

## Acknowledgments

- The Odysseia Greek project team
- Contributors to the Go ecosystem
- Ancient Greek scholars and their timeless wisdom
