// package: 
// file: system.proto

import * as jspb from "google-protobuf";
import * as common_pb from "./common_pb";

export class Settings extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
  }
}

export class LogRequest extends jspb.Message {
  getCaller(): string;
  setCaller(value: string): void;

  getLevel(): LogLevelMap[keyof LogLevelMap];
  setLevel(value: LogLevelMap[keyof LogLevelMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  clearPropertiesList(): void;
  getPropertiesList(): Array<LogProperty>;
  setPropertiesList(value: Array<LogProperty>): void;
  addProperties(value?: LogProperty, index?: number): LogProperty;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LogRequest): LogRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LogRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogRequest;
  static deserializeBinaryFromReader(message: LogRequest, reader: jspb.BinaryReader): LogRequest;
}

export namespace LogRequest {
  export type AsObject = {
    caller: string,
    level: LogLevelMap[keyof LogLevelMap],
    message: string,
    propertiesList: Array<LogProperty.AsObject>,
  }
}

export class LogProperty extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogProperty.AsObject;
  static toObject(includeInstance: boolean, msg: LogProperty): LogProperty.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LogProperty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogProperty;
  static deserializeBinaryFromReader(message: LogProperty, reader: jspb.BinaryReader): LogProperty;
}

export namespace LogProperty {
  export type AsObject = {
    key: string,
    value: string,
  }
}

export interface LogLevelMap {
  INFO: 0;
  WARNING: 1;
  ERROR: 2;
  FATAL: 3;
}

export const LogLevel: LogLevelMap;

