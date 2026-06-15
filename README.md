# SecureApp

Demo de 3 vulnerabilidades OWASP Top 10 y sus correcciones en Go.

## Estructura

```
secure-app-go/
├── vulnerable/    # Versión con SQLi, XSS, Command Injection
├── fixed/         # Versión corregida
└── .github/workflows/verify-fixes.yml  # Pipeline CI/CD
```

## Vulnerabilidades

### 1. SQL Injection → Prepared Statements

**vulnerable/main.go:**
```go
query := fmt.Sprintf("SELECT id FROM users WHERE username = '%s' AND password = '%s'", username, password)
row := db.QueryRow(query)
```

**fixed/main.go:**
```go
row := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password)
```

### 2. Reflected XSS → Escape de Output

**vulnerable/main.go:**
```go
fmt.Fprintf(w, "<h1>Search results for: %s</h1>", q)
```

**fixed/main.go:**
```go
w.Write([]byte("<h1>Search results for: " + html.EscapeString(q) + "</h1>"))
```

### 3. Command Injection → Validación + exec directo

**vulnerable/main.go:**
```go
out, err := exec.Command("bash", "-c", "ping -c1 "+host).Output()
```

**fixed/main.go:**
```go
var validHost = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)
if !validHost.MatchString(host) {
    http.Error(w, "Invalid host", http.StatusBadRequest)
    return
}
out, err := exec.Command("ping", "-c1", host).Output()
```
