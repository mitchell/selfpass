import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:grpc/grpc.dart';
import 'package:selfpass_protobuf/credentials.pbgrpc.dart' as grpc;
import 'package:selfpass_protobuf/credentials.pb.dart' as protobuf;

import '../types/abstracts.dart';
import '../types/connection_config.dart';
import '../types/credential.dart';

class GRPCCredentialsClient implements CredentialsRepo {
  static GRPCCredentialsClient _cached;
  grpc.CredentialsClient _client;

  GRPCCredentialsClient(ConnectionConfig config) {
    final caCert = utf8.encode(config.caCertificate);
    final cert = utf8.encode(config.certificate);
    final privateCert = utf8.encode(config.privateCertificate);

    final splitHost = config.host.split(':');
    final hostname = splitHost[0];
    final port = int.parse(splitHost[1]);

    _client = grpc.CredentialsClient(ClientChannel(
      hostname,
      port: port,
      options: ChannelOptions(
        credentials: _ChannelCredentials(caCert, cert, privateCert),
      ),
    ));
  }

  factory GRPCCredentialsClient.cached({ConnectionConfig config}) =>
      config == null ? _cached : _cached = GRPCCredentialsClient(config);

  Stream<Metadata> getAllMetadata(String sourceHost) {
    final request = grpc.SourceHostRequest();
    request.sourceHost = sourceHost;

    return _client.getAllMetadata(request).map<Metadata>(
        (protobuf.Metadata pbMetadata) => Metadata.fromProtobuf(pbMetadata));
  }

  Future<Credential> get(String id) async {
    final request = grpc.IdRequest();
    request.id = id;

    return Credential.fromProtobuf(await _client.get(request));
  }

  Future<Credential> create(CredentialInput input) async {
    return Credential();
  }

  Future<Credential> update(String id, CredentialInput input) async {
    return Credential();
  }

  Future<void> delete(String id) {
    final request = grpc.IdRequest();
    request.id = id;

    return _client.delete(request);
  }
}

class _ChannelCredentials extends ChannelCredentials {
  final List<int> _key;
  final List<int> _cert;

  const _ChannelCredentials(List<int> caCert, this._cert, this._key)
      : super.secure(certificates: caCert);

  @override
  SecurityContext get securityContext {
    return super.securityContext
      ..usePrivateKeyBytes(_key)
      ..useCertificateChainBytes(_cert);
  }
}
