import 'package:flutter/cupertino.dart';
import 'package:flutter/services.dart';
import 'package:otp/otp.dart';

import '../types/credential.dart' as types;

import '../widgets/text_field.dart';

class Credential extends StatefulWidget {
  final types.Credential credential;

  const Credential(this.credential, {Key key}) : super(key: key);

  @override
  State createState() => _CredentialState(credential);
}

class _CredentialState extends State<Credential> {
  _CredentialControllers _controllers;
  Map<String, _FieldBuildConfig> _fieldMap;
  types.Credential _credential;

  _CredentialState(this._credential) : super();

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      navigationBar: CupertinoNavigationBar(),
      child: Container(
        padding: const EdgeInsets.only(top: 15, bottom: 30, left: 30),
        child: ListView(
          children: _buildFieldRows(context),
        ),
      ),
    );
  }

  Map<String, _FieldBuildConfig> _buildFieldMap(
    _CredentialControllers controllers,
    types.Credential credential,
  ) {
    final fieldMap = {
      'Id:': _FieldBuildConfig(mutable: false, text: credential.meta.id),
      'Created:': _FieldBuildConfig(
        mutable: false,
        copyable: false,
        text: credential.meta.createdAt.toString(),
      ),
      'Updated:': _FieldBuildConfig(
        mutable: false,
        copyable: false,
        text: credential.meta.updatedAt.toString(),
      ),
      'Host:': _FieldBuildConfig(controller: controllers.sourceHost),
      'Primary:': _FieldBuildConfig(controller: controllers.primary),
    };

    if (credential.meta.tag != null && credential.meta.tag != '') {
      fieldMap['Tag'] = _FieldBuildConfig(controller: controllers.tag);
    }

    if (credential.username != null && credential.username != '') {
      fieldMap['User:'] = _FieldBuildConfig(controller: controllers.username);
    }

    if (credential.email != null && credential.email != '') {
      fieldMap['Email:'] = _FieldBuildConfig(controller: controllers.email);
    }

    fieldMap['Password:'] =
        _FieldBuildConfig(controller: controllers.password, obscured: true);

    if (credential.otpSecret != null && credential.otpSecret != '') {
      fieldMap['OTP Secret:'] = _FieldBuildConfig(
          controller: controllers.otpSecret, obscured: true, otp: true);
    }

    return fieldMap;
  }

  List<Widget> _buildFieldRows(BuildContext context) {
    List<Widget> rows = [];

    _controllers = _CredentialControllers.fromCredential(_credential);
    _fieldMap = _buildFieldMap(_controllers, _credential);

    _fieldMap.forEach((key, value) {
      rows.add(Container(
        margin: EdgeInsets.only(top: 2.5),
        child: Text(key, style: TextStyle(fontWeight: FontWeight.w600)),
      ));

      final List<Widget> widgets = [
        Expanded(
          flex: 3,
          child: value.mutable
              ? TextField(
                  maxLines: 1,
                  controller: value.controller,
                  obscure: value.obscured)
              : Container(
                  margin: EdgeInsets.symmetric(vertical: 10),
                  child: Text(value.text)),
        ),
      ];

      if (value.copyable) {
        widgets.add(Expanded(
          child: CupertinoButton(
            child: Text(value.otp ? 'OTP' : 'Copy'),
            onPressed: () => Clipboard.setData(ClipboardData(
              text: value.otp
                  ? OTP
                      .generateTOTPCode(value.controller.text,
                          DateTime.now().millisecondsSinceEpoch)
                      .toString()
                  : value.mutable ? value.controller.text : value.text,
            )),
          ),
        ));
      }

      rows.add(Row(children: widgets));
    });

    return rows;
  }
}

class _FieldBuildConfig {
  final TextEditingController controller;
  final String text;
  final bool mutable;
  final bool copyable;
  final bool obscured;
  final bool otp;

  const _FieldBuildConfig({
    this.mutable = true,
    this.copyable = true,
    this.obscured = false,
    this.otp = false,
    this.controller,
    this.text,
  });
}

class _CredentialControllers {
  final TextEditingController sourceHost;
  final TextEditingController primary;
  final TextEditingController tag;
  final TextEditingController username;
  final TextEditingController email;
  final TextEditingController password;
  final TextEditingController otpSecret;

  const _CredentialControllers({
    this.sourceHost,
    this.primary,
    this.tag,
    this.username,
    this.email,
    this.password,
    this.otpSecret,
  });

  factory _CredentialControllers.fromCredential(types.Credential credential) =>
      _CredentialControllers(
        sourceHost: TextEditingController(text: credential.meta.sourceHost),
        primary: TextEditingController(text: credential.meta.primary),
        tag: TextEditingController(text: credential.meta.tag),
        username: TextEditingController(text: credential.username),
        email: TextEditingController(text: credential.email),
        password: TextEditingController(text: credential.password),
        otpSecret: TextEditingController(text: credential.otpSecret),
      );
}
