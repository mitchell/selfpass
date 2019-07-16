import 'dart:async';

import 'credential.dart';
import 'connection_config.dart';

abstract class CredentialsRepo {
  Stream<Metadata> getAllMetadata(String sourceHost);
  Future<Credential> get(String id);
  Future<Credential> create(CredentialInput input);
  Future<Credential> update(String id, CredentialInput input);
  Future<void> delete(String id);
}

abstract class ConfigRepo {
  Future<void> setPrivateKey(String key);
  Future<String> get privateKey;

  String get password;
  Future<void> setPassword(String password);
  Future<bool> get passwordIsSet;
  Future<bool> matchesPasswordHash(String password);

  Future<void> setConnectionConfig(ConnectionConfig config);
  Future<ConnectionConfig> get connectionConfig;

  Future<void> deleteAll();

  void reset();
}
