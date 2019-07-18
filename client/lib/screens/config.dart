import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/connection_config.dart';

import '../widgets/text_field.dart';

class Config extends StatefulWidget {
  final ConnectionConfig connectionConfig;
  final String privateKey;

  const Config(this.connectionConfig, this.privateKey, {Key key})
      : super(key: key);

  @override
  State createState() => _ConfigState(this.connectionConfig, this.privateKey);
}

class _ConfigState extends State<Config> {
  TextEditingController hostController;
  TextEditingController caCertController;
  TextEditingController certController;
  TextEditingController privateCertController;
  TextEditingController privateKeyController;
  ConnectionConfig connectionConfig;
  String privateKey;
  ConfigRepo config;

  _ConfigState(this.connectionConfig, this.privateKey) {
    if (connectionConfig == null) {
      connectionConfig = ConnectionConfig();
    }

    hostController = TextEditingController(text: connectionConfig.host);
    certController = TextEditingController(text: connectionConfig.certificate);
    caCertController =
        TextEditingController(text: connectionConfig.caCertificate);
    privateCertController =
        TextEditingController(text: connectionConfig.privateCertificate);

    privateKeyController = TextEditingController(text: privateKey);
  }

  @override
  didChangeDependencies() async {
    super.didChangeDependencies();

    config = Provider.of<ConfigRepo>(context);
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      navigationBar: CupertinoNavigationBar(
        trailing: connectionConfig?.host == null
            ? null
            : CupertinoButton(
                child: Text(
                  'Reset',
                  style: TextStyle(color: CupertinoColors.destructiveRed),
                ),
                onPressed: buildResetAllHandler(context),
                padding: EdgeInsets.zero,
              ),
      ),
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 30),
        child: ListView(children: [
          Container(margin: EdgeInsets.only(top: 10), child: Text('Host:')),
          TextField(maxLines: 1, controller: hostController),
          Container(
              margin: EdgeInsets.only(top: 5), child: Text('Private key:')),
          TextField(maxLines: 1, controller: privateKeyController),
          Container(
              margin: EdgeInsets.only(top: 5), child: Text('CA certificate:')),
          TextField(maxLines: 5, controller: caCertController),
          Container(
              margin: EdgeInsets.only(top: 5),
              child: Text('Client certificate:')),
          TextField(maxLines: 5, controller: certController),
          Container(
              margin: EdgeInsets.only(top: 5),
              child: Text('Private certificate:')),
          TextField(maxLines: 5, controller: privateCertController),
          Container(
            padding: EdgeInsets.symmetric(vertical: 20),
            margin: EdgeInsets.symmetric(horizontal: 70),
            child: CupertinoButton.filled(
              child: Text('Save'),
              onPressed: makeSaveOnPressed(context),
            ),
          ),
        ]),
      ),
    );
  }

  GestureTapCallback buildResetAllHandler(BuildContext context) {
    return () {
      showCupertinoDialog(
        context: context,
        builder: (BuildContext context) => CupertinoAlertDialog(
          content: Text(
            'Are you sure you want to delete all config values and lock the app?',
          ),
          actions: [
            CupertinoDialogAction(
              isDefaultAction: true,
              child: Text('No'),
              onPressed: () => Navigator.of(context).pop(),
            ),
            CupertinoDialogAction(
              isDestructiveAction: true,
              child: Text('Yes'),
              onPressed: () async {
                connectionConfig = null;
                await config.deleteAll();
                Navigator.of(context)
                    .pushNamedAndRemoveUntil('/', ModalRoute.withName('/home'));
              },
            ),
          ],
        ),
      );
    };
  }

  VoidCallback makeSaveOnPressed(BuildContext context) {
    return () async {
      final connConfig = ConnectionConfig(
        host: hostController.text,
        certificate: certController.text,
        caCertificate: caCertController.text,
        privateCertificate: privateCertController.text,
      );

      await config.setConnectionConfig(connConfig);
      await config.setPrivateKey(privateKeyController.text);

      Navigator.of(context)
          .pushReplacementNamed('/home', arguments: connConfig);
    };
  }
}
