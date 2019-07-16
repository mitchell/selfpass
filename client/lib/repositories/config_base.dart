class ConfigBase {
  static const keyPrivateKey = "private_key";
  static const keyConnectionConfig = "connection_config";
  static const keyPassword = "password";

  bool passwordMatched = false;
  String _password;

  String get password {
    checkIfPasswordMatched();
    return _password;
  }

  set password(String password) => _password = password;

  void checkIfPasswordMatched() {
    if (passwordMatched) return;
    throw Exception('password not matched yet');
  }

  void reset() {
    passwordMatched = false;
    _password = null;
  }
}
