import { create } from "zustand";
import { type ProjectHandle } from "@/lib/grpcClient/pb/project";

interface Tab {
  handle: ProjectHandle;
  name: string;
}

interface TabStore {
  addTab: (handle: ProjectHandle, name: string) => void;
  removeTab: (handle: ProjectHandle) => void;
  selectedTab?: ProjectHandle;
  tabs: Record<string, Tab>;
}

export const useTabStore = create<TabStore>((set) => ({
  addTab: (handle: ProjectHandle, name: string) => {
    set((state) => {
      state.tabs[handle.id] = {
        handle,
        name,
      };

      return {
        tabs: state.tabs,
      };
    });
  },
  removeTab: (handle: ProjectHandle) => {
    set((state) => {
      // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
      delete state.tabs[handle.id];

      return {
        tabs: state.tabs,
      };
    });
  },
  selectedTab: undefined,
  tabs: {},
}));
