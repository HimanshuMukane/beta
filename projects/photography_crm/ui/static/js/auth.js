(function () {
  const STORAGE_KEY = 'crm_user';

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
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
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

    async logout() {
      try {
        await fetch('/api/auth/logout', { method: 'POST' });
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
