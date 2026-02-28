# ChromaDB Storage & Chat Issues - Root Cause Analysis

## Issue 1: Embeddings Not Being Saved to ChromaDB on Upload

### Root Cause:
**The handler is not logging response body errors from ChromaDB!**

In `chroma.go` line 143:
```go
if resp.StatusCode != 200 && resp.StatusCode != 201 {
    return fmt.Errorf("failed to add chunks: %d", resp.StatusCode)
}
```

This only returns the status code but **doesn't read the response body** which contains the actual error message from ChromaDB. So if ChromaDB rejects the request, you only see "failed to add chunks: 400" but no details.

### Likely Actual Problems:

1. **ChromaDB Collection Not Found**
   - Collection creation happens during `NewChromaClient()` initialization
   - If initialization fails silently, `c.enabled = false`
   - But code still tries to use the client
   - **The `IsHealthy()` method doesn't validate collection exists**

2. **Response Status Code Handling**
   - ChromaDB might return 204 (No Content) after successful add
   - Code waits for 200 or 201
   - **Missing error response body logging**

3. **StoreChunks is Called But Not Visible**
   - Because no detailed error logging
   - Silent failure with only "WARNING - Failed to store in ChromaDB"
   - No indication of what went wrong

---

## Issue 2: Chat Not Clearly Showing Whether ChromaDB Was Used

### Problems Found:

1. **Backend Chat Handler Issues**:
   - No logging to distinguish ChromaDB vs FTS5 search
   - Response doesn't indicate which search method was used
   - Frontend can't tell if vector search worked

2. **Frontend Issues**:
   - Chat page doesn't show search method in UI
   - No indicator that ChromaDB was queried
   - Sources are displayed but no indication if semantic or keyword

3. **Silent Fallback**:
   - If ChromaDB search fails, falls back to FTS5 with only log message
   - User sees results without knowing search method changed
   - Query embedding error is silently ignored

### Code Evidence:

**Backend** (handler.go line 216-226):
```go
if h.chroma != nil && h.chroma.IsHealthy() && h.ollama != nil && h.ollama.IsHealthy() {
    queryEmbedding, err := h.ollama.GetEmbedding(req.Question)
    if err == nil {  // ← If error here, silently continues without logging!
        searchResults, err = h.chroma.Search(queryEmbedding, 5)
        if err == nil && len(searchResults) > 0 {
            log.Printf("Chat: Found %d results via ChromaDB vector search", len(searchResults))
```

**The problem**: If `h.ollama.GetEmbedding()` fails, there's no error log, and search falls back to FTS5 without user knowing.

---

## Summary of Fixes Needed:

### Priority 1 (Critical for saving to ChromaDB):
1. Add detailed error response logging in StoreChunks
2. Make IsHealthy() validate collection existence
3. Handle all ChromaDB HTTP response codes correctly

### Priority 2 (Critical for chat visibility):
1. Add `search_method` to chat response (vector_search vs keyword_search)
2. Add detailed error logging when embedding generation fails
3. Display search method indicator in frontend

### Priority 3 (Quality of Life):
1. Add metrics/stats endpoint showing ChromaDB health
2. Add query timing metrics (ChromaDB vs FTS5)
3. Add frontend toast notifications for search method fallbacks
