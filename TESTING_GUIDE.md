# Quick Testing Guide - ChromaDB & Chat Fixes

## Pre-Test Checklist

```bash
# 1. Ensure services are running
docker-compose ps
# Should show: chroma, ollama, backend, frontend all running

# 2. Pull Ollama model (first time only)
docker exec azeru-ollama ollama pull nomic-embed-text

# 3. Verify services are healthy
curl http://localhost:8000/api/v1/heartbeat    # ChromaDB
curl http://localhost:11434/api/tags            # Ollama
curl http://localhost:8080/health               # Backend
```

---

## Test Scenario 1: Upload PDF & Check ChromaDB Storage

### Step 1: Upload PDF
- Go to http://localhost:3000
- Click "Upload Document"
- Select a PDF file
- Wait for processing to complete

### Step 2: Check Backend Logs
```bash
# In one terminal
docker logs -f azeru-backend
```

### Step 3: Look for Success Message
✅ **Good Log**:
```
PDF 1: Extracted 12 chunks
PDF 1: Generating embeddings for 12 chunks...
PDF 1: Successfully generated 12 embeddings
PDF 1: Stored 12 chunks in ChromaDB
PDF 1: Processing complete - 12 chunks indexed
```

❌ **Bad Log (Error Details Now Visible)**:
```
PDF 1: WARNING - Failed to store in ChromaDB: failed to add chunks: status=422, body={"code":"invalid_request","message":"Embedding dimension mismatch"}
```

---

## Test Scenario 2: Chat with ChromaDB Vector Search

### Step 1: Navigate to Chat
- Click "Chat" in the left sidebar
- Enter your Groq API key

### Step 2: Ask a Question
- Ask something related to the uploaded document
- Example: "What is this document about?"

### Step 3: Check Backend Logs
```bash
docker logs -f azeru-backend | grep "Chat:"
```

### Step 4: Look for Vector Search Success
✅ **Expected Log**:
```
Chat: Starting search for question: "What is this document about?"
Chat: ChromaDB healthy=true, Ollama healthy=true
Chat: Attempting ChromaDB vector search...
Chat: Generated query embedding, dimension=768
Chat: SUCCESS - Found 5 results via ChromaDB vector search
Chat: Successfully returned answer using vector_search
```

### Step 5: Check Frontend UI
- Look for the **blue badge** below the answer:
  ```
  🔍 Vector Search (Semantic)
  ```

---

## Test Scenario 3: Fallback to FTS5 (Optional - Stop ChromaDB)

### Step 1: Stop ChromaDB
```bash
docker stop azeru-chroma
```

### Step 2: Ask a Question
- In chat, ask another question
- Same question format as before

### Step 3: Check Backend Logs
```bash
docker logs -f azeru-backend | grep "Chat:"
```

### Step 4: Look for Fallback Message
✅ **Expected Log**:
```
Chat: Starting search for question: "Why is this important?"
Chat: ChromaDB healthy=false, Ollama healthy=true
Chat: ChromaDB/Ollama not available, skipping vector search
Chat: Falling back to FTS5 full-text search
Chat: FTS5 search returned 5 results
Chat: Successfully returned answer using keyword_search
```

### Step 5: Check Frontend UI
- Look for the **blue badge** below the answer:
  ```
  🔍 Keyword Search (FTS5)
  ```

### Step 6: Restart ChromaDB
```bash
docker start azeru-chroma
```

---

## Debugging Commands

### View All Logs
```bash
# Backend logs with timestamp
docker logs -f azeru-backend --tail=100

# ChromaDB logs
docker logs -f azeru-chroma

# Ollama logs
docker logs -f azeru-ollama
```

### Check ChromaDB Collection
```bash
# List collections
curl http://localhost:8000/api/v1/collections

# Get specific collection
curl http://localhost:8000/api/v1/collections/enterprise_brain

# Count items in collection
curl http://localhost:8000/api/v1/collections/enterprise_brain/count
```

### Check Ollama Models
```bash
# List available models
curl http://localhost:11434/api/tags

# For docker
docker exec azeru-ollama ollama list
```

### Test Backend API Directly
```bash
# Test chat endpoint (requires GROQ_API_KEY)
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "question": "What is this about?",
    "api_key": "gsk_YOUR_KEY_HERE"
  }'

# Should return JSON with:
# - answer: string
# - sources: array
# - search_method: "vector_search" | "keyword_search" | "none"
```

---

## Expected Outputs

### Successful Vector Search Response
```json
{
  "answer": "Based on the documents, this is about...",
  "sources": [
    {
      "id": 1,
      "document_id": 1,
      "content": "...",
      "page_num": 0
    }
  ],
  "search_method": "vector_search"
}
```

### Fallback FTS5 Response
```json
{
  "answer": "Based on the documents, this is about...",
  "sources": [
    {
      "id": 5,
      "document_id": 1,
      "content": "...",
      "page_num": 2
    }
  ],
  "search_method": "keyword_search"
}
```

### No Results Response
```json
{
  "answer": "I couldn't find any relevant information in the uploaded documents.",
  "sources": [],
  "search_method": "none"
}
```

---

## Troubleshooting

### Problem: "ChromaDB not healthy"
```bash
# Check if ChromaDB is running
docker ps | grep chroma

# Check health endpoint
curl http://localhost:8000/api/v1/heartbeat

# View logs
docker logs azeru-chroma

# Restart if needed
docker restart azeru-chroma
```

### Problem: "Ollama not reachable"
```bash
# Check if Ollama is running
docker ps | grep ollama

# Check available models
docker exec azeru-ollama ollama list

# Ensure model is pulled
docker exec azeru-ollama ollama pull nomic-embed-text

# View logs
docker logs azeru-ollama
```

### Problem: "Failed to add chunks: 422"
This usually means embedding dimension mismatch. Check:
1. Is Ollama using the same model? (should be `nomic-embed-text`)
2. Are embeddings being generated with correct dimensions?
3. Check Ollama logs: `docker logs azeru-ollama`

### Problem: Chat only returns keyword search results
Check backend logs to see why ChromaDB isn't being used:
```bash
docker logs azeru-backend | grep "Chat: ERROR"
```

---

## Performance Baseline

### Expected Times
- **PDF Upload** (10KB): 1-3 seconds
- **Embedding Generation** (100 chunks): 10-30 seconds
- **Vector Search**: 100-500ms
- **FTS5 Search**: 50-100ms
- **AI Response**: 3-10 seconds

### Monitor Performance
```bash
# Watch for timing in logs
grep "ChatTime\|EmbeddingTime\|SearchTime" container.log
```

---

## Success Criteria

✅ **All Tests Pass When**:
1. PDF uploads generate embeddings without errors
2. ChromaDB receives and stores embeddings
3. Chat shows "Vector Search (Semantic)" badge
4. Backend logs show successful ChromaDB operations
5. Fallback to FTS5 works when ChromaDB is unavailable
6. Frontend displays search method indicator

❌ **Issues If**:
1. Logs show "failed to add chunks" with no error details
2. Chat doesn't show search method badge
3. Only seeing FTS5 searches even though ChromaDB is running
4. Silent failures in ChromaDB integration
