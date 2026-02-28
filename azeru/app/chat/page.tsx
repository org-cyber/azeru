'use client'

import { useState, useRef, useEffect } from 'react'
import Link from 'next/link'
import { sendChat } from '@/lib/api'
import styles from '../css/chat.module.css'

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
  sources?: any[]
  searchMethod?: 'vector_search' | 'keyword_search' | 'none'
}

export default function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [apiKey, setApiKey] = useState('')
  const [showSettings, setShowSettings] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const saved = localStorage.getItem('enterprise_brain_api_key')
    if (saved) setApiKey(saved)
  }, [])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || !apiKey.trim()) {
      if (!apiKey.trim()) setShowSettings(true)
      return
    }

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input.trim(),
      timestamp: new Date(),
    }

    setMessages(prev => [...prev, userMessage])
    setInput('')
    setIsLoading(true)

    try {
      const response = await sendChat(userMessage.content, apiKey)

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response.answer,
        timestamp: new Date(),
        sources: response.sources,
        searchMethod: (response as any).search_method || 'none',
      }

      setMessages(prev => [...prev, assistantMessage])
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: `Error: ${error instanceof Error ? error.message : 'Failed to get response'}`,
        timestamp: new Date(),
      }
      setMessages(prev => [...prev, errorMessage])
    } finally {
      setIsLoading(false)
    }
  }

  const saveApiKey = () => {
    localStorage.setItem('enterprise_brain_api_key', apiKey)
    setShowSettings(false)
  }

  const clearChat = () => {
    if (confirm('Clear all messages?')) {
      setMessages([])
    }
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <Link href="/" className={styles.backLink}>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            <span>Back</span>
          </Link>

          <div className={styles.headerTitle}>
            <div className={styles.titleContainer}>
              <h1>Knowledge Assistant</h1>
              <p>Enterprise Intelligence Platform</p>
            </div>
          </div>

          <div className={styles.headerActions}>
            <button
              onClick={() => setShowSettings(!showSettings)}
              className={`${styles.headerButton} ${showSettings ? styles.active : ''}`}
              title="Settings"
              aria-label="Settings"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="3"/>
                <path d="M12 1v6m0 6v6M4.22 4.22l4.24 4.24m4.24 4.24l4.24 4.24M1 12h6m6 0h6M4.22 19.78l4.24-4.24m4.24-4.24l4.24-4.24"/>
              </svg>
            </button>
            <button
              onClick={clearChat}
              className={styles.headerButton}
              title="Clear conversation"
              aria-label="Clear conversation"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
                <line x1="10" y1="11" x2="10" y2="17"/>
                <line x1="14" y1="11" x2="14" y2="17"/>
              </svg>
            </button>
          </div>
        </div>

        {/* Settings Panel - Integrated */}
        {showSettings && (
          <div className={styles.settingsPanel}>
            <div className={styles.settingsContent}>
              <div className={styles.settingsHeader}>
                <h3>API Configuration</h3>
                <button
                  onClick={() => setShowSettings(false)}
                  className={styles.closeButton}
                  aria-label="Close settings"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                    <line x1="18" y1="6" x2="6" y2="18"/>
                    <line x1="6" y1="6" x2="18" y2="18"/>
                  </svg>
                </button>
              </div>
              <p className={styles.settingsDescription}>
                Enter your Groq API key. Your credentials are stored locally and never transmitted externally.
              </p>
              <div className={styles.apiKeyInput}>
                <input
                  type="password"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="gsk_..."
                  aria-label="API Key"
                  className={styles.input}
                />
                <button onClick={saveApiKey} className={styles.saveButton}>
                  Save Configuration
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Warning Banner */}
        {!apiKey && !showSettings && (
          <div className={styles.warningBanner}>
            <div className={styles.warningContent}>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              <span>API key required to begin</span>
              <button onClick={() => setShowSettings(true)} className={styles.warningButton}>
                Configure
              </button>
            </div>
          </div>
        )}
      </header>

      {/* Chat Container */}
      <div className={styles.chatContainer}>
        {messages.length === 0 ? (
          <div className={styles.emptyState}>
            <div className={styles.emptyIcon}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
            </div>
            <h2>Start a Conversation</h2>
            <p>Ask questions about your indexed documents and receive intelligent insights powered by semantic search.</p>
            
            <div className={styles.suggestionsGrid}>
              <button 
                className={styles.suggestionCard}
                onClick={() => setInput('What is this document about?')}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2m0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8m3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z"/>
                </svg>
                <span>Summarize Documents</span>
              </button>
              
              <button 
                className={styles.suggestionCard}
                onClick={() => setInput('What are the main findings?')}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M3 3h18v2H3V3zm0 4h18v2H3V7zm0 4h18v2H3v-2zm0 4h18v2H3v-2z"/>
                </svg>
                <span>Key Findings</span>
              </button>

              <button 
                className={styles.suggestionCard}
                onClick={() => setInput('Compare and analyze the content')}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M9 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h4m0-18v18m0-18h10c1.1 0 2 .9 2 2v14c0 1.1-.9 2-2 2h-10"/>
                </svg>
                <span>Data Analysis</span>
              </button>

              <button 
                className={styles.suggestionCard}
                onClick={() => setInput('Extract specific information')}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M21 21H3V3h9V1H3a2 2 0 0 0-2 2v18a2 2 0 0 0 2 2h18a2 2 0 0 0 2-2v-9h-2v9z"/>
                  <path d="M19 1v5h-5"/>
                  <path d="M20 1l-9.2 9.2"/>
                </svg>
                <span>Extract Information</span>
              </button>
            </div>
          </div>
        ) : (
          <div className={styles.messagesWrapper}>
            {messages.map((msg) => (
              <div
                key={msg.id}
                className={`${styles.messageGroup} ${styles[msg.role]}`}
              >
                <div className={styles.messageRow}>
                  {msg.role === 'assistant' && (
                    <div className={styles.avatar}>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                        <circle cx="12" cy="12" r="10"/>
                        <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
                        <circle cx="9" cy="9" r="1"/>
                        <circle cx="15" cy="9" r="1"/>
                      </svg>
                    </div>
                  )}

                  <div className={styles.messageBubble}>
                    <p>{msg.content}</p>
                    
                    {msg.searchMethod && msg.searchMethod !== 'none' && (
                      <div className={styles.methodTag}>
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <circle cx="11" cy="11" r="8"/>
                          <path d="m21 21-4.35-4.35"/>
                        </svg>
                        <span>
                          {msg.searchMethod === 'vector_search' ? 'Semantic Search' : 'Keyword Search'}
                        </span>
                      </div>
                    )}

                    {msg.sources && msg.sources.length > 0 && (
                      <div className={styles.sourceTag}>
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                        </svg>
                        <span>{msg.sources.length} source{msg.sources.length !== 1 ? 's' : ''}</span>
                      </div>
                    )}
                  </div>

                  {msg.role === 'user' && (
                    <div className={styles.avatar}>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                        <circle cx="12" cy="7" r="4"/>
                      </svg>
                    </div>
                  )}
                </div>
                <span className={styles.timestamp}>
                  {msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </span>
              </div>
            ))}

            {isLoading && (
              <div className={`${styles.messageGroup} ${styles.assistant}`}>
                <div className={styles.messageRow}>
                  <div className={styles.avatar}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                      <circle cx="12" cy="12" r="10"/>
                      <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
                      <circle cx="9" cy="9" r="1"/>
                      <circle cx="15" cy="9" r="1"/>
                    </svg>
                  </div>
                  <div className={styles.messageBubble}>
                    <div className={styles.typingIndicator}>
                      <span></span>
                      <span></span>
                      <span></span>
                    </div>
                  </div>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {/* Input Area */}
      <div className={styles.inputArea}>
        <form onSubmit={handleSubmit} className={styles.inputForm}>
          <div className={styles.inputWrapper}>
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={apiKey ? 'Ask a question...' : 'Configure API key to begin...'}
              disabled={isLoading || !apiKey}
              className={styles.inputField}
              aria-label="Message input"
            />
            <button
              type="submit"
              disabled={isLoading || !input.trim() || !apiKey}
              className={styles.sendButton}
              title="Send message"
              aria-label="Send message"
            >
              {isLoading ? (
                <span className={styles.loadingSpinner}></span>
              ) : (
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z"/>
                </svg>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}