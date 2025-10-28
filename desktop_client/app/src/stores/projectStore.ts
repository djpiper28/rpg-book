import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
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
    (set, get) => ({
      addCharacter: (
        handle: ProjectHandle,
        character: BasicCharacterDetails,
      ): void => {
        const projects = get();
        const project = projects.projects[asId(handle)]?.project;

        if (!project) {
          return;
        }

        const oldCharacter = project.characters.findIndex(
          (x) => x.handle?.id === character.handle?.id,
        );

        if (oldCharacter) {
          project.characters[oldCharacter] = structuredClone(character);
        } else {
          project.characters.push(character);
        }

        set(projects);
      },
      deleteCharacter: (
        handle: ProjectHandle,
        character: CharacterHandle,
      ): void => {
        const projects = get();
        const project = projects.projects[asId(handle)]?.project;

        if (!project) {
          return;
        }

        project.characters = project.characters.filter(
          (x) => x.handle?.id !== character.id,
        );

        set(projects);
      },
      getProject: (handle: ProjectHandle): Project | undefined => {
        const projects = get();
        return projects.projects[asId(handle)];
      },
      newProject: (handle: ProjectHandle, project: OpenProjectResp): void => {
        const projects = get();

        projects.projects[asId(handle)] = {
          handle,
          project,
        };

        set(projects);
      },
      projects: {},
      updateProject: (
        handle: ProjectHandle,
        project: OpenProjectResp,
      ): void => {
        const projects = get();
        const projectRef = projects.projects[asId(handle)];

        if (!projectRef) {
          return;
        }

        projectRef.project = project;
        set(projects);
      },
    }),
    {
      name: "project-storage",
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
);
