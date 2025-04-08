import { EnvVarCertificate, EnvVarPort } from "../launcherTypes";
import { credentials, load } from "grpc";

function mustHave(x: string | undefined): string {
  if (x) {
    return x;
  }

  throw new Error("Cannot find env vars");
}

export const certificate = mustHave(process.env[EnvVarCertificate]);
export const port = mustHave(process.env[EnvVarPort]);

const systemProto = load("../../../pubblic/protos/system.proto");
export const client = new systemProto.SystemSvc(
  `localhost:${port}`,
  credentials.createSsl(Buffer.from(certificate)),
);

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
  const req = new systemProto.LogRequest();
  req.setLevel(level);
  req.setCaller(logAtLevel.caller.caller.name);
  req.setMessage(msg);

  const properties: systemProto.LogProperty[] = [];
  for (const key of Object.keys(props)) {
    const item = new systemProto.LogProperty();
    item.setKey(key);
    item.setValue(props[key]);

    properties.push(item);
  }
  req.setPropertiesList(properties);

  console.log(`${level} - ${msg} - ${props}`);
  grpcClient.log(req, () => {});
}
