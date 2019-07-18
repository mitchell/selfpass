import 'package:flutter/cupertino.dart';
import 'package:flutter/services.dart';
import 'package:otp/otp.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/credential.dart' as types;

import '../widgets/text_field.dart';

class Credential extends StatefulWidget {
  final types.Credential credential;

  const Credential(this.credential, {Key key}) : super(key: key);

  @override
  State createState() => _CredentialState(credential);
}

class _CredentialState extends State<Credential> {
  _CredentialControllers controllers;
  Map<String, _FieldBuildConfig> fieldMap;
  types.Credential credential;
  CredentialsRepo client;

  _CredentialState(this.credential) : super() {
    controllers = _CredentialControllers.fromCredential(credential);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    client = Provider.of<CredentialsRepo>(context);
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      navigationBar: CupertinoNavigationBar(
        trailing: CupertinoButton(
          child: Text(
            'Delete',
            style: TextStyle(color: CupertinoColors.destructiveRed),
          ),
          onPressed: makeDeleteHandler(context),
          padding: EdgeInsets.zero,
        ),
      ),
      child: Container(
        padding: const EdgeInsets.only(top: 15, bottom: 30, left: 30),
        child: ListView(
          children: buildFieldRows(context),
        ),
      ),
    );
  }

  Function makeDeleteHandler(BuildContext context) {
    return () {
      showCupertinoDialog(
        context: context,
        builder: (BuildContext context) => CupertinoAlertDialog(
          content: Text('Are you sure you want to delete this credential?'),
          actions: <Widget>[
            CupertinoDialogAction(
              child: Text('No'),
              isDefaultAction: true,
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
            CupertinoDialogAction(
              child: Text('Yes'),
              isDestructiveAction: true,
              onPressed: () async {
                await client.delete(credential.meta.id);
                Navigator.of(context).pushNamedAndRemoveUntil(
                  '/home',
                  ModalRoute.withName('/home'),
                );
              },
            ),
          ],
        ),
      );
    };
  }

  List<Widget> buildFieldRows(BuildContext context) {
    List<Widget> rows = [];

    fieldMap = buildFieldMap(controllers, credential);

    fieldMap.forEach((String prefix, _FieldBuildConfig config) {
      rows.add(Container(
        margin: EdgeInsets.only(top: 2.5),
        child: Text(prefix, style: TextStyle(fontWeight: FontWeight.w600)),
      ));

      final List<Widget> widgets = [
        Expanded(
          flex: 3,
          child: config.mutable
              ? TextField(
                  maxLines: 1,
                  controller: config.controller,
                  obscure: config.obscured)
              : Container(
                  margin: EdgeInsets.symmetric(vertical: 10),
                  child: Text(config.text)),
        ),
      ];

      if (config.copyable) {
        widgets.add(Expanded(
          child: CupertinoButton(
            child: Text(config.otp ? 'OTP' : 'Copy'),
            onPressed: () => Clipboard.setData(ClipboardData(
              text: config.otp
                  ? OTP
                      .generateTOTPCode(config.controller.text,
                          DateTime.now().millisecondsSinceEpoch)
                      .toString()
                  : config.mutable ? config.controller.text : config.text,
            )),
          ),
        ));
      }

      rows.add(Row(children: widgets));
    });

    return rows;
  }

  Map<String, _FieldBuildConfig> buildFieldMap(
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

    if (credential.meta.tag?.isNotEmpty ?? false) {
      fieldMap['Tag'] = _FieldBuildConfig(controller: controllers.tag);
    }

    if (credential.username?.isNotEmpty ?? false) {
      fieldMap['User:'] = _FieldBuildConfig(controller: controllers.username);
    }

    if (credential.email?.isNotEmpty ?? false) {
      fieldMap['Email:'] = _FieldBuildConfig(controller: controllers.email);
    }

    fieldMap['Password:'] =
        _FieldBuildConfig(controller: controllers.password, obscured: true);

    if (credential.otpSecret?.isNotEmpty ?? false) {
      fieldMap['OTP Key:'] = _FieldBuildConfig(
          controller: controllers.otpSecret, obscured: true, otp: true);
    }

    return fieldMap;
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
