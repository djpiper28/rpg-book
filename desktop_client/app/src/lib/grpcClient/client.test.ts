import { describe, it, expect } from "vitest";
import { EnvVarCertificate, EnvVarPort } from "../launcherTypes";

const baseEnv = {
  APP_ROOT: "",
  VITE_PUBLIC: "",
};

describe("client", () => {
  it("should fail if it cannot load port and certificate", async () => {
    process.env = { ...baseEnv };
    const res = import("./client.ts");
    await expect(res).rejects.toThrowError();
  });

  it("should fail if it cannot connect to the server", async () => {
    process.env = {
      [EnvVarPort]: "9000",
      [EnvVarCertificate]: "Testing-cert",
      ...baseEnv,
    };
    const res = import("./client.ts");
    await expect(res).rejects.toThrowError();
  });
});
