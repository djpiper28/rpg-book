// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var system_pb = require('./system_pb.js');
var common_pb = require('./common_pb.js');

function serialize_Empty(arg) {
  if (!(arg instanceof common_pb.Empty)) {
    throw new Error('Expected argument of type Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Empty(buffer_arg) {
  return common_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_LogRequest(arg) {
  if (!(arg instanceof system_pb.LogRequest)) {
    throw new Error('Expected argument of type LogRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_LogRequest(buffer_arg) {
  return system_pb.LogRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_Settings(arg) {
  if (!(arg instanceof system_pb.Settings)) {
    throw new Error('Expected argument of type Settings');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Settings(buffer_arg) {
  return system_pb.Settings.deserializeBinary(new Uint8Array(buffer_arg));
}


var SystemSvcService = exports.SystemSvcService = {
  getSettings: {
    path: '/SystemSvc/GetSettings',
    requestStream: false,
    responseStream: false,
    requestType: common_pb.Empty,
    responseType: system_pb.Settings,
    requestSerialize: serialize_Empty,
    requestDeserialize: deserialize_Empty,
    responseSerialize: serialize_Settings,
    responseDeserialize: deserialize_Settings,
  },
  log: {
    path: '/SystemSvc/Log',
    requestStream: false,
    responseStream: false,
    requestType: system_pb.LogRequest,
    responseType: common_pb.Empty,
    requestSerialize: serialize_LogRequest,
    requestDeserialize: deserialize_LogRequest,
    responseSerialize: serialize_Empty,
    responseDeserialize: deserialize_Empty,
  },
};

exports.SystemSvcClient = grpc.makeGenericClientConstructor(SystemSvcService, 'SystemSvc');
