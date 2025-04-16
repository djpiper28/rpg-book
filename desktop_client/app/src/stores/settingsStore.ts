import { Settings } from "@/lib/grpcClient/pb/system";
import { create } from "zustand";

const defaultSettings: Settings = {
  devMode: false,
  darkMode: true,
};

interface SettingsStore {
  settings: Settings;
  setSettings: (settings: Settings) => void;
}

export const useSettingsStore = create<SettingsStore>(
  (set): SettingsStore => ({
    settings: defaultSettings,
    setSettings: (settings: Settings) => {
      set({ settings: settings });
      console.log(settings);
    },
  }),
);
