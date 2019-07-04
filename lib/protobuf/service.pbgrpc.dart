///
//  Generated code. Do not modify.
//  source: service.proto
///
// ignore_for_file: camel_case_types,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'dart:core' as $core show int, String, List;

import 'package:grpc/service_api.dart' as $grpc;
import 'service.pb.dart' as $0;
export 'service.pb.dart';

class CredentialServiceClient extends $grpc.Client {
  static final _$getAllMetadata =
      $grpc.ClientMethod<$0.GetAllMetadataRequest, $0.Metadata>(
          '/selfpass.credentials.CredentialService/GetAllMetadata',
          ($0.GetAllMetadataRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) => $0.Metadata.fromBuffer(value));
  static final _$get = $grpc.ClientMethod<$0.IdRequest, $0.Credential>(
      '/selfpass.credentials.CredentialService/Get',
      ($0.IdRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.Credential.fromBuffer(value));
  static final _$create =
      $grpc.ClientMethod<$0.CredentialRequest, $0.Credential>(
          '/selfpass.credentials.CredentialService/Create',
          ($0.CredentialRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) => $0.Credential.fromBuffer(value));
  static final _$update = $grpc.ClientMethod<$0.UpdateRequest, $0.Credential>(
      '/selfpass.credentials.CredentialService/Update',
      ($0.UpdateRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.Credential.fromBuffer(value));
  static final _$delete = $grpc.ClientMethod<$0.IdRequest, $0.DeleteResponse>(
      '/selfpass.credentials.CredentialService/Delete',
      ($0.IdRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.DeleteResponse.fromBuffer(value));

  CredentialServiceClient($grpc.ClientChannel channel,
      {$grpc.CallOptions options})
      : super(channel, options: options);

  $grpc.ResponseStream<$0.Metadata> getAllMetadata(
      $0.GetAllMetadataRequest request,
      {$grpc.CallOptions options}) {
    final call = $createCall(
        _$getAllMetadata, $async.Stream.fromIterable([request]),
        options: options);
    return $grpc.ResponseStream(call);
  }

  $grpc.ResponseFuture<$0.Credential> get($0.IdRequest request,
      {$grpc.CallOptions options}) {
    final call = $createCall(_$get, $async.Stream.fromIterable([request]),
        options: options);
    return $grpc.ResponseFuture(call);
  }

  $grpc.ResponseFuture<$0.Credential> create($0.CredentialRequest request,
      {$grpc.CallOptions options}) {
    final call = $createCall(_$create, $async.Stream.fromIterable([request]),
        options: options);
    return $grpc.ResponseFuture(call);
  }

  $grpc.ResponseFuture<$0.Credential> update($0.UpdateRequest request,
      {$grpc.CallOptions options}) {
    final call = $createCall(_$update, $async.Stream.fromIterable([request]),
        options: options);
    return $grpc.ResponseFuture(call);
  }

  $grpc.ResponseFuture<$0.DeleteResponse> delete($0.IdRequest request,
      {$grpc.CallOptions options}) {
    final call = $createCall(_$delete, $async.Stream.fromIterable([request]),
        options: options);
    return $grpc.ResponseFuture(call);
  }
}

abstract class CredentialServiceBase extends $grpc.Service {
  $core.String get $name => 'selfpass.credentials.CredentialService';

  CredentialServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetAllMetadataRequest, $0.Metadata>(
        'GetAllMetadata',
        getAllMetadata_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.GetAllMetadataRequest.fromBuffer(value),
        ($0.Metadata value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.IdRequest, $0.Credential>(
        'Get',
        get_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.IdRequest.fromBuffer(value),
        ($0.Credential value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CredentialRequest, $0.Credential>(
        'Create',
        create_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CredentialRequest.fromBuffer(value),
        ($0.Credential value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRequest, $0.Credential>(
        'Update',
        update_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateRequest.fromBuffer(value),
        ($0.Credential value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.IdRequest, $0.DeleteResponse>(
        'Delete',
        delete_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.IdRequest.fromBuffer(value),
        ($0.DeleteResponse value) => value.writeToBuffer()));
  }

  $async.Stream<$0.Metadata> getAllMetadata_Pre(
      $grpc.ServiceCall call, $async.Future request) async* {
    yield* getAllMetadata(call, (await request) as $0.GetAllMetadataRequest);
  }

  $async.Future<$0.Credential> get_Pre(
      $grpc.ServiceCall call, $async.Future request) async {
    return get(call, await request);
  }

  $async.Future<$0.Credential> create_Pre(
      $grpc.ServiceCall call, $async.Future request) async {
    return create(call, await request);
  }

  $async.Future<$0.Credential> update_Pre(
      $grpc.ServiceCall call, $async.Future request) async {
    return update(call, await request);
  }

  $async.Future<$0.DeleteResponse> delete_Pre(
      $grpc.ServiceCall call, $async.Future request) async {
    return delete(call, await request);
  }

  $async.Stream<$0.Metadata> getAllMetadata(
      $grpc.ServiceCall call, $0.GetAllMetadataRequest request);
  $async.Future<$0.Credential> get(
      $grpc.ServiceCall call, $0.IdRequest request);
  $async.Future<$0.Credential> create(
      $grpc.ServiceCall call, $0.CredentialRequest request);
  $async.Future<$0.Credential> update(
      $grpc.ServiceCall call, $0.UpdateRequest request);
  $async.Future<$0.DeleteResponse> delete(
      $grpc.ServiceCall call, $0.IdRequest request);
}
