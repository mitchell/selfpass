part of 'repositories.dart';

class EncryptedSharedPreferences extends ConfigBase implements ConfigRepo {
  @override
  Future<ConnectionConfig> get connectionConfig async {
    checkIfPasswordMatched();

    final prefs = await SharedPreferences.getInstance();
    final cipherText = prefs.getString(ConfigBase.keyConnectionConfig);

    if (cipherText == null) return null;

    final configJson = crypto.decrypt(cipherText, _password);

    return ConnectionConfig.fromJson(json.decode(configJson));
  }

  @override
  Future<void> deleteAll() async {
    checkIfPasswordMatched();

    final prefs = await SharedPreferences.getInstance();

    prefs.remove(ConfigBase.keyConnectionConfig);
    prefs.remove(ConfigBase.keyPassword);
    prefs.remove(ConfigBase.keyPrivateKey);
  }

  @override
  Future<bool> matchesPasswordHash(String password) async {
    final prefs = await SharedPreferences.getInstance();

    passwordMatched = crypto.matchHashedPassword(
      prefs.getString(ConfigBase.keyPassword),
      password,
    );

    if (passwordMatched) _password = password;

    return passwordMatched;
  }

  @override
  Future<bool> get passwordIsSet async {
    final prefs = await SharedPreferences.getInstance();

    final isSet = prefs.containsKey(ConfigBase.keyPassword);
    passwordMatched = !isSet;

    return isSet;
  }

  @override
  Future<String> get privateKey async {
    checkIfPasswordMatched();

    final prefs = await SharedPreferences.getInstance();
    final cipherText = prefs.getString(ConfigBase.keyPrivateKey);

    return crypto.decrypt(cipherText, _password);
  }

  @override
  Future<void> setConnectionConfig(ConnectionConfig config) async {
    checkIfPasswordMatched();

    final prefs = await SharedPreferences.getInstance();

    final configJson = json.encode(config);

    prefs.setString(
      ConfigBase.keyConnectionConfig,
      crypto.encrypt(configJson, _password),
    );
  }

  @override
  Future<void> setPassword(String password) async {
    final prefs = await SharedPreferences.getInstance();

    _password = password;
    passwordMatched = true;

    prefs.setString(ConfigBase.keyPassword, crypto.hashPassword(password));
  }

  @override
  Future<void> setPrivateKey(String key) async {
    final prefs = await SharedPreferences.getInstance();

    prefs.setString(ConfigBase.keyPrivateKey, crypto.encrypt(key, _password));
  }
}
