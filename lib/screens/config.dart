import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/connection_config.dart';

import '../widgets/text_field.dart';

class Config extends StatefulWidget {
  final ConnectionConfig connectionConfig;

  const Config(this.connectionConfig, {Key key}) : super(key: key);

  @override
  State createState() => _ConfigState(this.connectionConfig);
}

class _ConfigState extends State<Config> {
  TextEditingController _hostController;
  TextEditingController _caCertController;
  TextEditingController _certController;
  TextEditingController _privateCertController;
  ConnectionConfig _connectionConfig;
  ConfigRepo _config;

  _ConfigState(this._connectionConfig) {
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
                child: Text('Reset App',
                    style: TextStyle(color: CupertinoColors.destructiveRed)),
              ),
            ),
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 50.0),
        child: Column(children: [
          Spacer(flex: 3),
          Flexible(child: Text('Host:')),
          Flexible(child: TextField(maxLines: 1, controller: _hostController)),
          Flexible(child: Text('CA certificate:')),
          Flexible(
              child: TextField(maxLines: 3, controller: _caCertController)),
          Flexible(child: Text('Client certificate:')),
          Flexible(child: TextField(maxLines: 3, controller: _certController)),
          Flexible(child: Text('Private certificate:')),
          Flexible(
              child:
                  TextField(maxLines: 3, controller: _privateCertController)),
          CupertinoButton(
              child: Text('Save'), onPressed: _makeSaveOnPressed(context))
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

      Navigator.of(context)
          .pushReplacementNamed('/home', arguments: connConfig);
    };
  }
}
