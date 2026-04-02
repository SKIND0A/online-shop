package handlers

import "net/http"

// ShopPage — витрина после входа: шапка с поиском и профилем.
func (h *UIHandler) ShopPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(shopPageHTML))
}

const shopPageHTML = `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Магазин</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      font-family: system-ui, -apple-system, "Segoe UI", Roboto, Arial, sans-serif;
      background: #fff;
      color: #111827;
    }
    .top {
      position: sticky;
      top: 0;
      z-index: 10;
      background: #fff;
      border-bottom: 1px solid #e5e7eb;
      padding: 12px 16px;
    }
    .top-inner {
      max-width: 1100px;
      margin: 0 auto;
      display: flex;
      align-items: center;
      gap: 12px;
      flex-wrap: wrap;
    }
    .logo {
      font-weight: 700;
      font-size: 1.125rem;
      color: #111827;
      text-decoration: none;
      margin-right: 4px;
    }
    .search-wrap {
      flex: 1;
      min-width: 200px;
      display: flex;
      align-items: center;
      gap: 10px;
    }
    .search-input {
      flex: 1;
      padding: 10px 14px;
      font-size: 1rem;
      border: 2px solid #2563eb;
      border-radius: 10px;
      background: #fff;
      color: #111827;
    }
    .search-input::placeholder { color: #9ca3af; }
    .search-input:focus {
      outline: none;
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
    }
    .profile-slot {
      position: relative;
      flex-shrink: 0;
    }
    .profile-btn {
      width: 44px;
      height: 44px;
      padding: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      border: 2px solid #2563eb;
      border-radius: 10px;
      background: #fff;
      color: #2563eb;
      cursor: pointer;
      transition: background 0.15s, color 0.15s;
    }
    .profile-btn:hover, .profile-btn[aria-expanded="true"] {
      background: #eff6ff;
      color: #1d4ed8;
    }
    .profile-btn svg {
      width: 22px;
      height: 22px;
    }
    .profile-menu {
      display: none;
      position: absolute;
      top: calc(100% + 8px);
      right: 0;
      min-width: 240px;
      background: #fff;
      border: 1px solid #e5e7eb;
      border-radius: 10px;
      box-shadow: 0 8px 24px rgba(0,0,0,.12);
      padding: 14px;
      z-index: 30;
    }
    .profile-menu.open { display: block; }
    .profile-menu-head {
      padding-bottom: 12px;
      margin-bottom: 12px;
      border-bottom: 1px solid #e5e7eb;
    }
    .profile-menu-head strong {
      display: block;
      font-size: 1rem;
      font-weight: 600;
      color: #111827;
      word-break: break-word;
    }
    .menu-email {
      display: block;
      margin-top: 4px;
      font-size: 0.8125rem;
      color: #6b7280;
      word-break: break-all;
    }
    .btn-logout {
      width: 100%;
      padding: 10px 14px;
      font-size: 0.9375rem;
      font-weight: 600;
      color: #fff;
      background: #dc2626;
      border: none;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-logout:hover { background: #b91c1c; }
    main {
      max-width: 1100px;
      margin: 0 auto;
      padding: 32px 16px 48px;
    }
    .hero {
      text-align: center;
      padding: 48px 16px;
      background: linear-gradient(180deg, #f8fafc 0%, #fff 100%);
      border-radius: 16px;
      border: 1px solid #e5e7eb;
    }
    .hero h1 {
      margin: 0 0 12px;
      font-size: 1.5rem;
      font-weight: 600;
    }
    .hero p {
      margin: 0;
      color: #6b7280;
      font-size: 1rem;
    }
  </style>
</head>
<body>
  <header class="top">
    <div class="top-inner">
      <a class="logo" href="/shop">Online Shop</a>
      <div class="search-wrap">
        <input type="search" class="search-input" id="shopSearch" placeholder="Поиск товаров…" autocomplete="off" aria-label="Поиск">
        <div class="profile-slot">
          <button type="button" class="profile-btn" id="profileBtn" aria-label="Профиль" aria-expanded="false" aria-haspopup="true" title="Профиль">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </button>
          <div class="profile-menu" id="profileMenu" role="menu" aria-hidden="true">
            <div class="profile-menu-head">
              <strong id="menuUserName">—</strong>
              <span class="menu-email" id="menuUserEmail"></span>
            </div>
            <button type="button" class="btn-logout" id="logoutBtn" role="menuitem">Выйти</button>
          </div>
        </div>
      </div>
    </div>
  </header>

  <main>
    <section class="hero">
      <h1>Добро пожаловать</h1>
      <p>Каталог товаров появится здесь. Пока можно пользоваться поиском в шапке.</p>
    </section>
  </main>

  <script>
    (function () {
      var token;
      try {
        token = localStorage.getItem('access_token');
        if (!token) {
          window.location.replace('/');
          return;
        }
      } catch (_) {
        window.location.replace('/');
        return;
      }

      function displayNameFromMe(me) {
        var n = (me.name || '').trim();
        if (n) return n;
        var e = me.email || '';
        var at = e.indexOf('@');
        return at > 0 ? e.slice(0, at) : 'Пользователь';
      }

      var profileBtn = document.getElementById('profileBtn');
      var profileMenu = document.getElementById('profileMenu');
      var menuUserName = document.getElementById('menuUserName');
      var menuUserEmail = document.getElementById('menuUserEmail');
      var logoutBtn = document.getElementById('logoutBtn');

      function setMenuOpen(open) {
        profileMenu.classList.toggle('open', open);
        profileBtn.setAttribute('aria-expanded', open ? 'true' : 'false');
        profileMenu.setAttribute('aria-hidden', open ? 'false' : 'true');
      }

      profileBtn.addEventListener('click', function (e) {
        e.stopPropagation();
        setMenuOpen(!profileMenu.classList.contains('open'));
      });

      document.addEventListener('click', function () {
        setMenuOpen(false);
      });

      profileMenu.addEventListener('click', function (e) {
        e.stopPropagation();
      });

      logoutBtn.addEventListener('click', function () {
        try { localStorage.removeItem('access_token'); } catch (_) {}
        window.location.href = '/';
      });

      fetch('/api/v1/auth/me', {
        headers: { 'Authorization': 'Bearer ' + token }
      })
        .then(function (res) {
          if (res.status === 401) {
            try { localStorage.removeItem('access_token'); } catch (_) {}
            window.location.replace('/');
            return null;
          }
          return res.json();
        })
        .then(function (body) {
          if (!body || !body.success || !body.data) return;
          var me = body.data;
          menuUserName.textContent = displayNameFromMe(me);
          menuUserEmail.textContent = me.email || '';
        })
        .catch(function () {
          try { localStorage.removeItem('access_token'); } catch (_) {}
          window.location.replace('/');
        });
    })();
  </script>
</body>
</html>`
