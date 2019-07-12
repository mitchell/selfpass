import 'package:selfpass_protobuf/credentials.pb.dart' as protobuf;

class Metadata {
  String id;
  String sourceHost;
  DateTime createdAt;
  DateTime updatedAt;
  String primary;
  String loginUrl;
  String tag;

  Metadata({
    this.id,
    this.sourceHost,
    this.createdAt,
    this.updatedAt,
    this.primary,
    this.loginUrl,
    this.tag,
  });

  Metadata.fromProtobuf(protobuf.Metadata metadata) {
    id = metadata.id;
    createdAt = metadata.createdAt.toDateTime();
    updatedAt = metadata.updatedAt.toDateTime();
    sourceHost = metadata.sourceHost;
    primary = metadata.primary;
    loginUrl = metadata.loginUrl;
    tag = metadata.tag;
  }

  @override
  String toString() => "id: $id";
}

class MetadataInput {
  String sourceHost;
  String primary;
  String loginUrl;
  String tag;

  MetadataInput({this.sourceHost, this.primary, this.loginUrl, this.tag});
}

class Credential {
  Metadata meta;
  String username;
  String email;
  String password;
  String otpSecret;

  Credential({
    this.meta,
    this.username,
    this.email,
    this.password,
    this.otpSecret,
  });

  Credential.fromProtobuf(protobuf.Credential credential) {
    meta = Metadata(
      id: credential.id,
      createdAt: credential.createdAt.toDateTime(),
      updatedAt: credential.updatedAt.toDateTime(),
      sourceHost: credential.sourceHost,
      primary: credential.primary,
      loginUrl: credential.loginUrl,
      tag: credential.tag,
    );
    username = credential.username;
    email = credential.email;
    password = credential.password;
    otpSecret = credential.otpSecret;
  }

  @override
  String toString() => "meta: $meta\n";
}

class CredentialInput {
  MetadataInput meta;
  String username;
  String email;
  String password;
  String otpSecret;

  CredentialInput({
    this.meta,
    this.username,
    this.email,
    this.password,
    this.otpSecret,
  });
}
