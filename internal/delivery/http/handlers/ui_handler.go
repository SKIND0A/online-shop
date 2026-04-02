package handlers

import "net/http"

type UIHandler struct{}

func NewUIHandler() *UIHandler {
	return &UIHandler{}
}

// AuthPage — главная страница: вход и регистрация.
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
  <title>Вход и регистрация</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      font-family: system-ui, -apple-system, "Segoe UI", Roboto, Arial, sans-serif;
      background: #fff;
      color: #111827;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px 16px;
    }
    .wrap {
      width: 100%;
      max-width: 400px;
    }
    h1 {
      font-size: 1.5rem;
      font-weight: 600;
      margin: 0 0 8px;
      text-align: center;
    }
    .sub {
      text-align: center;
      font-size: 0.9rem;
      color: #6b7280;
      margin-bottom: 28px;
    }
    .tabs {
      display: flex;
      gap: 0;
      margin-bottom: 24px;
      border-radius: 10px;
      padding: 4px;
      background: #f3f4f6;
    }
    .tabs button {
      flex: 1;
      border: none;
      padding: 10px 16px;
      font-size: 0.95rem;
      font-weight: 500;
      cursor: pointer;
      border-radius: 8px;
      background: transparent;
      color: #6b7280;
      transition: background 0.15s, color 0.15s;
    }
    .tabs button.active {
      background: #fff;
      color: #111827;
      box-shadow: 0 1px 2px rgba(0,0,0,.06);
    }
    .panel { display: none; }
    .panel.active { display: block; }
    label {
      display: block;
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
      margin-bottom: 6px;
    }
    input {
      width: 100%;
      padding: 10px 12px;
      font-size: 1rem;
      border: 1px solid #e5e7eb;
      border-radius: 8px;
      margin-bottom: 16px;
      background: #fff;
    }
    input:focus {
      outline: none;
      border-color: #2563eb;
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
    }
    .btn-primary {
      width: 100%;
      margin-top: 8px;
      padding: 12px 16px;
      font-size: 1rem;
      font-weight: 600;
      color: #fff;
      background: #2563eb;
      border: none;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-primary:hover { background: #1d4ed8; }
    .btn-primary:disabled {
      background: #93c5fd;
      cursor: not-allowed;
    }
    .msg {
      margin-top: 16px;
      padding: 12px 14px;
      border-radius: 8px;
      font-size: 0.875rem;
      line-height: 1.45;
      display: none;
    }
    .msg.show { display: block; }
    .msg.ok {
      background: #ecfdf5;
      color: #065f46;
      border: 1px solid #a7f3d0;
    }
    .msg.err {
      background: #fef2f2;
      color: #991b1b;
      border: 1px solid #fecaca;
    }
  </style>
</head>
<body>
  <div class="wrap">
    <h1>Online Shop</h1>
    <p class="sub">Войдите или создайте аккаунт</p>

    <div class="tabs" role="tablist">
      <button type="button" class="active" id="tabLogin" role="tab" aria-selected="true">Вход</button>
      <button type="button" id="tabRegister" role="tab" aria-selected="false">Регистрация</button>
    </div>

    <div id="panelLogin" class="panel active" role="tabpanel">
      <form id="formLogin" novalidate>
        <label for="lEmail">Email</label>
        <input id="lEmail" name="email" type="email" autocomplete="email" required placeholder="you@example.com">

        <label for="lPassword">Пароль</label>
        <input id="lPassword" name="password" type="password" autocomplete="current-password" required placeholder="••••••••">

        <button type="submit" class="btn-primary" id="btnLogin">Войти</button>
      </form>
    </div>

    <div id="panelRegister" class="panel" role="tabpanel" hidden>
      <form id="formRegister" novalidate>
        <label for="rName">Имя</label>
        <input id="rName" name="name" type="text" autocomplete="name" placeholder="Иван">

        <label for="rEmail">Email</label>
        <input id="rEmail" name="email" type="email" autocomplete="email" required placeholder="you@example.com">

        <label for="rPassword">Пароль</label>
        <input id="rPassword" name="password" type="password" autocomplete="new-password" required placeholder="Не менее 8 символов" minlength="8">

        <button type="submit" class="btn-primary" id="btnRegister">Зарегистрироваться</button>
      </form>
    </div>

    <div id="msg" class="msg" role="status"></div>
  </div>

  <script>
    const tabLogin = document.getElementById('tabLogin');
    const tabRegister = document.getElementById('tabRegister');
    const panelLogin = document.getElementById('panelLogin');
    const panelRegister = document.getElementById('panelRegister');
    const msg = document.getElementById('msg');

    function showMsg(text, ok) {
      msg.textContent = text;
      msg.className = 'msg show ' + (ok ? 'ok' : 'err');
    }

    function hideMsg() {
      msg.className = 'msg';
      msg.textContent = '';
    }

    function switchTab(register) {
      tabLogin.classList.toggle('active', !register);
      tabRegister.classList.toggle('active', register);
      tabLogin.setAttribute('aria-selected', !register);
      tabRegister.setAttribute('aria-selected', register);
      panelLogin.classList.toggle('active', !register);
      panelRegister.classList.toggle('active', register);
      panelRegister.hidden = !register;
      hideMsg();
    }

    tabLogin.addEventListener('click', () => switchTab(false));
    tabRegister.addEventListener('click', () => switchTab(true));

    async function postJSON(path, body) {
      const res = await fetch(path, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });
      const text = await res.text();
      let data;
      try { data = JSON.parse(text); } catch (_) { data = null; }
      return { res, data, text };
    }

    document.getElementById('formLogin').addEventListener('submit', async (e) => {
      e.preventDefault();
      hideMsg();
      const email = document.getElementById('lEmail').value.trim();
      const password = document.getElementById('lPassword').value;
      if (!email || !password) {
        showMsg('Введите email и пароль.', false);
        return;
      }
      const btn = document.getElementById('btnLogin');
      btn.disabled = true;
      try {
        const { res, data } = await postJSON('/api/v1/auth/login', { email, password });
        if (res.ok && data && data.success && data.data && data.data.access_token) {
          try { localStorage.setItem('access_token', data.data.access_token); } catch (_) {}
          window.location.href = '/shop';
          return;
        } else {
          const m = data && data.error && data.error.message ? data.error.message : 'Не удалось войти.';
          showMsg(m, false);
        }
      } catch (err) {
        showMsg('Сеть недоступна или сервер не отвечает.', false);
      }
      btn.disabled = false;
    });

    document.getElementById('formRegister').addEventListener('submit', async (e) => {
      e.preventDefault();
      hideMsg();
      const email = document.getElementById('rEmail').value.trim();
      const password = document.getElementById('rPassword').value;
      const name = document.getElementById('rName').value.trim();
      if (!email || password.length < 8) {
        showMsg('Укажите email и пароль не короче 8 символов.', false);
        return;
      }
      const btn = document.getElementById('btnRegister');
      btn.disabled = true;
      try {
        const { res, data } = await postJSON('/api/v1/auth/register', { email, password, name });
        if (res.ok && data && data.success) {
          showMsg('Аккаунт создан. Можно перейти на вкладку «Вход».', true);
        } else {
          const m = data && data.error && data.error.message ? data.error.message : 'Регистрация не удалась.';
          showMsg(m, false);
        }
      } catch (err) {
        showMsg('Сеть недоступна или сервер не отвечает.', false);
      }
      btn.disabled = false;
    });
  </script>
</body>
</html>`
