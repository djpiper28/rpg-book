import "@testing-library/jest-dom/vitest";
import { EnvVarCertificate, EnvVarPort } from "@/lib/launcherTypes/index.ts";
import MatchMediaMock from "vitest-matchmedia-mock";
import { vi } from "vitest";

process.env = {
  [EnvVarPort]: "9000",
  [EnvVarCertificate]: "Testing-cert",
};

let matchMediaMock = new MatchMediaMock();
beforeEach(() => {
  matchMediaMock.clear();
});

const mockElectron = {
  getEnv: vi.fn(() => {
    console.log("getEnv mock called");
    return {};
  }),
  dialog: {
    showSaveDialog: vi.fn(
      async (
        opts: Electron.SaveDialogOptions,
      ): Promise<Electron.SaveDialogReturnValue> => {
        console.log(`Save dialog called ${JSON.stringify(opts)}`);
        return Promise.resolve({
          canceled: false,
          filePath: "Mock file path",
        });
      },
    ),
  },
  getBuildType: vi.fn(() => "Test"),
  getSystemVersion: vi.fn(() => "Mocked OS"),
  shell: {
    openExternal: vi.fn(async (url: string) => {
      console.log(`Open external called ${url}`);
    }),
  },
};

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
window.electron = mockElectron;
vi.mock("@/lib/electron", () => mockElectron);

vi.mock("@/lib/grpcClient/client", () => ({
  getLogger: () => ({
    info: (msg: string, props) => console.log("INFO:", msg, props),
    warn: (msg: string, props) => console.log("WARN:", msg, props),
    error: (msg: string, props) => console.log("ERROR:", msg, props),
    fatal: (msg: string, props) => console.log("FATAL:", msg, props),
  }),
  getSystemClient: vi.fn(() => ({
    log: vi.fn(),
  })),
  getProjectClient: vi.fn(() => ({
    createProject: vi.fn(),
  })),
  initializeClients: vi.fn(),
}));

vi.mock("@electron/remote", () => {});
