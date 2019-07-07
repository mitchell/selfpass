import 'package:flutter/cupertino.dart';
import 'package:flutter/services.dart';

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

  _CredentialState(this._credential) : super() {
    _controllers = _CredentialControllers.fromCredential(_credential);
    _fieldMap = _buildFieldMap(_controllers, _credential);
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      navigationBar: CupertinoNavigationBar(),
      child: Container(
        margin: const EdgeInsets.only(top: 30, left: 30),
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
      'Created: ': _FieldBuildConfig(
        mutable: false,
        copyable: false,
        text: credential.meta.createdAt.toString(),
      ),
      'Updated: ': _FieldBuildConfig(
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

    return fieldMap;
  }

  List<Widget> _buildFieldRows(BuildContext context) {
    List<Widget> rows = [];

    _fieldMap.forEach((key, value) {
      rows.add(Container(
        margin: EdgeInsets.only(top: 10),
        child: Text(key, style: TextStyle(fontWeight: FontWeight.w600)),
      ));

      final List<Widget> widgets = [
        Expanded(
          flex: 3,
          child: value.mutable
              ? TextField(maxLines: 1, controller: value.controller)
              : Container(
                  margin: EdgeInsets.symmetric(vertical: 10),
                  child: Text(value.text),
                ),
        ),
      ];

      if (value.copyable) {
        widgets.add(Flexible(
          child: CupertinoButton(
            child: Text('Copy'),
            onPressed: () => Clipboard.setData(ClipboardData(
              text: value.mutable ? value.controller.text : value.text,
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

  const _FieldBuildConfig({
    this.mutable = true,
    this.copyable = true,
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
      );
}
