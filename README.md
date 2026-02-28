# Azeru - Enterprise Brain 🧠

Azeru is a cutting-edge, AI-driven document intelligence platform designed to transform your PDFs into actionable knowledge. By leveraging advanced Retrieval-Augmented Generation (RAG) techniques, Azeru provides precise answers based on your private documents, making it an indispensable tool for enterprises and professionals.

## 🚀 Features

- **PDF Intelligence**: Upload and process PDF documents seamlessly, extracting valuable insights.
- **Automated Text Extraction**: Intelligent chunking and text extraction for optimal context retrieval.
- **RAG-Powered Chat**: Engage in meaningful conversations with your documents using Llama-powered AI (via Groq).
- **Bring Your Own Key (BYOK)**: Ensure enhanced security and cost management by using your own Groq API keys.
- **Enterprise-Ready Architecture**: Built with a scalable Go backend and a modern Next.js frontend.
- **Customizable Search**: Supports both keyword-based and vector-based search methods for document queries.

## 🛠️ Technology Stack

### Backend

- **Language**: Go 1.25+
- **Framework**: [Gin](https://github.com/gin-gonic/gin) (HTTP Web Framework)
- **ORM**: [GORM](https://gorm.io/) with SQLite
- **AI Integration**: [Groq](https://groq.com/) (Llama 3/4 Models)
- **Database**: SQLite for lightweight and efficient data storage

### Frontend

- **Framework**: [Next.js](https://nextjs.org/) (React)
- **Language**: TypeScript
- **Styling**: CSS Modules for scoped and maintainable styles

## 📂 Project Structure

```text
azeru/
├── backend/            # Go Backend Service
│   ├── cmd/            # Application entry points
│   ├── internal/       # Private library code (handlers, services, models)
│   └── data/           # SQLite database storage
├── frontend/           # Next.js Frontend
│   └── app/            # Source code for the web app
└── architecture.md     # Detailed system architecture
```

## 🏁 Getting Started

### Prerequisites

To run Azeru, ensure you have the following installed:

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/) (recommended)
- Alternatively, install manually:
  - [Go](https://golang.org/dl/) (1.25 or later)
  - [Node.js](https://nodejs.org/) (v18 or later)
  - [Ollama](https://ollama.ai/) with the `nomic-embed-text` model
  - [ChromaDB](https://www.trychroma.com/) running locally

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repo/azeru.git
   cd azeru
   ```

2. **Start Services with Docker**:
   ```bash
   docker-compose up
   ```

3. **Manual Setup** (if not using Docker):
   - Start the backend:
     ```bash
     cd backend
     go run cmd/server/main.go
     ```
   - Start the frontend:
     ```bash
     cd frontend
     npm install
     npm run dev
     ```

4. **Access the Application**:
   Open your browser and navigate to `http://localhost:3000`.

### Usage

- **Upload Documents**: Navigate to the upload section and add your PDFs.
- **Chat with Documents**: Use the chat interface to ask questions and receive AI-generated answers.
- **Monitor Logs**: Check the backend logs for detailed processing information.

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with a detailed description of your changes.

## 📄 License

Azeru is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

For more information, visit our [documentation](https://your-docs-link.com) or contact us at support@azeru.com.
