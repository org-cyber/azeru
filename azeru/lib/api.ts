const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface Document {
  id: number
  file_name: string
  file_size: number
  status: string
  chunk_count: number
  created_at: string
}

export interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

export interface ChatResponse {
  answer: string
  sources: any[]
}

// Upload PDF
export async function uploadPDF(file: File): Promise<{ document: Document }> {
  const formData = new FormData()
  formData.append('file', file)

  const res = await fetch(`${API_BASE}/api/upload`, {
    method: 'POST',
    body: formData,
  })

  if (!res.ok) throw new Error('Upload failed')
  return res.json()
}

// Get all documents
export async function getDocuments(): Promise<Document[]> {
  const res = await fetch(`${API_BASE}/api/documents`)
  if (!res.ok) throw new Error('Failed to fetch documents')
  return res.json()
}

// Send chat message
export async function sendChat(
  question: string, 
  apiKey: string
): Promise<ChatResponse> {
  const res = await fetch(`${API_BASE}/api/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ question, api_key: apiKey }),
  })

  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error || 'Chat failed')
  }
  return res.json()
}