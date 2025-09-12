import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { type ProjectHandle } from "@/lib/grpcClient/pb/project";

interface Tab {
  handle: ProjectHandle;
  name: string;
}

interface TabStore {
  addTab: (handle: ProjectHandle, name: string) => void;
  removeTab: (handle: ProjectHandle) => void;
  selectedTab?: ProjectHandle;
  setSelectedTab: (handle: ProjectHandle) => void;
  tabs: Record<string, Tab>;
}

export const useTabStore = create<TabStore>()(
  persist(
    (set, get) => ({
      addTab: (handle: ProjectHandle, name: string): void => {
        const newTabs = structuredClone(get().tabs);

        newTabs[handle.id] = {
          handle,
          name,
        };

        set({
          selectedTab: handle,
          tabs: newTabs,
        });
      },
      removeTab: (handle: ProjectHandle): void => {
        const newTabs = structuredClone(get().tabs);
        // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
        delete newTabs[handle.id];

        set({
          tabs: newTabs,
        });
      },
      selectedTab: undefined,
      setSelectedTab: (handle: ProjectHandle): void => {
        set({
          selectedTab: handle,
        });
      },
      tabs: {},
    }),
    {
      name: "tabs-storage",
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
);
