**创建短链接**

```sh
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl": "https://example.com", "customSuffix": "x", "expiration": 60}'
```

**访问短链接**

```sh
curl -v http://localhost:8080/x1
```

**查看统计**

```sh
curl http://localhost:8080/api/v1/stats/x1
```

