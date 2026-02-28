'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import styles from '../css/settings.module.css'

export default function SettingsPage() {
  const [apiKey, setApiKey] = useState('')
  const [showKey, setShowKey] = useState(false)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    const savedKey = localStorage.getItem('enterprise_brain_api_key')
    if (savedKey) setApiKey(savedKey)
  }, [])

  const handleSave = () => {
    localStorage.setItem('enterprise_brain_api_key', apiKey)
    setSaved(true)
    setTimeout(() => setSaved(false), 3000)
  }

  const handleClear = () => {
    if (confirm('Clear API key? You will need to re-enter it to use the chat.')) {
      localStorage.removeItem('enterprise_brain_api_key')
      setApiKey('')
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
        <h1>Configuration</h1>
        <span className={styles.divider}></span>
        <p>API keys, system info, and compliance settings</p>
      </header>

      <div className={styles.body}>
        <div className={styles.grid}>

          {/* API Key */}
          <section className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardHeaderLeft}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                <h2>API Configuration</h2>
              </div>
              <span className={styles.badge}>Required</span>
            </div>

            <div className={styles.cardBody}>
              <div className={styles.field}>
                <label htmlFor="api-key-input">Groq API Key</label>
                <p className={styles.help}>
                  Your API key is stored locally in your browser and never transmitted to our servers.{' '}
                  <a href="https://console.groq.com/keys" target="_blank" rel="noopener noreferrer">
                    Obtain key from Groq Console &rarr;
                  </a>
                </p>

                <div className={styles.inputGroup}>
                  <input
                    id="api-key-input"
                    type={showKey ? 'text' : 'password'}
                    value={apiKey}
                    onChange={(e) => setApiKey(e.target.value)}
                    placeholder="gsk_..."
                    className={styles.input}
                  />
                  <button
                    onClick={() => setShowKey(!showKey)}
                    className={styles.toggleBtn}
                    title={showKey ? 'Hide key' : 'Show key'}
                  >
                    {showKey ? (
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/>
                      </svg>
                    ) : (
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
                      </svg>
                    )}
                  </button>
                </div>

                <div className={styles.actions}>
                  <button
                    onClick={handleSave}
                    disabled={!apiKey.trim()}
                    className={styles.saveBtn}
                  >
                    {saved ? 'Saved' : 'Save Key'}
                  </button>
                  {apiKey && (
                    <button onClick={handleClear} className={styles.clearBtn}>
                      Clear
                    </button>
                  )}
                </div>
              </div>
            </div>
          </section>

          {/* System Info */}
          <section className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardHeaderLeft}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/>
                </svg>
                <h2>System Information</h2>
              </div>
            </div>

            <div className={styles.cardBody}>
              <div className={styles.infoGrid}>
                <div className={styles.infoItem}>
                  <span className={styles.label}>Version</span>
                  <span className={styles.value}>enterprise-brain v1.0.0</span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.label}>Backend</span>
                  <span className={styles.value}>Go + SQLite + FTS5</span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.label}>AI Model</span>
                  <span className={styles.value}>meta-llama/llama-4-scout-17b-16e-instruct</span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.label}>Data Storage</span>
                  <span className={styles.value}>Local SQLite (Encrypted)</span>
                </div>
              </div>
            </div>
          </section>

          {/* Compliance */}
          <section className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardHeaderLeft}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                </svg>
                <h2>Compliance &amp; Security</h2>
              </div>
            </div>

            <div className={styles.cardBody}>
              <ul className={styles.complianceList}>
                {[
                  { title: 'NDPA 2023 Compliant', desc: 'Full compliance with the Nigeria Data Protection Act' },
                  { title: 'Self-Hosted Infrastructure', desc: 'All data remains entirely on your servers' },
                  { title: 'BYOK Architecture', desc: 'Bring Your Own Key — API credentials never leave your browser' },
                  { title: 'No External Data Retention', desc: 'Documents are never transmitted to third-party servers' },
                ].map(({ title, desc }) => (
                  <li key={title}>
                    <div className={styles.checkWrap}>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="20 6 9 17 4 12"/>
                      </svg>
                    </div>
                    <div>
                      <strong>{title}</strong>
                      <p>{desc}</p>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          </section>

          {/* Support */}
          <section className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardHeaderLeft}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/>
                </svg>
                <h2>Support</h2>
              </div>
            </div>

            <div className={styles.cardBody}>
              <div className={styles.support}>
                <p>Need assistance with your Enterprise Brain deployment?</p>
                <div className={styles.supportLinks}>
                  <a href="mailto:support@enterprisebrain.ng" className={styles.link}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/>
                    </svg>
                    Email Support
                  </a>
                  <a href="#" className={styles.link}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/>
                    </svg>
                    Documentation
                  </a>
                  <a href="#" className={styles.link}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
                    </svg>
                    Video Tutorials
                  </a>
                </div>
              </div>
            </div>
          </section>

        </div>

        <footer className={styles.footer}>
          &copy; 2026 Enterprise Brain. Built for Nigerian Enterprises.
          <span className={styles.motto}>Secure. Compliant. Intelligent.</span>
        </footer>
      </div>
    </div>
  )
}