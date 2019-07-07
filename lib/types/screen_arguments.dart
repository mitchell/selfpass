import 'connection_config.dart';

class ConfigScreenArguments {
  final ConnectionConfig connectionConfig;
  final String privateKey;

  const ConfigScreenArguments(this.connectionConfig, this.privateKey);
}
