library repositories;

import 'dart:convert';
import 'dart:async';
import 'dart:io';

import 'package:grpc/grpc.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:selfpass_protobuf/credentials.pbgrpc.dart' as grpc;
import 'package:selfpass_protobuf/credentials.pb.dart' as protobuf;

import '../types/abstracts.dart';
import '../types/connection_config.dart';
import '../types/credential.dart';

import '../utils/crypto.dart' as crypto;

part 'config_base.dart';
part 'encrypted_shared_preferences.dart';
part 'grpc_credentials_client.dart';