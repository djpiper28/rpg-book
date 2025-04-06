// GENERATED CODE -- DO NOT EDIT!

// package: 
// file: settings.proto

import * as settings_pb from "./settings_pb";
import * as common_pb from "./common_pb";
import * as grpc from "@grpc/grpc-js";

interface ISettingsSvcService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  getSettings: grpc.MethodDefinition<common_pb.Empty, settings_pb.Settings>;
}

export const SettingsSvcService: ISettingsSvcService;

export interface ISettingsSvcServer extends grpc.UntypedServiceImplementation {
  getSettings: grpc.handleUnaryCall<common_pb.Empty, settings_pb.Settings>;
}

export class SettingsSvcClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  getSettings(argument: common_pb.Empty, callback: grpc.requestCallback<settings_pb.Settings>): grpc.ClientUnaryCall;
  getSettings(argument: common_pb.Empty, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<settings_pb.Settings>): grpc.ClientUnaryCall;
  getSettings(argument: common_pb.Empty, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<settings_pb.Settings>): grpc.ClientUnaryCall;
}
