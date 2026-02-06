import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { type ProjectHandle } from "@/lib/grpcClient/pb/project";
import { type CharacterHandle } from "@/lib/grpcClient/pb/project_character";
import { type NoteHandle } from "@/lib/grpcClient/pb/project_note";

interface Tab {
  handle: ProjectHandle;
  name: string;
  selectedCharacter?: CharacterHandle;
  selectedNote?: NoteHandle;
  // Currently selected tab
  value: string;
}

interface TabStore {
  addTab: (handle: ProjectHandle, name: string) => void;
  removeTab: (handle: ProjectHandle) => void;
  selected: {
    setSelectedCharacter: (tabId: string, handle?: CharacterHandle) => void;
    setSelectedNote: (tabId: string, handle?: NoteHandle) => void;
  };
  selectedTab?: ProjectHandle;
  setSelectedTab: (handle: ProjectHandle) => void;
  setValue: (tabId: string, value: string) => void;
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
      selected: {
        setSelectedCharacter: (
          tabId: string,
          handle?: CharacterHandle,
        ): void => {
          const newTabs = structuredClone(get().tabs);
          newTabs[tabId].selectedCharacter = handle;
          set({ tabs: newTabs });
        },
        setSelectedNote: (tabId: string, handle?: NoteHandle): void => {
          const newTabs = structuredClone(get().tabs);
          newTabs[tabId].selectedNote = handle;
          set({ tabs: newTabs });
        },
      },
      selectedTab: undefined,
      setSelectedTab: (handle: ProjectHandle): void => {
        set({
          selectedTab: handle,
        });
      },
      setValue: (tabId: string, value: string): void => {
        const newTabs = structuredClone(get().tabs);
        newTabs[tabId].value = value;

        set({
          tabs: newTabs,
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
