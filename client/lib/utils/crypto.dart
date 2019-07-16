import 'dart:convert';
import 'dart:math';
import 'dart:typed_data';

import 'package:crypt/crypt.dart';
import 'package:encrypt/encrypt.dart';
import 'package:password_hash/password_hash.dart';

String hashPassword(String password) {
  final salt = Salt.generateAsBase64String(saltSize);
  return Crypt.sha256(password, salt: salt).toString();
}

bool matchHashedPassword(String hashedPassword, String password) =>
    Crypt(hashedPassword).match(password);

String decrypt(String cipherText, String masterpass, [String privateKey]) {
  var cipherBytes = base64.decode(cipherText);

  if (privateKey == null) {
    final saltLength = cipherBytes[0];
    cipherBytes = cipherBytes.sublist(1);

    privateKey = base64.encode(cipherBytes.sublist(0, saltLength));
    cipherBytes = cipherBytes.sublist(saltLength);
  }

  final key = PBKDF2().generateKey(
    masterpass,
    privateKey,
    pbkdf2Rounds,
    keySize,
  );

  final ivBytes = cipherBytes.sublist(0, aesBlockSize);
  cipherBytes = cipherBytes.sublist(aesBlockSize);

  final iv = IV(ivBytes);
  final encrypter = Encrypter(AES(Key(key), mode: AESMode.cbc));
  return encrypter.decrypt(Encrypted(cipherBytes), iv: iv);
}

String encrypt(String plainText, String masterpass, [String privateKey]) {
  bool privateKeyWasEmpty = false;

  if (privateKey == null) {
    privateKey = Salt.generateAsBase64String(saltSize);
    privateKeyWasEmpty = true;
  }

  final key = PBKDF2().generateKey(
    masterpass,
    privateKey,
    pbkdf2Rounds,
    keySize,
  );

  final random = Random.secure();
  final ivBytes = List<int>.generate(aesBlockSize, (_) => random.nextInt(byteIntMax));
  final iv = IV(Uint8List.fromList(ivBytes));

  final encrypter = Encrypter(AES(Key(key), mode: AESMode.cbc));
  final cipherBytes = List<int>.from(encrypter.encrypt(plainText, iv: iv).bytes);
  cipherBytes.insertAll(0, ivBytes);

  if (privateKeyWasEmpty) {
    final base64PrivKey = base64.decode(privateKey);
    cipherBytes.insertAll(0, base64PrivKey);
    cipherBytes.insert(0, base64PrivKey.length);
  }

  return base64.encode(cipherBytes);
}

const saltSize = 16;
const pbkdf2Rounds = 4096;
const keySize = 32;
const aesBlockSize = 16;
const byteIntMax = 256;
