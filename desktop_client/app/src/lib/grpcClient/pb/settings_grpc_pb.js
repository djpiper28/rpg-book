// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var settings_pb = require('./settings_pb.js');
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

function serialize_Settings(arg) {
  if (!(arg instanceof settings_pb.Settings)) {
    throw new Error('Expected argument of type Settings');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Settings(buffer_arg) {
  return settings_pb.Settings.deserializeBinary(new Uint8Array(buffer_arg));
}


var SettingsSvcService = exports.SettingsSvcService = {
  getSettings: {
    path: '/SettingsSvc/GetSettings',
    requestStream: false,
    responseStream: false,
    requestType: common_pb.Empty,
    responseType: settings_pb.Settings,
    requestSerialize: serialize_Empty,
    requestDeserialize: deserialize_Empty,
    responseSerialize: serialize_Settings,
    responseDeserialize: deserialize_Settings,
  },
};

exports.SettingsSvcClient = grpc.makeGenericClientConstructor(SettingsSvcService, 'SettingsSvc');
