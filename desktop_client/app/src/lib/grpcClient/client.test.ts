import { getEnv } from "@/lib/electron";
import { getSystemClient, initializeClients } from "./client.ts";
import { vi } from "vitest";
import { EnvVarCertificate, EnvVarPort } from "../launcherTypes";

const mockedGetEnv = vi.mocked(getEnv);
vi.unmock("./client.ts");

describe("client", () => {
  beforeEach(() => {
    mockedGetEnv.mockClear();
  });

  it("should fail if it cannot load port and certificate", async () => {
    mockedGetEnv.mockReturnValueOnce({} as any);
    await expect(async () => initializeClients()).rejects.toThrowError();
  });

  it("should fail if it cannot connect to the server", async () => {
    mockedGetEnv.mockReturnValueOnce({
      [EnvVarPort]: "9000",
      [EnvVarCertificate]: "Testing-cert",
    });

    await expect(async () => {
      initializeClients();
      const version = getSystemClient().getVersion({});
      console.log(
        `The above call should have failed ${JSON.stringify(version)}`,
      );
    }).rejects.toThrowError();
  });
});
