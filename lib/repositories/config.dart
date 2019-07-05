import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../types/abstracts.dart';
import '../types/connection_config.dart';

import '../utils/crypto.dart' as crypto;

class Config implements ConfigRepo {
  static const _keyPrivateKey = "private_key";
  static const _keyConnectionConfig = "connection_config";
  static const _keyPassword = "password";
  final FlutterSecureStorage _storage = FlutterSecureStorage();
  var _passwordMatched = false;

  Future<void> setPrivateKey(String key) {
    _checkIfPasswordMatched();
    return _storage.write(key: _keyPrivateKey, value: key);
  }

  Future<String> get privateKey {
    _checkIfPasswordMatched();
    return _storage.read(key: _keyPrivateKey);
  }

  Future<void> setPassword(String password) {
    _checkIfPasswordMatched();
    return _storage.write(
        key: _keyPassword, value: crypto.hashPassword(password));
  }

  Future<bool> get passwordSet async {
    var passHash = await _storage.read(key: _keyPassword);

    if (passHash != null) {
      return true;
    }

    _passwordMatched = true;

    return false;
  }

  Future<bool> matchesPasswordHash(String password) async =>
      _passwordMatched = crypto.matchHashedPassword(
          await _storage.read(key: _keyPassword), password);

  Future<void> setConnectionConfig(ConnectionConfig config) {
    _checkIfPasswordMatched();
    return _storage.write(
        key: _keyConnectionConfig, value: json.encode(config));
  }

  Future<ConnectionConfig> get connectionConfig async {
    _checkIfPasswordMatched();
    final connConfig = await _storage.read(key: _keyConnectionConfig);

    if (connConfig == null) {
      return null;
    }

    return ConnectionConfig.fromJson(json.decode(connConfig));
  }

  Future<void> deleteAll() {
    _checkIfPasswordMatched();
    return _storage.deleteAll();
  }

  void _checkIfPasswordMatched() {
    if (_passwordMatched) return;
    throw Exception('password not matched yet');
  }
}
