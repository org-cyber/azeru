# ChromaDB Storage & Chat - Issues Fixed

## Issues Identified & Resolved

### Issue #1: Embeddings Not Being Saved to ChromaDB

**Root Cause**: 
The `StoreChunks()` method in `chroma.go` was returning HTTP status codes without reading the response body. When ChromaDB returned an error (e.g., 400, 422), only the status code was logged, not the actual error message explaining what went wrong.

**Lines 143 (OLD)**:
```go
if resp.StatusCode != 200 && resp.StatusCode != 201 {
    return fmt.Errorf("failed to add chunks: %d", resp.StatusCode)  // Only shows "400", not why
}
```

**Fix Applied** ✅:
1. Modified `StoreChunks()` to read response body and include it in error message
2. Enhanced `IsHealthy()` to validate collection actually exists (not just that client initialized)
3. Added detailed error logging in upload handler
4. Import added: `"io"` for `io.ReadAll()`

**Lines 143-155 (NEW)**:
```go
// Read response body for error details
respBody, _ := io.ReadAll(resp.Body)

if resp.StatusCode != 200 && resp.StatusCode != 201 {
    errMsg := string(respBody)
    if errMsg == "" {
        errMsg = "no error details provided"
    }
    return fmt.Errorf("failed to add chunks: status=%d, body=%s", resp.StatusCode, errMsg)
}
```

**Result**: Now you'll see actual error messages like:
```
failed to add chunks: status=422, body={"error": "Invalid embedding dimension"}
```

---

### Issue #2: Chat Not Showing Whether ChromaDB Was Used

**Root Cause**:
Backend silently fell back to FTS5 if ChromaDB search failed, with minimal logging. Frontend had no way to display which search method was used.

**Fixed With**:

#### Backend Changes (handler.go):

1. **Added detailed logging throughout Chat handler**:
   - Logs when ChromaDB/Ollama availability checked
   - Logs if embedding generation fails
   - Logs if ChromaDB search returns 0 results
   - Logs fallback to FTS5
   - Reports final search method used

2. **Added `search_method` field to response**:
   ```go
   c.JSON(http.StatusOK, gin.H{
       "answer":        answer,
       "sources":       chunks,
       "search_method": searchMethod,  // "vector_search", "keyword_search", or "none"
   })
   ```

3. **Example log output for debugging**:
   ```
   Chat: Starting search for question: "What is the summary?"
   Chat: ChromaDB healthy=true, Ollama healthy=true
   Chat: Attempting ChromaDB vector search...
   Chat: Generated query embedding, dimension=768
   Chat: SUCCESS - Found 5 results via ChromaDB vector search
   Chat: Successfully returned answer using vector_search
   ```

#### Frontend Changes (chat/page.tsx):

1. **Updated Message interface** to include search method:
   ```typescript
   interface Message {
     searchMethod?: 'vector_search' | 'keyword_search' | 'none'
   }
   ```

2. **Capture search method from API response**:
   ```typescript
   const assistantMessage: Message = {
     searchMethod: response.search_method || 'none',
   }
   ```

3. **Display search method badge**:
   - Shows "Vector Search (Semantic)" in blue badge if ChromaDB was used
   - Shows "Keyword Search (FTS5)" if keyword search fallback occurred
   - Badge displays below the answer text

#### CSS Styling (chat.module.css):

Added `.searchMethodBadge` class:
```css
.searchMethodBadge {
  margin-top: 10px;
  padding: 6px 10px;
  background: #f0f9ff;
  border: 1px solid #bfdbfe;
  border-radius: 4px;
  font-size: 0.75rem;
  color: #1e40af;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
```

---

## How to Debug Now

### When Embeddings Aren't Saving to ChromaDB:

1. **Check backend logs** - Now shows detailed error:
   ```
   PDF 1: WARNING - Failed to store in ChromaDB: failed to add chunks: status=422, body={"error": "..."}
   ```

2. **Check ChromaDB health** - IsHealthy() now verifies collection:
   ```
   docker exec azeru-chroma curl http://localhost:8000/api/v1/collections/enterprise_brain
   ```

3. **Verify Ollama embeddings generated** - Look for:
   ```
   PDF 1: Successfully generated 512 embeddings
   ```

### When Chat Isn't Using ChromaDB:

1. **Check browser** - Look for search method badge showing:
   - 🔵 "Vector Search (Semantic)" = ChromaDB was used ✅
   - 🔵 "Keyword Search (FTS5)" = Fallback occurred ⚠️

2. **Check backend logs** - Shows search path:
   ```
   Chat: Attempting ChromaDB vector search...
   Chat: ERROR in ChromaDB search: connection refused
   Chat: Falling back to FTS5 full-text search
   ```

3. **Verify services are healthy**:
   ```bash
   # Check ChromaDB
   curl http://localhost:8000/api/v1/heartbeat
   
   # Check Ollama  
   curl http://localhost:11434/api/tags
   
   # Check Backend
   curl http://localhost:8080/health
   ```

---

## Testing Steps

### Test 1: Verify Embeddings Save to ChromaDB

1. Upload a PDF through frontend
2. Check backend logs for "Stored X chunks in ChromaDB"
3. If not appearing, logs now show the actual error

### Test 2: Verify Chat Uses ChromaDB

1. Ask a question in chat
2. Look for blue badge showing "Vector Search (Semantic)"
3. If fallback to FTS5, badge shows "Keyword Search (FTS5)"
4. Check backend console for search logs

### Test 3: Simulate ChromaDB Failure

1. Stop ChromaDB: `docker stop azeru-chroma`
2. Ask a question
3. Should see "Keyword Search (FTS5)" badge
4. Backend logs show fallback reasoning

---

## Files Modified

| File | Changes |
|------|---------|
| `backend/internal/services/chroma.go` | Added `io` import, enhanced error logging in StoreChunks(), improved IsHealthy() to validate collection |
| `backend/internal/handlers/handler.go` | Added detailed logging throughout Chat handler, added search_method to response |
| `azeru/app/chat/page.tsx` | Added searchMethod to Message type, capture from response, render badge |
| `azeru/app/css/chat.module.css` | Added .searchMethodBadge styling |

---

## What Information is Now Visible

### In Backend Logs:
- ✅ Detailed ChromaDB error messages (not just status codes)
- ✅ Ollama embedding generation success/failure
- ✅ ChromaDB search success/failure
- ✅ Fallback reasoning
- ✅ Which search method was actually used

### In Frontend UI:
- ✅ Visual indicator of search method (Vector vs Keyword)
- ✅ Number of source documents used
- ✅ Color-coded badge (blue) for easy identification

### In API Response:
- ✅ `search_method` field: "vector_search" | "keyword_search" | "none"
- ✅ Allows clients to react differently based on search quality

---

## Next Steps (Optional Enhancements)

1. **Add timing metrics**: Log how long vector search took vs FTS5
2. **Add fallback notifications**: Toast notification when falling back to FTS5
3. **Add ChromaDB stats endpoint**: Return health metrics via `/api/health/detailed`
4. **Add search confidence score**: Include similarity scores from vector search
5. **Add user preference**: Let users choose search method (vector-only vs hybrid)
