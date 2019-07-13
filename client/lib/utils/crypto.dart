import 'dart:math';
import 'dart:convert';

import 'package:crypt/crypt.dart';
import 'package:encrypt/encrypt.dart';
import 'package:password_hash/password_hash.dart';

String hashPassword(String password) {
  const saltSize = 16;
  const saltIntMax = 256;

  final random = Random.secure();
  final saltInts =
  List<int>.generate(saltSize, (_) => random.nextInt(saltIntMax));
  final salt = base64.encode(saltInts);

  return Crypt.sha256(password, salt: salt).toString();
}

bool matchHashedPassword(String hashedPassword, String password) =>
    Crypt(hashedPassword).match(password);

String decryptPassword(String masterpass, privateKey, cipherText) {
  final key = PBKDF2().generateKey(
    masterpass, privateKey, pbkdf2Rounds, keySize,
  );

  var cipherBytes = base64.decode(cipherText);
  final ivBytes = cipherBytes.sublist(0, aesBlockSize);
  cipherBytes = cipherBytes.sublist(aesBlockSize);

  final iv = IV(ivBytes);
  final encrypter = Encrypter(AES(Key(key), mode: AESMode.cbc));
  return encrypter.decrypt(Encrypted(cipherBytes), iv: iv);
}

const pbkdf2Rounds = 4096;
const keySize = 32;
const aesBlockSize = 16;
