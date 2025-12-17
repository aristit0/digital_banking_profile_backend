# Customer 360 API - Golang Backend

RESTful API untuk mengakses Customer 360 data dari Couchbase.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Couchbase Server running di localhost
- Bucket `customer_360` dengan data

### Installation

```bash
cd backend

# Install dependencies
go mod download

# Run server
go run main.go
```

Server akan running di `http://localhost:8080`

## 📡 API Endpoints

### 1. Health Check
```bash
GET http://localhost:8080/api/v1/health
```

Response:
```json
{
  "status": "healthy",
  "service": "customer-360-api",
  "time": "2024-12-17T10:00:00Z"
}
```

### 2. Get Statistics
```bash
GET http://localhost:8080/api/v1/stats
```

Response:
```json
{
  "service": "customer-360-api",
  "time": "2024-12-17T10:00:00Z",
  "total_customers": 10000
}
```

### 3. Get Customer 360 Data (Main Endpoint)
```bash
POST http://localhost:8080/api/v1/customers
Content-Type: application/json

{
  "customer_id": ["CUST0000001", "CUST0000002", "CUST0000003"]
}
```

Response:
```json
{
  "success": true,
  "count": 3,
  "data": [
    {
      "customer_id": "CUST0000001",
      "customer": {
        "customer_id": "CUST0000001",
        "personal_info": {
          "full_name": "Budi Santoso",
          "age": 39,
          "gender": "M",
          ...
        },
        "demographics": {
          "occupation": "professional",
          "income_range": "10M-25M",
          ...
        },
        "status": {
          "segment": "mass_affluent",
          "customer_status": "active",
          ...
        }
      },
      "address": {
        "city": "Jakarta",
        "province": "DKI Jakarta",
        ...
      },
      "contact": {
        "primary_phone": "081234567890",
        "email_primary": "budi@email.com",
        ...
      },
      "accounts": [
        {
          "account_number": "1234567890",
          "account_type": "savings",
          "balance": 50000000,
          ...
        }
      ],
      "deposits": [...],
      "loans": [...],
      "cards": [...],
      "investments": [...],
      "segment": {
        "current_segment": "mass_affluent",
        "lifetime_value": 75000000,
        ...
      },
      "behavior": {
        "product_ownership": {
          "total_products": 5,
          ...
        },
        "channel_usage": {
          "mobile_banking_active": true,
          ...
        }
      },
      "preference": {
        "channel_preferences": {
          "preferred_channel": "mobile_banking",
          ...
        }
      }
    },
    ...
  ]
}
```

### Error Response
```json
{
  "success": false,
  "count": 1,
  "data": [...],
  "errors": [
    {
      "customer_id": "CUST9999999",
      "error": "customer not found"
    }
  ]
}
```

## 🧪 Testing dengan cURL

### Single Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{"customer_id": ["CUST0000001"]}'
```

### Multiple Customers
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{"customer_id": ["CUST0000001","CUST0000002","CUST0000003"]}'
```

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

## 📦 Project Structure

```
backend/
├── main.go           # Main application
├── go.mod            # Go dependencies
├── go.sum            # Dependency checksums
└── README.md         # This file
```

## 🔧 Configuration

Edit constants di `main.go`:

```go
const (
    CouchbaseHost     = "localhost"
    CouchbaseUser     = "admin"
    CouchbasePassword = "T1ku$H1t4m"
    BucketName        = "customer_360"
)
```

## 🚀 Build for Production

```bash
# Build binary
go build -o customer360-api main.go

# Run binary
./customer360-api
```

## 🐳 Docker (Optional)

Create `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o customer360-api main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/customer360-api .
EXPOSE 8080
CMD ["./customer360-api"]
```

Build and run:
```bash
docker build -t customer360-api .
docker run -p 8080:8080 customer360-api
```

## 📝 API Response Structure

Setiap customer memiliki data lengkap:

| Field | Type | Description |
|-------|------|-------------|
| customer_id | string | ID customer |
| customer | object | Profile demografis |
| address | object | Alamat residential |
| contact | object | Kontak & preferences |
| accounts | array | List rekening |
| deposits | array | List deposito |
| loans | array | List pinjaman |
| cards | array | List kartu |
| investments | array | List investasi |
| segment | object | Segmentasi data |
| behavior | object | Perilaku & channel usage |
| preference | object | Preferensi customer |

## 🔒 CORS Configuration

API sudah dikonfigurasi untuk accept request dari:
- `http://localhost:5173` (Vite default)
- `http://localhost:3000` (React default)

Untuk production, update di `setupRouter()`:
```go
config.AllowOrigins = []string{"https://your-domain.com"}
```

## 📈 Performance Tips

1. **Connection Pooling**: Gocb SDK sudah handle connection pooling
2. **Parallel Queries**: Bisa optimize dengan goroutines untuk multiple customers
3. **Caching**: Implement Redis untuk frequently accessed data
4. **Rate Limiting**: Add middleware untuk limit requests

## 🐛 Troubleshooting

### Connection Failed
```
Failed to initialize Couchbase: failed to connect to cluster
```
**Solution**: 
- Pastikan Couchbase Server running
- Check connection string dan credentials
- Verify bucket exists

### Customer Not Found
```
customer not found: document not found
```
**Solution**: 
- Verify customer_id exists di database
- Check document key format: `customer::CUST0000001`

### Port Already in Use
```
bind: address already in use
```
**Solution**: 
- Change port di `main()`: `router.Run(":8081")`
- Or kill process using port 8080

## 📚 Dependencies

- **gin-gonic/gin**: Web framework
- **couchbase/gocb**: Couchbase SDK
- **gin-contrib/cors**: CORS middleware

## 🎯 Next Steps

- [ ] Add authentication/authorization
- [ ] Implement rate limiting
- [ ] Add request logging
- [ ] Add metrics/monitoring
- [ ] Implement caching layer
- [ ] Add pagination support
- [ ] Add filtering/sorting options
