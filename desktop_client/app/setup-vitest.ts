import "@testing-library/jest-dom/vitest";
import { EnvVarCertificate, EnvVarPort } from "@/lib/launcherTypes/index.ts";
import MatchMediaMock from "vitest-matchmedia-mock";

process.env = {
  [EnvVarPort]: "9000",
  [EnvVarCertificate]: "Testing-cert",
};

let matchMediaMock = new MatchMediaMock();
afterAll(() => {
  matchMediaMock.clear();
});

vi.mock("./src/lib/electron/index.ts", {
  electronShell: {
    openExternal: async (url: string) => {
      console.log(`Open external called ${url}`);
    },
  },
  electronDialog: {
    showSaveDialog: async (
      opts: Electron.SaveDialogOptions,
    ): Promise<Electron.SaveDialogReturnValue> => {
      console.log(`Save dialog called ${JSON.stringify(opts)}`);
      return Promise.resolve({
        canceled: false,
        filePath: "Mock file path",
      });
    },
  },
  env: {},
  getSystemVersion: () => "Mocked OS",
  getBuildType: () => "Test",
});

vi.mock("@electron/remote", () => {});
