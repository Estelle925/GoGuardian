import { create } from 'zustand';
import request from '@/utils/request';

interface UserState {
  user: any;
  permissions: string[];
  menus: any[];
  setUser: (user: any) => void;
  setPermissions: (permissions: string[]) => void;
  setMenus: (menus: any[]) => void;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  fetchUserInfo: () => Promise<void>;
  fetchUserPermissions: () => Promise<void>;
  fetchUserMenus: () => Promise<void>;
}

const useUserStore = create<UserState>((set) => ({
  user: null,
  permissions: [],
  menus: [],

  setUser: (user) => set({ user }),
  setPermissions: (permissions) => set({ permissions }),
  setMenus: (menus) => set({ menus }),

  login: async (username, password) => {
    try {
      const response = await request.post('/login', { username, password });
      localStorage.setItem('token', response.token);
      set({ user: response.user });
    } catch (error) {
      throw error;
    }
  },

  logout: () => {
    localStorage.removeItem('token');
    set({ user: null, permissions: [], menus: [] });
  },

  fetchUserInfo: async () => {
    try {
      const user = await request.get('/users/info');
      set({ user });
    } catch (error) {
      throw error;
    }
  },

  fetchUserPermissions: async () => {
    try {
      const permissions = await request.get('/users/permissions');
      set({ permissions });
    } catch (error) {
      throw error;
    }
  },

  fetchUserMenus: async () => {
    try {
      const menus = await request.get('/users/menus');
      set({ menus });
    } catch (error) {
      throw error;
    }
  },
}));

export default useUserStore;