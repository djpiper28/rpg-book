// GENERATED CODE -- DO NOT EDIT!

// package: 
// file: system.proto

import * as system_pb from "./system_pb";
import * as common_pb from "./common_pb";
import * as grpc from "@grpc/grpc-js";

interface ISystemSvcService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  getSettings: grpc.MethodDefinition<common_pb.Empty, system_pb.Settings>;
  log: grpc.MethodDefinition<system_pb.LogRequest, common_pb.Empty>;
}

export const SystemSvcService: ISystemSvcService;

export interface ISystemSvcServer extends grpc.UntypedServiceImplementation {
  getSettings: grpc.handleUnaryCall<common_pb.Empty, system_pb.Settings>;
  log: grpc.handleUnaryCall<system_pb.LogRequest, common_pb.Empty>;
}

export class SystemSvcClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  getSettings(argument: common_pb.Empty, callback: grpc.requestCallback<system_pb.Settings>): grpc.ClientUnaryCall;
  getSettings(argument: common_pb.Empty, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<system_pb.Settings>): grpc.ClientUnaryCall;
  getSettings(argument: common_pb.Empty, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<system_pb.Settings>): grpc.ClientUnaryCall;
  log(argument: system_pb.LogRequest, callback: grpc.requestCallback<common_pb.Empty>): grpc.ClientUnaryCall;
  log(argument: system_pb.LogRequest, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<common_pb.Empty>): grpc.ClientUnaryCall;
  log(argument: system_pb.LogRequest, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<common_pb.Empty>): grpc.ClientUnaryCall;
}
