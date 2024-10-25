// src/App.jsx
import { useState, useEffect } from 'react'
import './App.css'
import axios from 'axios'
import { FiLink, FiCopy, FiCheck, FiLoader, FiMoon, FiSun } from 'react-icons/fi'
import config from './config'; // Import the config

function App() {
  const [longUrl, setLongUrl] = useState('')
  const [shortUrl, setShortUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [copied, setCopied] = useState(false)
  const [darkMode, setDarkMode] = useState(() => {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  })

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', darkMode ? 'dark' : 'light')
  }, [darkMode])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    
    try {
      const response = await axios.post(`${config.BACKEND_URL}/shorten?longUrl=${longUrl}`)

      setLongUrl(longUrl)
      setShortUrl(response.data)
      
    } catch (err) {
      setError('Failed to shorten URL. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(shortUrl)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="app-container">
      <button 
        className="theme-toggle"
        onClick={() => setDarkMode(!darkMode)}
        aria-label="Toggle theme"
      >
        {darkMode ? <FiSun /> : <FiMoon />}
      </button>

      <div className="glass-card">
        <div className="content-wrapper">
          <div className="header">
            <div className="logo-container">
              <FiLink className="logo-icon" />
            </div>
            <h1>URL Shortener</h1>
            <p className="subtitle">Transform your long URLs into short, shareable links</p>
          </div>

          <form onSubmit={handleSubmit} className="url-form">
            <div className="input-wrapper">
              <input
                type="url"
                value={longUrl}
                onChange={(e) => setLongUrl(e.target.value)}
                placeholder="Paste your long URL here"
                required
                className="url-input"
              />
              <button 
                type="submit" 
                className="submit-button"
                disabled={loading}
              >
                {loading ? (
                  <FiLoader className="spinner" />
                ) : (
                  'Shorten'
                )}
              </button>
            </div>
          </form>
          
          {error && (
            <div className="error-message">
              <p>{error}</p>
            </div>
          )}
          
          {shortUrl && (
            <div className="result-card">
              <div className="result-header">
                <h2>Your shortened URL is ready!</h2>
              </div>
              <div className="url-result">
                <a href={shortUrl} target="_blank" rel="noopener noreferrer">
                  {shortUrl}
                </a>
                <button 
                  onClick={copyToClipboard}
                  className={`copy-button ${copied ? 'copied' : ''}`}
                >
                  {copied ? (
                    <>
                      <FiCheck /> Copied
                    </>
                  ) : (
                    <>
                      <FiCopy /> Copy
                    </>
                  )}
                </button>
              </div>
              <div className="stats">
                <p>Original URL length: {longUrl.length} characters</p>
                <p>Shortened URL length: {shortUrl.length} characters</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default App