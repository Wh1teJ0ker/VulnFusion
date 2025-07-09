// src/store/user.js
import { create } from 'zustand';

export const useUserStore = create((set) => ({
    token: '',
    username: '',
    role: '',

    setUser: ({ token, username, role }) => set({ token, username, role }),

    clearUser: () => set({ token: '', username: '', role: '' }),
}));
