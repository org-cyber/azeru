'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { getDocuments } from '@/lib/api'
import styles from '../css/document.module.css'

interface Document {
  id: number
  file_name: string
  file_size: number
  status: string
  chunk_count: number
  created_at: string
}

export default function DocumentsPage() {
  const [documents, setDocuments] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadDocuments()
  }, [])

  const loadDocuments = async () => {
    try {
      const docs = await getDocuments()
      setDocuments(docs)
    } catch (err) {
      setError('Failed to load documents')
    } finally {
      setLoading(false)
    }
  }

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-NG', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const isCompleted = (status: string) => status === 'completed'

  return (
    <div className={styles.container}>

      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerLeft}>
          <Link href="/" className={styles.backLink}>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="15 18 9 12 15 6"/>
            </svg>
            Dashboard
          </Link>
          <span className={styles.divider}></span>
          <h1>Document Library</h1>
          <span className={styles.divider}></span>
          <p>All indexed documents in your knowledge base</p>
        </div>
        <Link href="/upload" className={styles.uploadButton}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          Upload Document
        </Link>
      </header>

      <div className={styles.body}>
        {loading ? (
          <div className={styles.loading}>Loading documents...</div>
        ) : error ? (
          <div className={styles.error}>{error}</div>
        ) : documents.length === 0 ? (
          <div className={styles.empty}>
            <div className={styles.emptyWrap}>
              <div className={styles.emptyIcon}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                </svg>
              </div>
              <h3>No documents yet</h3>
              <p>Upload your first PDF to begin building your knowledge base.</p>
              <Link href="/upload" className={styles.uploadLink}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Upload Document
              </Link>
            </div>
          </div>
        ) : (
          <div className={styles.tableContainer}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Status</th>
                  <th>Filename</th>
                  <th>Size</th>
                  <th>Chunks</th>
                  <th>Uploaded</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {documents.map((doc) => (
                  <tr key={doc.id}>
                    <td>
                      <span
                        className={`${styles.statusBadge} ${isCompleted(doc.status) ? styles.statusOk : styles.statusPending}`}
                        title={doc.status}
                      >
                        <span className={styles.statusDot}></span>
                        {doc.status}
                      </span>
                    </td>
                    <td className={styles.filename}>{doc.file_name}</td>
                    <td>{formatSize(doc.file_size)}</td>
                    <td>
                      <span className={doc.chunk_count > 0 ? styles.chunksOk : styles.chunksEmpty}>
                        {doc.chunk_count}
                      </span>
                    </td>
                    <td className={styles.date}>{formatDate(doc.created_at)}</td>
                    <td>
                      <div className={styles.actions}>
                        <button className={styles.actionBtn} title="View Details">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
                          </svg>
                        </button>
                        <button className={styles.actionBtn} title="Delete">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            <div className={styles.summary}>
              <strong>{documents.length}</strong> documents &middot;{' '}
              <strong>{documents.reduce((sum, d) => sum + d.chunk_count, 0)}</strong> chunks &middot;{' '}
              <strong>{formatSize(documents.reduce((sum, d) => sum + d.file_size, 0))}</strong> total
            </div>
          </div>
        )}
      </div>
    </div>
  )
}