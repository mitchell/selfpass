///
//  Generated code. Do not modify.
//  source: service.proto
///
// ignore_for_file: camel_case_types,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name

const DeleteResponse$json = const {
  '1': 'DeleteResponse',
  '2': const [
    const {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

const GetAllMetadataRequest$json = const {
  '1': 'GetAllMetadataRequest',
  '2': const [
    const {'1': 'source_host', '3': 1, '4': 1, '5': 9, '10': 'sourceHost'},
  ],
};

const IdRequest$json = const {
  '1': 'IdRequest',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

const UpdateRequest$json = const {
  '1': 'UpdateRequest',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    const {
      '1': 'credential',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.selfpass.credentials.CredentialRequest',
      '10': 'credential'
    },
  ],
};

const DumpResponse$json = const {
  '1': 'DumpResponse',
  '2': const [
    const {'1': 'contents', '3': 1, '4': 1, '5': 12, '10': 'contents'},
  ],
};

const EmptyRequest$json = const {
  '1': 'EmptyRequest',
};

const Metadata$json = const {
  '1': 'Metadata',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    const {
      '1': 'created_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    const {
      '1': 'updated_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
    const {'1': 'primary', '3': 4, '4': 1, '5': 9, '10': 'primary'},
    const {'1': 'source_host', '3': 5, '4': 1, '5': 9, '10': 'sourceHost'},
    const {'1': 'login_url', '3': 6, '4': 1, '5': 9, '10': 'loginUrl'},
    const {'1': 'tag', '3': 7, '4': 1, '5': 9, '10': 'tag'},
  ],
};

const Credential$json = const {
  '1': 'Credential',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    const {
      '1': 'created_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    const {
      '1': 'updated_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
    const {'1': 'primary', '3': 4, '4': 1, '5': 9, '10': 'primary'},
    const {'1': 'username', '3': 5, '4': 1, '5': 9, '10': 'username'},
    const {'1': 'email', '3': 6, '4': 1, '5': 9, '10': 'email'},
    const {'1': 'password', '3': 7, '4': 1, '5': 9, '10': 'password'},
    const {'1': 'source_host', '3': 8, '4': 1, '5': 9, '10': 'sourceHost'},
    const {'1': 'login_url', '3': 9, '4': 1, '5': 9, '10': 'loginUrl'},
    const {'1': 'tag', '3': 10, '4': 1, '5': 9, '10': 'tag'},
    const {'1': 'otp_secret', '3': 11, '4': 1, '5': 9, '10': 'otpSecret'},
  ],
};

const CredentialRequest$json = const {
  '1': 'CredentialRequest',
  '2': const [
    const {'1': 'primary', '3': 1, '4': 1, '5': 9, '10': 'primary'},
    const {'1': 'username', '3': 2, '4': 1, '5': 9, '10': 'username'},
    const {'1': 'email', '3': 3, '4': 1, '5': 9, '10': 'email'},
    const {'1': 'password', '3': 4, '4': 1, '5': 9, '10': 'password'},
    const {'1': 'source_host', '3': 5, '4': 1, '5': 9, '10': 'sourceHost'},
    const {'1': 'login_url', '3': 6, '4': 1, '5': 9, '10': 'loginUrl'},
    const {'1': 'tag', '3': 7, '4': 1, '5': 9, '10': 'tag'},
    const {'1': 'otp_secret', '3': 8, '4': 1, '5': 9, '10': 'otpSecret'},
  ],
};
