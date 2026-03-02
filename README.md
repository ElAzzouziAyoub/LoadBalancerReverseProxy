# Load Balancer Reverse Proxy

A high-performance Go-based reverse proxy server with built-in load balancing, rate limiting, and dynamic backend management.

## 🚀 Overview

This project implements a **reverse proxy server** that distributes incoming HTTP requests across multiple backend servers using **round-robin load balancing**. It includes:

- **Round-Robin Load Balancing**: Evenly distributes requests across healthy backends
- **Rate Limiting**: Limits requests to 10 per minute per IP address
- **Health Tracking**: Monitors active connections and backend status
- **Dynamic Backend Management**: Add/remove backends without restarting via Admin API
- **Swagger/OpenAPI Documentation**: Interactive API documentation at `/swagger/`

## 📋 Architecture

```
Client Requests
     ↓
  Port 9090 (Reverse Proxy)
     ↓ (Round-Robin Distribution)
  ┌──────┬──────┬──────┐
  ↓      ↓      ↓      ↓
Port 8081  Port 8082  Port 8083
Backend 1  Backend 2  Backend 3
```

## 🛠️ Prerequisites

- Go 1.25.6 or higher
- Linux/Mac/Windows with standard utilities

## 📦 Installation

1. **Clone or navigate to the project directory**:
```bash
cd LoadBalancerReverseProxy
```

2. **Download dependencies**:
```bash
go mod download
```

## 🚀 Running the Application

### Step 1: Start the Backend Servers

Open three separate terminals and run each backend server:

**Terminal 1 - Backend Server 1 (Port 8081)**:
```bash
cd BackEndServers
go run backend1.go
```

**Terminal 2 - Backend Server 2 (Port 8082)**:
```bash
cd BackEndServers
go run backend2.go
```

**Terminal 3 - Backend Server 3 (Port 8083)**:
```bash
cd BackEndServers
go run backend3.go
```

### Step 2: Start the Reverse Proxy

**Terminal 4 - Main Reverse Proxy**:
```bash
go run ReverseProxy.go
```

You should see output like:
```
Started Server 1 : 
Started Server 2 : 
Started Server 3 : 
Admin API listening on :9091
Proxy listening on :9090
```

## 📡 API Usage

### 1. Send Requests Through the Proxy

The proxy listens on **port 9090**. All requests are balanced across the three backends:

```bash
# Simple GET request
curl http://localhost:9090/

# Output will show which backend handled the request
# Forwarded to http://localhost:8081
# Connection established with server
```

**Test Round-Robin Load Balancing**:
```bash
# Run this multiple times to see requests distributed across backends
for i in {1..6}; do
  echo "Request $i:"
  curl http://localhost:9090/
  echo ""
done
```

### 2. Admin API - Manage Backends

The admin API listens on **port 9091** and provides endpoints to manage backends dynamically.

**List all backends**:
```bash
curl http://localhost:9091/admin/backends
```

Response:
```json
[
  {
    "URL": "http://localhost:8081",
    "Alive": true,
    "Connections": 0
  },
  {
    "URL": "http://localhost:8082",
    "Alive": true,
    "Connections": 0
  },
  {
    "URL": "http://localhost:8083",
    "Alive": true,
    "Connections": 0
  }
]
```

**Add a new backend**:
```bash
curl -X POST "http://localhost:9091/admin/backends/add?url=http://localhost:8084"
```

Response: `backend added`

**Remove a backend**:
```bash
curl -X DELETE "http://localhost:9091/admin/backends/remove?url=http://localhost:8081"
```

Response: `backend removed`

### 3. Swagger API Documentation

View interactive API documentation at:
```
http://localhost:9091/swagger/
```

This provides a visual interface to explore all available endpoints.

## ⚙️ Configuration

### Rate Limiting

The rate limiter is configured in [ReverseProxy.go](ReverseProxy.go#L30):
```go
limitePerMinute = 10  // Requests per minute per IP
```

To modify, edit the `limitePerMinute` variable and recompile.

### Backend Servers

Default backends are configured in [ReverseProxy.go](ReverseProxy.go#L25-L28):
```go
backends = []*Backend{
    {URL: mustParse("http://localhost:8081"), Alive: true},
    {URL: mustParse("http://localhost:8082"), Alive: true},
    {URL: mustParse("http://localhost:8083"), Alive: true},
}
```

Modify as needed before running.

## 📊 Features Explained

### Round-Robin Load Balancing
Requests are distributed sequentially across healthy backends. If you send 6 requests, they'll be distributed as: 8081 → 8082 → 8083 → 8081 → 8082 → 8083

### Rate Limiting
Each client IP is limited to 10 requests per minute. Exceeding this limit returns a `429 Too Many Requests` response. The rate limiter resets every minute.

### Connection Tracking
The proxy tracks active connections per backend to monitor load distribution.

### Health Management
Backends marked as `Alive: true` receive requests. A backend can be removed dynamically via the Admin API.

## 🧪 Testing Examples

**Test rate limiting** (should block after 10 requests):
```bash
for i in {1..15}; do
  echo "Request $i:"
  curl http://localhost:9090/
  sleep 0.1
done
```

**Test with POST requests**:
```bash
curl -X POST http://localhost:9090/ -d "test data"
```

**Monitor real-time**:
```bash
# Watch backend server logs in their respective terminals
# You'll see "Connection established with [IP_ADDRESS]" for each request
```

## 📁 Project Structure

```
LoadBalancerReverseProxy/
├── ReverseProxy.go          # Main reverse proxy implementation
├── go.mod                   # Go module definition
├── openapi.yaml             # OpenAPI/Swagger specification
├── README.md                # This file
└── BackEndServers/
    ├── backend1.go          # Backend server 1 (Port 8081)
    ├── backend2.go          # Backend server 2 (Port 8082)
    └── backend3.go          # Backend server 3 (Port 8083)
```

## 🔧 Troubleshooting

**Port already in use**:
```bash
# Kill existing process on port 9090
lsof -ti:9090 | xargs kill -9
# Or modify backend ports in the source code
```

**Backends not responding**:
- Ensure all 3 backend servers are running
- Check backend logs for errors
- Verify ports 8081, 8082, 8083 are available

**Rate limiter not resetting**:
- Check server logs for "Rate limiter reset" message every minute
- This is normal behavior

## 📝 Development

To extend this proxy:

1. **Add new endpoints**: Modify handler functions in [ReverseProxy.go](ReverseProxy.go)
2. **Change load balancing algorithm**: Modify `getNextBackend()` function
3. **Add health checks**: Implement endpoint monitoring
4. **Add authentication**: Wrap handler functions with auth middleware

## 📄 License

This project is provided as-is for educational and development purposes.

## 💡 Notes

- The proxy currently forwards requests without modifying headers
- Connections are tracked but not actively monitored for backend health
- Rate limiting resets occur in-memory and are not persisted