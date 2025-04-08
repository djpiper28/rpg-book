import * as jspb from 'google-protobuf'

import * as common_pb from './common_pb'; // proto import: "common.proto"


export class Settings extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
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
  setCaller(value: string): LogRequest;

  getLevel(): LogLevel;
  setLevel(value: LogLevel): LogRequest;

  getMessage(): string;
  setMessage(value: string): LogRequest;

  getPropertiesList(): Array<LogProperty>;
  setPropertiesList(value: Array<LogProperty>): LogRequest;
  clearPropertiesList(): LogRequest;
  addProperties(value?: LogProperty, index?: number): LogProperty;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LogRequest): LogRequest.AsObject;
  static serializeBinaryToWriter(message: LogRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogRequest;
  static deserializeBinaryFromReader(message: LogRequest, reader: jspb.BinaryReader): LogRequest;
}

export namespace LogRequest {
  export type AsObject = {
    caller: string,
    level: LogLevel,
    message: string,
    propertiesList: Array<LogProperty.AsObject>,
  }
}

export class LogProperty extends jspb.Message {
  getKey(): string;
  setKey(value: string): LogProperty;

  getValue(): string;
  setValue(value: string): LogProperty;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogProperty.AsObject;
  static toObject(includeInstance: boolean, msg: LogProperty): LogProperty.AsObject;
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

export enum LogLevel { 
  INFO = 0,
  WARNING = 1,
  ERROR = 2,
  FATAL = 3,
}
