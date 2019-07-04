import 'dart:math';
import 'dart:convert';

import 'package:crypt/crypt.dart';

String hashPassword(String password) {
  final random = Random.secure();
  final saltInts = List<int>.generate(16, (_) => random.nextInt(256));
  final salt = base64.encode(saltInts);

  return Crypt.sha256(password, salt: salt).toString();
}

bool matchHashedPassword(String hashedPassword, String password) =>
    Crypt(hashedPassword).match(password);
