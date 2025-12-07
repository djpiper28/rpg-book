import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";
import {
  type OpenProjectResp,
  type ProjectHandle,
} from "@/lib/grpcClient/pb/project";
import {
  type BasicCharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";

export interface Project {
  handle: ProjectHandle;
  project: OpenProjectResp;
}

interface ProjectStore {
  addCharacter: (
    handle: ProjectHandle,
    character: BasicCharacterDetails,
  ) => void;
  deleteCharacter: (handle: ProjectHandle, character: CharacterHandle) => void;
  getProject: (handle: ProjectHandle) => Project | undefined;
  newProject: (handle: ProjectHandle, project: OpenProjectResp) => void;
  projects: Record<string, Project | undefined>;
  updateProject: (handle: ProjectHandle, project: OpenProjectResp) => void;
}

function asId(handle: ProjectHandle): string {
  return handle.id;
}

export const useProjectStore = create<ProjectStore>()(
  persist(
    immer((set, get) => ({
      addCharacter: (
        handle: ProjectHandle,
        character: BasicCharacterDetails,
      ): void => {
        set((state) => {
          const project = state.projects[asId(handle)]?.project;

          if (!project) {
            return;
          }

          const oldCharacterIndex = project.characters.findIndex(
            (x) => x.handle?.id === character.handle?.id,
          );

          if (oldCharacterIndex === -1) {
            project.characters.push(character);
          } else {
            project.characters[oldCharacterIndex] = character;
          }
        });
      },
      deleteCharacter: (
        handle: ProjectHandle,
        character: CharacterHandle,
      ): void => {
        set((state) => {
          const project = state.projects[asId(handle)]?.project;

          if (!project) {
            return;
          }

          project.characters = project.characters.filter(
            (x) => x.handle?.id !== character.id,
          );
        });
      },
      getProject: (handle: ProjectHandle): Project | undefined => {
        return get().projects[asId(handle)];
      },
      newProject: (handle: ProjectHandle, project: OpenProjectResp): void => {
        set((state) => {
          state.projects[asId(handle)] = {
            handle,
            project,
          };
        });
      },
      projects: {},
      updateProject: (
        handle: ProjectHandle,
        project: OpenProjectResp,
      ): void => {
        set((state) => {
          const projectRef = state.projects[asId(handle)];

          if (!projectRef) {
            return;
          }

          projectRef.project = project;
        });
      },
    })),
    {
      name: "project-storage",
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
);
