import { EnvVarPort } from "../launcherTypes";
import { LogLevel, LogProperty, LogRequest } from "./pb/system";
import { ProjectSvcClient } from "./pb/project.client";
import { SystemSvcClient } from "./pb/system.client";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { getEnv } from "../electron";

let projectClient: ProjectSvcClient | undefined;
let systemClient: SystemSvcClient | undefined;

function mustHaveEnv(key: string): string {
  const val = getEnv()[key];

  if (!val) {
    throw new Error(
      `Cannot find env var ${key} - ${val} - ${JSON.stringify(globalThis.process.env)}`,
    );
  }

  return val;
}

export function initializeClients() {
  const port = mustHaveEnv(EnvVarPort);
  const transport = new GrpcWebFetchTransport({
    baseUrl: `https://127.0.0.1:${port}`,
    fetchInit: {
      keepalive: true,
      redirect: "error",
      cache: "no-cache",
      priority: "high",
    },
    format: "binary",
    timeout: 5_000,
  });

  systemClient = new SystemSvcClient(transport);
  projectClient = new ProjectSvcClient(transport);
  logger = {
    info: (msg: string, props: Record<string, string>) =>
      logAtLevel(LogLevel.INFO, msg, props),
    warn: (msg: string, props: Record<string, string>) =>
      logAtLevel(LogLevel.WARNING, msg, props),
    error: (msg: string, props: Record<string, string>) =>
      logAtLevel(LogLevel.ERROR, msg, props),
    fatal: (msg: string, props: Record<string, string>) =>
      logAtLevel(LogLevel.FATAL, msg, props),
  };

  const consoleLogger = (message: any, ...args: any) => {
    getLogger().info(JSON.stringify(message), {
      args: JSON.stringify(args),
      caller1: "console.log",
    });
  };

  const consoleErrorLogger = (message: any, ...args: any) => {
    getLogger().error(JSON.stringify(message), {
      args: JSON.stringify(args),
      caller1: "console.error",
    });
  };

  console.log = consoleLogger;
  console.error = consoleErrorLogger;
}

export function getProjectClient(): ProjectSvcClient {
  if (!projectClient) {
    throw new Error("ProjectClient not initialized");
  }

  return projectClient;
}

export function getSystemClient(): SystemSvcClient {
  if (!systemClient) {
    throw new Error("SystemClient not initialized");
  }

  return systemClient;
}

type logFunc = (msg: string, props: Record<string, string>) => void;

let logger:
  | {
      info: logFunc;
      warn: logFunc;
      error: logFunc;
      fatal: logFunc;
    }
  | undefined;

const clog = console.log;
const cerror = console.error;

export function getLogger(): {
  info: logFunc;
  warn: logFunc;
  error: logFunc;
  fatal: logFunc;
} {
  if (!logger) {
    throw new Error("Logger not initialized");
  }

  return logger;
}

function logAtLevel(
  level: LogLevel,
  msg: string,
  props: Record<string, string>,
) {
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

  clog(`${level} - ${msg} - ${JSON.stringify(props)}`);
  getSystemClient()
    .log(req)
    .then(() => {})
    .catch((e) => {
      cerror("Cannot log:", e);
    });
}
