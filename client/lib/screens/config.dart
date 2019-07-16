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
  TextEditingController _hostController;
  TextEditingController _caCertController;
  TextEditingController _certController;
  TextEditingController _privateCertController;
  TextEditingController _privateKeyController;
  ConnectionConfig _connectionConfig;
  String _privateKey;
  ConfigRepo _config;

  _ConfigState(this._connectionConfig, this._privateKey) {
    if (_connectionConfig == null) {
      _connectionConfig = ConnectionConfig();
    }

    _hostController = TextEditingController(text: _connectionConfig.host);
    _certController =
        TextEditingController(text: _connectionConfig.certificate);
    _caCertController =
        TextEditingController(text: _connectionConfig.caCertificate);
    _privateCertController =
        TextEditingController(text: _connectionConfig.privateCertificate);

    _privateKeyController = TextEditingController(text: _privateKey);
  }

  @override
  didChangeDependencies() async {
    super.didChangeDependencies();

    _config = Provider.of<ConfigRepo>(context);
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      navigationBar: _connectionConfig.host == null
          ? null
          : CupertinoNavigationBar(
              trailing: GestureDetector(
                onTap: _buildResetAllHandler(context),
                child: Text(
                  'Reset',
                  style: TextStyle(color: CupertinoColors.destructiveRed),
                ),
              ),
            ),
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 30),
        child: ListView(children: [
          Container(margin: EdgeInsets.only(top: 10), child: Text('Host:')),
          TextField(maxLines: 1, controller: _hostController),
          Container(
              margin: EdgeInsets.only(top: 5), child: Text('Private key:')),
          TextField(maxLines: 1, controller: _privateKeyController),
          Container(
              margin: EdgeInsets.only(top: 5), child: Text('CA certificate:')),
          TextField(maxLines: 5, controller: _caCertController),
          Container(
              margin: EdgeInsets.only(top: 5),
              child: Text('Client certificate:')),
          TextField(maxLines: 5, controller: _certController),
          Container(
              margin: EdgeInsets.only(top: 5),
              child: Text('Private certificate:')),
          TextField(maxLines: 5, controller: _privateCertController),
          CupertinoButton(
              child: Text('Save'), onPressed: _makeSaveOnPressed(context)),
        ]),
      ),
    );
  }

  GestureTapCallback _buildResetAllHandler(BuildContext context) {
    return () {
      showCupertinoDialog(
        context: context,
        builder: (BuildContext context) => CupertinoAlertDialog(
          content: Text('Are you sure?'),
          actions: [
            CupertinoDialogAction(
              isDefaultAction: true,
              child: Text('Cancel'),
              onPressed: () => Navigator.of(context).pop(),
            ),
            CupertinoDialogAction(
              isDestructiveAction: true,
              child: Text('Confirm'),
              onPressed: () async {
                _connectionConfig = null;
                await _config.deleteAll();
                Navigator.of(context)
                    .pushNamedAndRemoveUntil('/', ModalRoute.withName('/'));
              },
            ),
          ],
        ),
      );
    };
  }

  VoidCallback _makeSaveOnPressed(BuildContext context) {
    return () async {
      final connConfig = ConnectionConfig(
        host: _hostController.text,
        certificate: _certController.text,
        caCertificate: _caCertController.text,
        privateCertificate: _privateCertController.text,
      );

      await _config.setConnectionConfig(connConfig);
      await _config.setPrivateKey(_privateKeyController.text);

      Navigator.of(context)
          .pushReplacementNamed('/home', arguments: connConfig);
    };
  }
}
