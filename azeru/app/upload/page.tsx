'use client'

import { useState, useCallback } from 'react'
import Link from 'next/link'
import { uploadPDF } from '@/lib/api'
import styles from '../css/upload.module.css'

export default function UploadPage() {
  const [isDragging, setIsDragging] = useState(false)
  const [isUploading, setIsUploading] = useState(false)
  const [uploadStatus, setUploadStatus] = useState<{
    type: 'success' | 'error' | null
    message: string
  }>({ type: null, message: '' })

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }, [])

  const handleDrop = useCallback(async (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)

    const files = Array.from(e.dataTransfer.files)
    const pdfFile = files.find(f => f.type === 'application/pdf')

    if (!pdfFile) {
      setUploadStatus({ type: 'error', message: 'Please drop a PDF file only' })
      return
    }

    await uploadFile(pdfFile)
  }, [])

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    if (file.type !== 'application/pdf') {
      setUploadStatus({ type: 'error', message: 'Please select a PDF file only' })
      return
    }

    await uploadFile(file)
  }

  const uploadFile = async (file: File) => {
    setIsUploading(true)
    setUploadStatus({ type: null, message: '' })

    try {
      const result = await uploadPDF(file)
      setUploadStatus({
        type: 'success',
        message: `"${result.document.file_name}" uploaded successfully. Processing ${result.document.file_size} bytes...`,
      })
    } catch (error) {
      setUploadStatus({
        type: 'error',
        message: `Upload failed: ${error instanceof Error ? error.message : 'Unknown error'}`,
      })
    } finally {
      setIsUploading(false)
    }
  }

  return (
    <div className={styles.container}>

      {/* Header */}
      <header className={styles.header}>
        <Link href="/" className={styles.backLink}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
          Dashboard
        </Link>
        <span className={styles.divider}></span>
        <h1>Upload Documents</h1>
        <span className={styles.divider}></span>
        <p>Add PDFs to your knowledge base</p>
      </header>

      {/* Body */}
      <div className={styles.body}>

        {/* Drop zone */}
        <div
          className={`${styles.dropZone} ${isDragging ? styles.dragging : ''} ${isUploading ? styles.uploading : ''}`}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <input
            type="file"
            accept=".pdf,application/pdf"
            onChange={handleFileSelect}
            className={styles.fileInput}
            id="file-input"
            disabled={isUploading}
          />

          <div className={styles.dropContent}>
            {isUploading ? (
              <>
                <div className={styles.spinner}></div>
                <p className={styles.uploadingText}>Uploading and processing document...</p>
              </>
            ) : (
              <>
                <div className={styles.uploadIconWrap}>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="17 8 12 3 7 8"/>
                    <line x1="12" y1="3" x2="12" y2="15"/>
                  </svg>
                </div>
                <p className={styles.dropText}>
                  <strong>Drag and drop</strong> your PDF here, or{' '}
                  <label htmlFor="file-input" className={styles.browseLink}>
                    browse files
                  </label>
                </p>
                <p className={styles.hint}>PDF files only — up to 50 MB</p>
              </>
            )}
          </div>
        </div>

        {/* Status */}
        {uploadStatus.type && (
          <div className={`${styles.status} ${styles[uploadStatus.type]}`}>
            {uploadStatus.type === 'success' ? (
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
            ) : (
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
            )}
            {uploadStatus.message}
          </div>
        )}

        {/* Guidelines */}
        <div className={styles.info}>
          <div className={styles.infoHeader}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <h3>Upload Guidelines</h3>
          </div>
          <ul>
            <li>Text-based PDFs work best — scanned images require OCR pre-processing</li>
            <li>Maximum file size is 50 MB per document</li>
            <li>Documents are processed and indexed in the background</li>
            <li>All data remains on your infrastructure — NDPA 2023 compliant</li>
          </ul>
        </div>
      </div>
    </div>
  )
}