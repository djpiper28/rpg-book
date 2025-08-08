import { EnvVarPort } from "../launcherTypes";
import { ProjectSvcClient } from "./pb/project.client";
import { LogLevel, LogProperty, LogRequest } from "./pb/system";
import { SystemSvcClient } from "./pb/system.client";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";

function mustHaveEnv(key: string): string {
  const val = globalThis.process.env[key];

  if (val) {
    return val;
  }

  throw new Error(`Cannot find env var ${key} - ${val} - ${JSON.stringify(globalThis.process.env)}`);
}

export const port = mustHaveEnv(EnvVarPort);
const transport = new GrpcWebFetchTransport({
  baseUrl: `https://127.0.0.1:${port}`,
  fetchInit: {
    keepalive: true,
    redirect: "error",
    cache: "no-cache",
  },
  format: "binary",
  timeout: 5_000,
});

export const systemClient = new SystemSvcClient(transport);
export const projectClient = new ProjectSvcClient(transport);
export const logger = {
  info: (msg: string, props: Record<string, string>) =>
    logAtLevel(LogLevel.INFO, msg, props),
  warn: (msg: string, props: Record<string, string>) =>
    logAtLevel(LogLevel.WARNING, msg, props),
  error: (msg: string, props: Record<string, string>) =>
    logAtLevel(LogLevel.ERROR, msg, props),
  fatal: (msg: string, props: Record<string, string>) =>
    logAtLevel(LogLevel.FATAL, msg, props),
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function logAtLevel(level: any, msg: string, props: Record<string, string>) {
  const properties = [];
  for (const key of Object.keys(props)) {
    const item = LogProperty.create({
      key: key,
      value: props[key],
    });

    properties.push(item);
  }

  const req = LogRequest.create({
    level: level,
    message: msg,
    // caller: logAtLevel.caller.caller.name,
    properties: properties,
  });

  console.log(`${level} - ${msg} - ${props}`);
  systemClient
    .log(req)
    .then(() => {})
    .catch((e) => {
      console.error("Cannot log:", e);
    });
}
