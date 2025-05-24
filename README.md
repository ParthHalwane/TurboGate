
# TurboGate 🚀 - High Performance Reverse Proxy with Rate Limiting

![TurboGate Banner](https://github.com/user-attachments/assets/eba8710b-a376-42e8-a945-dc223630530d)


> A blazing-fast, production-grade **Reverse Proxy** built with **Golang**, featuring **dynamic route configuration**, **rate limiting**, **hot-reloadable YAML-based routing**, and **performance benchmarking up to 2 million requests per minute**.

---

## 🔗 What is a Reverse Proxy?

A **Reverse Proxy** is a server that sits between the client and one or more backend servers, forwarding client requests to those backend services and returning their responses. It is used for:

- Load balancing
- Caching
- SSL termination
- Web acceleration
- Rate limiting and security

Unlike a **forward proxy** (which sits in front of the client), a **reverse proxy** operates on behalf of the server.

---

## 🚀 About TurboGate

**TurboGate** is designed to be a plug-and-play, scalable reverse proxy solution that supports:

- Dynamic route additions at runtime
- YAML-based routing configuration
- IP-based token bucket rate limiting
- Prometheus-compatible metrics
- High concurrency with Go's net/http and goroutines

This project was created to demonstrate not just how reverse proxies work under the hood, but also how to build a **production-ready reverse proxy from scratch**, while emphasizing performance and extensibility.

---

## 📊 Key Features

- ✅ **Dynamic Routing**: Add/remove proxy routes in real-time via hot-reloadable YAML.
- ✅ **Token Bucket Rate Limiting**: Per-IP request throttling to prevent abuse.
- ✅ **Prometheus Metrics**: Integrated monitoring with `/metrics` endpoint.
- ✅ **High Throughput**: Benchmarked at up to **2 million requests/min**.
- ✅ **Fully Dockerized**: Easy deployment and local testing.
- ✅ **Simple Frontend**: Enter a domain and get a working proxy route.
- ✅ **API for Route Management**: Seamlessly integrate dynamic route control.

---

## 📹 Project Demo & Walkthrough

### 🔍 Project Walkthrough
<video width="640" height="360" src="https://github.com/user-attachments/assets/55df5440-684e-4013-9390-611a91337bcb">
</video>

---

### 🔢 JMeter Load Test
<video width="640" height="360" src="https://github.com/user-attachments/assets/ac352f88-68e5-41bf-b943-27714279dc2b">
</video>

---

## 🖼️ Performance Snapshot

![2 Million RPM](https://github.com/user-attachments/assets/e625e294-a5b9-4255-8d10-56523265d80e)  
> Achieved almost 2,000,000 requests per minute during JMeter benchmarking on a multi-core instance

---

## 📚 Tech Stack

- **Language**: Go (Golang)
- **Concurrency**: Native goroutines & channels
- **Routing Configuration**: YAML
- **Rate Limiting**: Token bucket algorithm
- **Metrics**: Prometheus-compatible
- **Deployment**: Docker & Render
- **Frontend**: React + Tailwind (Vercel hosted)

---

## ⚠️ Why Some Websites Can’t Be Proxied?

Due to modern web security protocols like **CORS (Cross-Origin Resource Sharing)** and **X-Frame-Options**, not all websites allow themselves to be reverse proxied, especially through browsers. This is a deliberate protection mechanism against:

- Clickjacking
- Cross-site forgery
- Content manipulation

So, even if TurboGate fetches the content server-side, some websites may still block rendering in the client browser.

---

## 🧪 Getting Started

Follow these steps to get the project running locally:

---

### 🔁 Backend Setup (Go)

```bash
# 1. Clone the repo
git clone https://github.com/yourusername/turbogate.git
cd turbogate

# 2. Build and run
go run cmd/main.go
```

Alternatively, using docker:
```bash
# Docker Build & Run
docker build -t turbogate .
docker run -p 8080:8080 turbogate
```

Frontend setup:
```bash
# 1. Navigate to frontend directory
cd frontend

# 2. Install dependencies
npm install

# 3. Run development server
npm run dev
```

🧪 API Usage
➕ Add New Route (POST)
```bash
POST /api/add-route
Content-Type: application/json

{
  "path": "/youtube",
  "target": "https://youtube.com"
}
```

📄 YAML Configuration Format
```bash
routes:
  - path: /github
    target: https://github.com
  - path: /openai
    target: https://openai.com
```
---

## 🌍 Some Publicly Proxyable Sites
These URLs typically allow reverse proxying:

https://github.com

https://openai.com

https://example.com

https://golang.org

https://jsonplaceholder.typicode.com

https://pokeapi.co

⚠️ Websites with strict CORS/X-Frame-Options may fail to load in browser.

---

## 💡 Future Enhancements
✅ TLS termination with custom domain support

✅ Basic Auth or API keys for route creation

✅ Caching layer (Redis or in-memory)

✅ WebSocket support

✅ Admin dashboard with metrics visualizations

✅ PostgreSQL or Redis-based persistent route store

✅ Vercel button for one-click deploy

---

## 📝 License
This project is licensed under the [MIT License](https://github.com/ParthHalwane/TurboGate/blob/main/LICENSE).
Feel free to use, modify, or extend it with proper attribution.

---

## 🙏 Acknowledgments
Inspired by tools like NGINX, Traefik, and Caddy

Special thanks to the Go community

Load testing powered by Apache JMeter

