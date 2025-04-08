import { describe, it, expect } from "vitest";
import { EnvVarCertificate, EnvVarPort } from "../launcherTypes";

describe("client", () => {
  it("should fail if it cannot load port and certificate", async () => {
    process.env = {};
    const res = import("./client.ts");
    await expect(res).rejects.toThrowError();
  });

  it("should fail if it cannot connect to the server", async () => {
    process.env = {
      [EnvVarPort]: "9000",
      [EnvVarCertificate]: "Testing-cert",
    };
    const res = import("./client.ts");
    await expect(res).rejects.toThrowError();
  });
});
