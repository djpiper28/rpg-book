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
