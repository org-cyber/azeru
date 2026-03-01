# Contributing to Azeru

First off, thank you for considering contributing to Azeru! It is people like you that make this Enterprise Brain stronger and more powerful.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful, welcoming, and inclusive to all contributors. Harassment or abusive behavior will not be tolerated.

## How Can I Contribute?

### Reporting Bugs

If you find a bug, please create an issue on our repository. Include:

- A clear and descriptive title.
- Steps to reproduce the issue.
- Expected versus actual behavior.
- Error logs or screenshots if applicable.
- Environment details (OS, Docker version, Go version, Node version).

### Suggesting Enhancements

We welcome new feature requests! Please submit an issue detailing:

- The problem your feature solves.
- A proposed solution and how it should work.
- Any UI/UX mockups if available.

### Pull Requests

1. Fork the repository and create your branch from `main`.
2. Name your branch descriptively (e.g., `feature/add-new-search`, `bugfix/fix-upload-crash`).
3. Ensure your code follows the existing style and conventions.
4. Update or add relevant documentation (README, architecture maps) if applicable.
5. Create the pull request and describe your changes thoroughly. Link any relevant issues.

## Development Setup

To contribute to Azeru, you will need to set up the environment locally.

### Prerequisites

- Go (1.25 or later)
- Node.js (v18 or later)
- Docker & Docker Compose (for spinning up auxiliary services like DB, Chroma, Ollama)

### Backend Setup

1. Navigate to the `backend` directory.
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the development server:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup

1. Navigate to the `frontend` directory.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the Next.js development server:
   ```bash
   npm run dev
   ```

## Testing Protocol

We have integrated full vector-search and fallback flows that need careful testing when changing core components. Before submitting a PR, verify functionality manually using the steps from our testing guide.

### Environment Readiness

Ensure services are running (ChromaDB, Ollama, Backend, Frontend). Verify health checks:

```bash
curl http://localhost:8000/api/v1/heartbeat
curl http://localhost:11434/api/tags
curl http://localhost:8080/health
```

### 1. Document Upload Test

1. Navigate to the upload page and upload a small PDF.
2. Monitor the backend logs (`docker logs -f azeru-backend`).
3. Confirm lines indicating extraction, embedding generation, and successful storage in ChromaDB without `422` or mismatch errors.

### 2. Semantic Search Test

1. Navigate to the chat page, input your Groq API key, and ask a relevant question about the document.
2. Monitor backend logs to confirm a `vector_search` was performed.
3. Confirm the frontend returns an accurate answer and displays a "Vector Search (Semantic)" badge.

### 3. Fallback Search Test

1. Turn off ChromaDB (`docker stop azeru-chroma`).
2. Ask another question in the chat interface.
3. Verify backend logs indicate ChromaDB was unreachable and the system fell back to FTS5 full-text search.
4. Confirm the frontend returns a "Keyword Search (FTS5)" badge.
5. Restart ChromaDB after testing (`docker start azeru-chroma`).

By following these instructions, you help ensure Azeru remains stable, performant, and reliable for all users.
