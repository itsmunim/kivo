import { useCallback, useEffect, useRef, useState } from 'react'

async function api(path, opts) {
  const res = await fetch(path, opts)
  if (!res.ok) throw new Error(`${path} -> HTTP ${res.status}`)
  return res.json()
}

function fmtTTL(ms) {
  if (ms <= 0) return ''
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${Math.floor(ms / 60000)}m${Math.round((ms % 60000) / 1000)}s`
}

function FormatArray({ items }) {
  if (!items || items.length === 0) return <span className="r-nil">(empty array)</span>
  return (
    <div className="array">
      {items.map((item, i) => (
        <div className="arr-row" key={i}>
          <span className="arr-idx">{i + 1}</span>
          <FormatValue result={item} />
        </div>
      ))}
    </div>
  )
}

function FormatValue({ result }) {
  const kind = result?.kind
  switch (kind) {
    case 'simple':
      return <span className="r-simple">+{result.value}</span>
    case 'error':
      return <span className="r-error">-{result.value}</span>
    case 'integer':
      return <span className="r-int">:{result.value}</span>
    case 'bulk':
      return result.value === null ? (
        <span className="r-nil">(nil)</span>
      ) : (
        <span className="r-bulk">"{result.value}"</span>
      )
    case 'array':
      return result.value === null ? (
        <span className="r-nil">(nil)</span>
      ) : (
        <FormatArray items={result.value} />
      )
    default:
      return <span className="r-nil">{JSON.stringify(result)}</span>
  }
}

function Entry({ entry }) {
  return (
    <div className="entry">
      <div className="cmd-line">
        <span className="prompt">$</span>
        <span className="cmd-text">{entry.command}</span>
        {entry.duration_ms >= 0 && (
          <span className="dur">{entry.duration_ms.toFixed(2)}ms</span>
        )}
      </div>
      <div className="result">
        <FormatValue result={entry.result} />
      </div>
    </div>
  )
}

export default function App() {
  const [keys, setKeys] = useState([])
  const [dbSize, setDbSize] = useState(0)
  const [entries, setEntries] = useState([])
  const [input, setInput] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState(null)

  const historyRef = useRef([])
  const histIdxRef = useRef(0)
  const endRef = useRef(null)

  const refreshKeys = useCallback(async () => {
    try {
      const data = await api('/api/keys')
      setKeys(data.keys || [])
      setDbSize(data.count ?? 0)
    } catch {
      /* transient — ignore, next poll will retry */
    }
  }, [])

  useEffect(() => {
    refreshKeys()
    const id = setInterval(refreshKeys, 5000)
    return () => clearInterval(id)
  }, [refreshKeys])

  useEffect(() => {
    endRef.current?.scrollIntoView({ block: 'end' })
  }, [entries])

  const run = async (raw) => {
    const cmd = raw.trim()
    if (!cmd || busy) return
    setBusy(true)
    try {
      const data = await api('/api/command', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ command: cmd }),
      })
      setEntries((prev) => [
        ...prev,
        { command: cmd, result: data.result, duration_ms: data.duration_ms ?? -1 },
      ])
      historyRef.current.push(cmd)
      histIdxRef.current = historyRef.current.length
      setError(null)
    } catch (err) {
      setError(String(err))
    } finally {
      setBusy(false)
      refreshKeys()
      setInput('')
    }
  }

  const handleKey = (e) => {
    if (e.key === 'Enter') {
      run(input)
      return
    }
    if (e.key === 'ArrowUp') {
      e.preventDefault()
      const h = historyRef.current
      if (!h.length) return
      histIdxRef.current = Math.max(0, histIdxRef.current - 1)
      setInput(h[histIdxRef.current] ?? '')
      return
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      const h = historyRef.current
      if (!h.length) return
      histIdxRef.current = Math.min(h.length, histIdxRef.current + 1)
      setInput(histIdxRef.current >= h.length ? '' : h[histIdxRef.current])
    }
  }

  return (
    <div className="app">
      <header className="topbar">
        <div className="brand">
          <span className="word">kivo</span>
          <span className="cursor" />
          <span className="sub">web console</span>
        </div>
        <div className="stats">
          <span className="stat">{dbSize} keys</span>
          <button className="btn-refresh" onClick={refreshKeys} disabled={busy}>
            Refresh
          </button>
          <a
            className="github-link"
            href="https://github.com/itsmunim/kivo"
            target="_blank"
            rel="noreferrer"
          >
            GitHub
          </a>
        </div>
      </header>

      <main className="main">
        <aside className="keys-panel">
          <div className="panel-head">
            <h2>Keys</h2>
            <span className="k-count">{keys.length}</span>
          </div>
          <ul className="key-list">
            {keys.map((k) => (
              <li
                key={k.key}
                className="key-row"
                onClick={() => setInput(`GET ${k.key}`)}
                title="Click to GET this key"
              >
                <span className={`type-badge t-${k.type}`}>{k.type}</span>
                <span className="key-name">{k.key}</span>
                {k.ttl_ms > -1 && <span className="ttl">{fmtTTL(k.ttl_ms)}</span>}
              </li>
            ))}
          </ul>
          {keys.length === 0 && (
            <p className="empty">
              No keys yet. Try <code>SET foo bar</code> in the terminal.
            </p>
          )}
        </aside>

        <section className="terminal-panel">
          <div className="terminal">
            <div className="entries">
              <div className="entry welcome">
                <span className="prompt">$</span>
                <span className="cmd-text">kivo web console — type Redis commands, press Enter. ↵↑/↓ recalls history. Click a key to GET it.</span>
              </div>
              {entries.map((en, i) => (
                <Entry key={i} entry={en} />
              ))}
              <div ref={endRef} />
            </div>
            <div className="input-row">
              <span className="prompt">$</span>
              <input
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKey}
                placeholder="Type a command, e.g. SET foo bar"
                autoFocus
                spellCheck={false}
                autoComplete="off"
              />
            </div>
          </div>
          {error && <div className="errbar">✕ {error}</div>}
        </section>
      </main>
    </div>
  )
}