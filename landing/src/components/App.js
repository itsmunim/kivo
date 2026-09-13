export function renderApp() {
  return `
    <div class="page">
      <nav class="nav">
        <div class="nav-inner">
          <a href="#" class="logo" aria-label="kivo">
            <span class="resp-word">kivo</span><span class="resp-cursor"></span>
          </a>
          <div class="nav-links">
            <a href="#features">Features</a>
            <a href="#docker">Docker</a>
            <a href="#kubernetes">Kubernetes</a>
            <a href="https://github.com/itsmunim/kivo" target="_blank">GitHub</a>
          </div>
        </div>
      </nav>

      <header class="hero">
        <div class="hero-glow"></div>
        <div class="hero-content">
          <div class="logo-resp" aria-label="kivo">
            <span class="resp-prompt">$</span>
            <span class="resp-word">kivo</span><span class="resp-cursor"></span>
          </div>
          <h1>A fast, lightweight, Redis-compatible in-memory data store</h1>
          <p class="tagline">Written in Go. Simple. Hackable. Reliable.</p>
          <div class="hero-cta">
            <a href="https://github.com/itsmunim/kivo" class="btn btn-primary" target="_blank">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
              View on GitHub
            </a>
            <a href="#docker" class="btn btn-secondary">Get Started</a>
          </div>
        </div>
      </header>

      <section class="features" id="features">
        <div class="container">
          <h2>What is kivo?</h2>
          <div class="feature-grid">
            <div class="feature-card">
              <div class="feature-icon">⚡</div>
              <h3>Redis Protocol Compatible</h3>
              <p>Drop-in replacement for Redis. Works with redis-cli, go-redis, ioredis, and every Redis client library. RESP2 protocol with pipelining support.</p>
            </div>
            <div class="feature-card">
              <div class="feature-icon">🔧</div>
              <h3>Simple & Hackable</h3>
              <p>Small, clean Go codebase. Easy to understand, easy to extend. Single-threaded async design with one goroutine per connection.</p>
            </div>
            <div class="feature-card">
              <div class="feature-icon">🛡️</div>
              <h3>Reliable</h3>
              <p>AOF persistence with configurable sync strategies. Graceful shutdown. Lazy + active expiration. Memory limits with safety checks.</p>
            </div>
            <div class="feature-card">
              <div class="feature-icon">📈</div>
              <h3>Built to Grow</h3>
              <p>Single-node and reliable today. Foundation for distributed sharding, clustering, and replication tomorrow.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="code-section" id="docker">
        <div class="container">
          <h2>Run with Docker</h2>
          <p>Pull the public image from GitHub Container Registry:</p>
          <div class="code-block">
            <pre><code>docker pull ghcr.io/itsmunim/kivo:latest
docker run -p 6379:6379 ghcr.io/itsmunim/kivo:latest</code></pre>
            <button class="copy-btn" onclick="copyCode(this)">Copy</button>
          </div>
          <p>Connect with any Redis client:</p>
          <div class="code-block">
            <pre><code>redis-cli -p 6379 PING
# PONG</code></pre>
            <button class="copy-btn" onclick="copyCode(this)">Copy</button>
          </div>
        </div>
      </section>

      <section class="code-section kubernetes" id="kubernetes">
        <div class="container">
          <h2>Deploy on Kubernetes</h2>
          <p>A 1-replica StatefulSet with an attached PVC — the AOF survives pod restarts. One command:</p>
          <div class="code-block">
            <pre><code>kubectl apply -f https://raw.githubusercontent.com/itsmunim/kivo/main/k8s/kivo-deploy.yaml</code></pre>
            <button class="copy-btn" onclick="copyCode(this)">Copy</button>
          </div>
          <div class="note">
            <strong>Note:</strong> the default manifest uses your cluster's default StorageClass, which on many clusters is node-local (hostPath). That is fine for testing, but data is lost if the node is deleted or crashes. For durable persistence, download <code>k8s/kivo-deploy.yaml</code>, set <code>storageClassName</code> to detachable network storage (AWS EBS <code>gp2</code>, GCE PD <code>standard</code>, Azure Disk <code>managed-csi</code>, Longhorn...), then apply your copy.
          </div>
        </div>
      </section>

      <section class="features">
        <div class="container">
          <h2>Supported Commands</h2>
          <div class="command-grid">
            <div class="command-group">
              <h4>Strings</h4>
              <code>GET</code> <code>SET</code> <code>MGET</code> <code>MSET</code> <code>INCR</code> <code>DECR</code> <code>APPEND</code> <code>STRLEN</code>
            </div>
            <div class="command-group">
              <h4>Lists</h4>
              <code>LPUSH</code> <code>RPUSH</code> <code>LPOP</code> <code>RPOP</code> <code>LRANGE</code> <code>LLEN</code> <code>LINDEX</code> <code>LREM</code> <code>LTRIM</code>
            </div>
            <div class="command-group">
              <h4>Sets</h4>
              <code>SADD</code> <code>SREM</code> <code>SMEMBERS</code> <code>SISMEMBER</code> <code>SCARD</code> <code>SPOP</code> <code>SUNION</code> <code>SINTER</code> <code>SDIFF</code>
            </div>
            <div class="command-group">
              <h4>Hashes</h4>
              <code>HSET</code> <code>HGET</code> <code>HGETALL</code> <code>HDEL</code> <code>HLEN</code> <code>HEXISTS</code> <code>HKEYS</code> <code>HVALS</code> <code>HMGET</code>
            </div>
            <div class="command-group">
              <h4>Sorted Sets</h4>
              <code>ZADD</code> <code>ZREM</code> <code>ZRANGE</code> <code>ZREVRANGE</code> <code>ZRANGEBYSCORE</code> <code>ZCARD</code> <code>ZSCORE</code> <code>ZINCRBY</code> <code>ZCOUNT</code>
            </div>
            <div class="command-group">
              <h4>Keys & Server</h4>
              <code>KEYS</code> <code>FLUSHDB</code> <code>DBSIZE</code> <code>TYPE</code> <code>EXPIRE</code> <code>TTL</code> <code>PING</code> <code>ECHO</code>
            </div>
          </div>
        </div>
      </section>

      <footer class="footer">
        <div class="container">
          <div class="footer-logo">
            <span class="resp-word">kivo</span><span class="resp-cursor"></span>
          </div>
          <p>MIT Licensed · Built with Go</p>
          <div class="footer-links">
            <a href="https://github.com/itsmunim/kivo" target="_blank">GitHub</a>
            <a href="https://github.com/itsmunim/kivo/issues" target="_blank">Issues</a>
            <a href="https://github.com/itsmunim/kivo/blob/main/plans/v1-design.md" target="_blank">Design Doc</a>
          </div>
        </div>
      </footer>
    </div>

    <script>
      function copyCode(btn) {
        const code = btn.previousElementSibling.querySelector('code').innerText;
        navigator.clipboard.writeText(code).then(() => {
          btn.textContent = 'Copied!';
          setTimeout(() => btn.textContent = 'Copy', 2000);
        });
      }
    </script>
  `;
}