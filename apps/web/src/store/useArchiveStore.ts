import { create } from 'zustand';
import { api } from '../api/client';

export interface Series {
    id: string;
    title: string;
    author: string;
    cover_color: string;
    description: string;
}

export interface LoreEntry {
    id: string;
    series_id: string;
    title: string;
    category: string;
    content: string;
    metadata: Record<string, unknown>;
    created_at: string;
}

interface ArchiveState {
    // Inventory state (Admin key)
    adminToken: string | null;
    isAdmin: boolean;
    authError: string | null;
    isLoading: boolean;

    // World Archive Content Data
    seriesList: Series[];
    activeSeriesEntries: LoreEntry[];

    // Spatial Interaction Nodes (For R3F Camera Targeting)
    focusedSeriesId: string | null;
    focusedEntryId: string | null;

    // Auth Operations
    loginAdmin: (credentials: { email: string; password: string }) => Promise<boolean>;
    logoutAdmin: () => void;
    clearAuthError: () => void;

    // Archive Loading Operations
    fetchSeriesList: () => Promise<void>;
    fetchSeriesDetails: (seriesId: string) => Promise<void>;

    // Interaction Focus Mutators
    setFocusedSeries: (seriesId: string | null) => void;
    setFocusedEntry: (entryId: string | null) => void;
}

export const useArchiveStore = create<ArchiveState>((set, get) => ({
    // Initialize state directly from browser runtime memory
    adminToken: localStorage.getItem('admin_key'),
    isAdmin: !!localStorage.getItem('admin_key'),
    authError: null,
    isLoading: false,

    seriesList: [],
    activeSeriesEntries: [],

    focusedSeriesId: null,
    focusedEntryId: null,

    loginAdmin: async (credentials) => {
        set({ isLoading: true, authError: null });
        try {
            const data = await api.post('/auth/login', { body: credentials }) as { token: string };
            localStorage.setItem('admin_key', data.token);
            set({ adminToken: data.token, isAdmin: true, isLoading: false });
            return true;
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Authentication failed';
            set({ authError: errorMessage, isLoading: false, isAdmin: false });
            return false;
        }
    },

    logoutAdmin: () => {
        localStorage.removeItem('admin_key');
        set({ adminToken: null, isAdmin: false, focusedSeriesId: null, focusedEntryId: null });
    },

    clearAuthError: () => set({ authError: null }),

    fetchSeriesList: async () => {
        try {
            const data = await api.get('/series') as Series[];
            set({ seriesList: data || [] });
        } catch (err) {
            console.error('Failed to sync bookshelves with backend archive context:', err);
        }
    },

    fetchSeriesDetails: async (seriesId) => {
        try {
            const data = await api.get(`/series/${seriesId}`) as { entries: LoreEntry[] };
            // Our Go backend returns an object containing { series: ..., entries: [...] }
            set({ activeSeriesEntries: data.entries || [] });
        } catch (err) {
            console.error(`Failed to pull layout map data for series node ${seriesId}:`, err);
        }
    },

    setFocusedSeries: (seriesId) => {
        set({ focusedSeriesId: seriesId });
        if (seriesId) {
            get().fetchSeriesDetails(seriesId);
        } else {
            set({ activeSeriesEntries: [] });
        }
    },

    setFocusedEntry: (entryId) => set({ focusedEntryId: entryId }),
}));