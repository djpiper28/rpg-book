import { EnvVarCertificate, EnvVarPort } from "../launcherTypes";

const baseEnv = {
  APP_ROOT: "",
  RPG_BOOK_CERTIFICATE: undefined,
  RPG_BOOK_PORT: undefined,
  VITE_PUBLIC: "",
};

describe("client", () => {
  it("should fail if it cannot load port and certificate", async () => {
    process.env = { ...baseEnv } as any;
    const res = import("./client.ts");
    await expect(res).rejects.toThrowError();
  });

  it("should fail if it cannot connect to the server", async () => {
    process.env = {
      [EnvVarPort]: "9000",
      [EnvVarCertificate]: "Testing-cert",
      ...baseEnv,
    } as any;
    const res = import("./client.ts");
    await expect(res).rejects.toThrowError();
  });
});
