import { create } from "zustand";
import { type Settings } from "@/lib/grpcClient/pb/system";

const defaultSettings: Settings = {
  compressImages: false,
  darkMode: true,
  devMode: false,
};

interface SettingsStore {
  setSettings: (settings: Settings) => void;
  settings: Settings;
}

export const useSettingsStore = create<SettingsStore>(
  (set): SettingsStore => ({
    setSettings: (settings: Settings): void => {
      set({ settings });
      console.log(settings);
    },
    settings: defaultSettings,
  }),
);
