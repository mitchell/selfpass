import 'dart:math';
import 'dart:convert';
import 'dart:typed_data';

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

String decryptPassword(String masterpass, privateKey, ciphertext) {
  final key =
      PBKDF2().generateKey(masterpass, privateKey, pbkdf2Rounds, keySize);

  var cipherbytes = base64.decode(ciphertext);
  final iv =
      IV(Uint8List.fromList(cipherbytes.getRange(0, aesBlockSize).toList()));

  cipherbytes = Uint8List.fromList(
      cipherbytes.getRange(aesBlockSize, cipherbytes.length).toList());

  final encrypter = Encrypter(AES(Key(key), mode: AESMode.cbc));

  return encrypter.decrypt(Encrypted(cipherbytes), iv: iv);
}

const pbkdf2Rounds = 4096;
const keySize = 32;
const aesBlockSize = 16;
