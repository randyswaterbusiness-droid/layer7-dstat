{
  "version": 2,
  "framework": null,
  "routes": [
    { "src": "/target", "dest": "/api/index.go" },
    { "src": "/api/target", "dest": "/api/index.go" },
    { "src": "/(.*)", "dest": "/public/$1" }
  ]
}

