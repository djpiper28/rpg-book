import { create } from "zustand";

export interface RpgBookError {
  body: string;
  title?: string;
}

interface GlobalErrorStore {
  currentError?: RpgBookError;
  setError: (error: RpgBookError) => void;
}

export const useGlobalErrorStore = create<GlobalErrorStore>()(
  (set): GlobalErrorStore => ({
    setError: (error: RpgBookError) => {
      set({ currentError: error });
    },
  }),
);
