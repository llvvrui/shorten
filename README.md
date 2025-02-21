**创建短链接**

```sh
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl": "https://example.com"}'
```

**访问短链接**

```sh
curl -v http://localhost:8080/abc123
```

**查看统计**

```sh
curl http://localhost:8080/api/v1/stats/abc123
```

