import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import {
  type OpenProjectResp,
  type ProjectHandle,
} from "@/lib/grpcClient/pb/project";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";

export interface Project {
  handle: ProjectHandle;
  project: OpenProjectResp;
}

interface ProjectStore {
  addCharacter: (
    handle: ProjectHandle,
    character: BasicCharacterDetails,
  ) => void;
  getProject: (handle: ProjectHandle) => Project | undefined;
  newProject: (handle: ProjectHandle, project: OpenProjectResp) => void;
  projects: Record<string, Project | undefined>;
  updateProject: (handle: ProjectHandle, project: OpenProjectResp) => void;
}

function asId(handle: ProjectHandle) {
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
        projects.projects[asId(handle)]?.project.characters.push(character);
        set(projects);
      },
      getProject: (handle: ProjectHandle) => {
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

        if (!projects.projects[asId(handle)]) {
          return;
        }

        projects.projects[asId(handle)].project = project;
        set(projects);
      },
    }),
    {
      name: "project-storage",
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
);
