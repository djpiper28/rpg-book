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
import { type Note, type NoteHandle } from "@/lib/grpcClient/pb/project_note";

export interface Project {
  handle: ProjectHandle;
  project: OpenProjectResp;
}

interface CharacterMethods {
  addCharacter: (
    handle: ProjectHandle,
    character: BasicCharacterDetails,
  ) => void;
  deleteCharacter: (handle: ProjectHandle, character: CharacterHandle) => void;
}

interface NoteMethods {
  addNote: (handle: ProjectHandle, note: Note) => void;
  deleteNote: (handle: ProjectHandle, note: NoteHandle) => void;
}

interface ProjectMethods {
  getProject: (handle: ProjectHandle) => Project | undefined;
  newProject: (handle: ProjectHandle, project: OpenProjectResp) => void;
  projects: Record<string, Project | undefined>;
  updateProject: (handle: ProjectHandle, project: OpenProjectResp) => void;
}

type ProjectStore = ProjectMethods & CharacterMethods & NoteMethods;

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
      addNote: (handle: ProjectHandle, note: Note): void => {
        set((state) => {
          const project = state.projects[asId(handle)]?.project;

          if (!project) {
            return;
          }

          const oldNoteIndex = project.notes.findIndex(
            (x) => x.handle?.id === note.handle?.id,
          );

          if (oldNoteIndex === -1) {
            project.notes.push(note);
          } else {
            project.notes[oldNoteIndex] = note;
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
      deleteNote: (handle: ProjectHandle, note: NoteHandle): void => {
        set((state) => {
          const project = state.projects[asId(handle)]?.project;

          if (!project) {
            return;
          }

          project.notes = project.notes.filter((x) => x.handle?.id !== note.id);
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
