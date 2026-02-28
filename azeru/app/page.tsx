import Link from 'next/link'
import styles from './page.module.css'

export default function Home() {
  return (
    <div className={styles.dashboard}>
      {/* Animated Background Elements */}
      <div className={styles.bgGradient}></div>
      <div className={styles.bgAccent}></div>

      {/* Top Navigation Bar */}
      <header className={styles.navbar}>
        <div className={styles.navContent}>
          <Link href="/" className={styles.logoBlock}>
            <div className={styles.logoIcon}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M12 2a9 9 0 0 1 9 9c0 3.6-2.1 6.7-5.2 8.2L14 22H10l-1.8-2.8C5.1 17.7 3 14.6 3 11a9 9 0 0 1 9-9Z"/>
                <circle cx="12" cy="11" r="3"/>
              </svg>
            </div>
            <div>
              <div className={styles.logoTitle}>Enterprise Brain</div>
              <div className={styles.logoSubtitle}>Knowledge Platform</div>
            </div>
          </Link>

          <div className={styles.navBadges}>
            <span className={styles.badge}>
              <span className={styles.badgeDot}></span>
              Self-hosted
            </span>
            <span className={styles.badge}>NDPA 2023</span>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className={styles.container}>
        
        {/* Left Column - Hero Content */}
        <section className={styles.heroSection}>
          <div className={styles.heroContent}>
            <span className={styles.eyebrow}>Knowledge Management</span>
            <h1 className={styles.heading}>
              Enterprise documents at your fingertips.
            </h1>
            <p className={styles.subtitle}>
              Upload, index, and query your internal knowledge base with semantic AI search. Powered by your infrastructure, zero data leaves your servers.
            </p>

            <div className={styles.heroStats}>
              <div className={styles.stat}>
                <span className={styles.statLabel}>Security</span>
                <span className={styles.statValue}>100%</span>
              </div>
              <div className={styles.statDivider}></div>
              <div className={styles.stat}>
                <span className={styles.statLabel}>Private</span>
                <span className={styles.statValue}>On-Prem</span>
              </div>
              <div className={styles.statDivider}></div>
              <div className={styles.stat}>
                <span className={styles.statLabel}>Data</span>
                <span className={styles.statValue}>BYOK</span>
              </div>
            </div>
          </div>
        </section>

        {/* Right Column - Cards Grid */}
        <section className={styles.cardsSection}>
          <div className={styles.cardsGrid}>
            {/* Card 1 - Upload */}
            <Link href="/upload" className={`${styles.card} ${styles.card1}`}>
              <div className={styles.cardHeader}>
                <div className={styles.cardIcon}>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                    <line x1="12" y1="18" x2="12" y2="12"/>
                    <line x1="9" y1="15" x2="15" y2="15"/>
                  </svg>
                </div>
                <span className={styles.cardBadge}>Primary</span>
              </div>
              <h3 className={styles.cardTitle}>Upload Documents</h3>
              <p className={styles.cardDescription}>Import PDFs and build your knowledge base</p>
              <span className={styles.cardArrow}>→</span>
            </Link>

            {/* Card 2 - Chat */}
            <Link href="/chat" className={`${styles.card} ${styles.card2}`}>
              <div className={styles.cardHeader}>
                <div className={styles.cardIcon}>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                  </svg>
                </div>
                <span className={styles.cardBadge}>Query</span>
              </div>
              <h3 className={styles.cardTitle}>Query Knowledge</h3>
              <p className={styles.cardDescription}>Ask questions and get cited answers</p>
              <span className={styles.cardArrow}>→</span>
            </Link>

            {/* Card 3 - Library */}
            <Link href="/documents" className={`${styles.card} ${styles.card3}`}>
              <div className={styles.cardHeader}>
                <div className={styles.cardIcon}>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                  </svg>
                </div>
                <span className={styles.cardBadge}>Library</span>
              </div>
              <h3 className={styles.cardTitle}>Document Library</h3>
              <p className={styles.cardDescription}>Browse and manage all indexed documents</p>
              <span className={styles.cardArrow}>→</span>
            </Link>

            {/* Card 4 - Settings */}
            <Link href="/settings" className={`${styles.card} ${styles.card4}`}>
              <div className={styles.cardHeader}>
                <div className={styles.cardIcon}>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
                  </svg>
                </div>
                <span className={styles.cardBadge}>Config</span>
              </div>
              <h3 className={styles.cardTitle}>Configuration</h3>
              <p className={styles.cardDescription}>Manage settings and API credentials</p>
              <span className={styles.cardArrow}>→</span>
            </Link>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className={styles.footer}>
        <div className={styles.footerContent}>
          <div className={styles.footerLeft}>
            <span className={styles.footerBadge}>Self-hosted</span>
            <span className={styles.footerBadge}>NDPA 2023</span>
            <span className={styles.footerBadge}>BYOK</span>
          </div>
          <div className={styles.footerRight}>
            <span>© 2024 Enterprise Brain</span>
          </div>
        </div>
      </footer>
    </div>
  )
}