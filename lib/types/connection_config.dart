class ConnectionConfig {
  String host;
  String caCertificate;
  String certificate;
  String privateCertificate;

  ConnectionConfig({
    this.host,
    this.caCertificate,
    this.certificate,
    this.privateCertificate,
  });

  ConnectionConfig.fromJson(Map<String, dynamic> json) {
    host = json['host'];
    caCertificate = json['caCertificate'];
    certificate = json['certificate'];
    privateCertificate = json['privateCertificate'];
  }

  Map<String, dynamic> toJson() => {
        'host': host,
        'caCertificate': caCertificate,
        'certificate': certificate,
        'privateCertificate': privateCertificate,
      };
}
