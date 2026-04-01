package handlers

import "net/http"

type UIHandler struct{}

func NewUIHandler() *UIHandler {
	return &UIHandler{}
}

func (h *UIHandler) AuthPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(authPageHTML))
}

const authPageHTML = `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Auth Debug UI</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 24px; background: #f7f7f7; color: #111; }
    .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; max-width: 1000px; }
    .card { background: #fff; border: 1px solid #ddd; border-radius: 8px; padding: 16px; }
    h1 { margin-top: 0; }
    label { display: block; margin-top: 10px; font-size: 14px; }
    input { width: 100%; padding: 8px; margin-top: 4px; box-sizing: border-box; }
    button { margin-top: 12px; padding: 10px 14px; cursor: pointer; }
    pre { background: #111; color: #0f0; padding: 12px; border-radius: 6px; min-height: 180px; overflow: auto; }
  </style>
</head>
<body>
  <h1>Проверка регистрации/логина</h1>
  <p>Эта страница шлет запросы в <code>/api/v1/auth/register</code> и <code>/api/v1/auth/login</code>.</p>

  <div class="grid">
    <div class="card">
      <h3>Register</h3>
      <label>Email <input id="rEmail" type="email" value="user@example.com"></label>
      <label>Password <input id="rPassword" type="password" value="StrongPass123!"></label>
      <label>Name <input id="rName" type="text" value="Ivan"></label>
      <button id="registerBtn">POST /register</button>
    </div>

    <div class="card">
      <h3>Login</h3>
      <label>Email <input id="lEmail" type="email" value="user@example.com"></label>
      <label>Password <input id="lPassword" type="password" value="StrongPass123!"></label>
      <button id="loginBtn">POST /login</button>
    </div>
  </div>

  <div class="card" style="margin-top:16px; max-width:1000px;">
    <h3>Ответ сервера</h3>
    <pre id="output">Ожидание запроса...</pre>
  </div>

  <script>
    const output = document.getElementById('output');

    async function callAPI(path, payload) {
      output.textContent = 'Запрос: ' + path + '\n\n' + JSON.stringify(payload, null, 2);
      try {
        const res = await fetch(path, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        const text = await res.text();
        let parsed = text;
        try { parsed = JSON.parse(text); } catch (_) {}
        output.textContent =
          'HTTP ' + res.status + '\n\n' +
          (typeof parsed === 'string' ? parsed : JSON.stringify(parsed, null, 2));
      } catch (e) {
        output.textContent = 'Ошибка запроса: ' + e;
      }
    }

    document.getElementById('registerBtn').addEventListener('click', () => {
      callAPI('/api/v1/auth/register', {
        email: document.getElementById('rEmail').value,
        password: document.getElementById('rPassword').value,
        name: document.getElementById('rName').value
      });
    });

    document.getElementById('loginBtn').addEventListener('click', () => {
      callAPI('/api/v1/auth/login', {
        email: document.getElementById('lEmail').value,
        password: document.getElementById('lPassword').value
      });
    });
  </script>
</body>
</html>`
