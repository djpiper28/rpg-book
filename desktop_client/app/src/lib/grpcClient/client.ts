import { EnvVarPort } from "../launcherTypes";
require("google-closure-library");
import * as system from "./pb/system_pb";
import { SystemSvcClient } from "./pb/SystemServiceClientPb";
const { LogLevel, LogProperty, LogRequest } = system;

function mustHave(x: string | undefined): string {
  if (x) {
    return x;
  }

  throw new Error("Cannot find env vars");
}

export const port = mustHave(process.env[EnvVarPort]);
export const client = new SystemSvcClient(`localhost:${port}`, null, null);
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
  const req = new LogRequest();
  req.setLevel(level);
  req.setCaller(logAtLevel.caller.caller.name);
  req.setMessage(msg);

  const properties = [];
  for (const key of Object.keys(props)) {
    const item = new LogProperty();
    item.setKey(key);
    item.setValue(props[key]);

    properties.push(item);
  }
  req.setPropertiesList(properties);

  console.log(`${level} - ${msg} - ${props}`);
  client.log(req);
}
