(function () {
  const STORAGE_KEY = 'crm_user';

  function getCsrfToken() {
    const match = document.cookie.match(new RegExp('(^|;\\s*)app_csrf=([^;]+)'));
    return match ? decodeURIComponent(match[2]) : '';
  }

  window.CRMAuth = {
    getUser() {
      try {
        const data = localStorage.getItem(STORAGE_KEY);
        return data ? JSON.parse(data) : null;
      } catch (e) {
        return null;
      }
    },

    setUser(user) {
      if (!user) {
        localStorage.removeItem(STORAGE_KEY);
      } else {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(user));
      }
    },

    clearUser() {
      localStorage.removeItem(STORAGE_KEY);
    },

    async login(email, password) {
      const headers = {
        'Content-Type': 'application/json',
      };
      const csrf = getCsrfToken();
      if (csrf) {
        headers['X-CSRF-Token'] = csrf;
      }

      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        throw new Error(data.error || 'Failed to sign in');
      }

      if (data.user) {
        this.setUser(data.user);
      }

      return data;
    },

    async signup(firstName, lastName, email, password) {
      const headers = {
        'Content-Type': 'application/json',
      };
      const csrf = getCsrfToken();
      if (csrf) {
        headers['X-CSRF-Token'] = csrf;
      }

      const response = await fetch('/api/auth/signup', {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({
          first_name: firstName,
          last_name: lastName,
          email: email,
          password: password,
        }),
      });

      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        throw new Error(data.error || 'Failed to create account');
      }

      if (data.user) {
        this.setUser(data.user);
      }

      return data;
    },

    async logout() {
      try {
        const headers = {};
        const csrf = getCsrfToken();
        if (csrf) {
          headers['X-CSRF-Token'] = csrf;
        }
        await fetch('/api/auth/logout', { method: 'POST', headers: headers });
      } finally {
        this.clearUser();
        window.location.href = '/signin';
      }
    },

    async checkMe() {
      try {
        const response = await fetch('/api/auth/me');
        if (response.ok) {
          const data = await response.json();
          if (data.user) {
            this.setUser(data.user);
            return data.user;
          }
        } else {
          this.clearUser();
          return null;
        }
      } catch (e) {
        return null;
      }
    }
  };
})();
