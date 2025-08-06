import { EnvVarCertificate, EnvVarPort } from "@/lib/launcherTypes/index.ts";
import MatchMediaMock from "vitest-matchmedia-mock";
import * as remote from "@electron/remote";

process.env = {
  [EnvVarPort]: "9000",
  [EnvVarCertificate]: "Testing-cert",
};

let matchMediaMock = new MatchMediaMock();
afterAll(() => {
  matchMediaMock.clear();
});

vi.mock("@electron/remote", () => ({ exec: vi.fn() }));
