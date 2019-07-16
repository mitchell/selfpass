import 'dart:async';

import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/credential.dart';
import '../types/screen_arguments.dart';

import '../widgets/tappable_text_list.dart';

class Home extends StatefulWidget {
  const Home({Key key}) : super(key: key);

  @override
  State createState() => _HomeState();
}

class _HomeState extends State<Home> with WidgetsBindingObserver {
  CredentialsRepo _client;
  ConfigRepo _config;
  Future<List<Metadata>> _metadatas;
  bool _stateIsPaused = false;
  Timer _pausedStateTimer;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    _config = Provider.of<ConfigRepo>(context);
    _client = Provider.of<CredentialsRepo>(context);

    _metadatas = _client.getAllMetadata('').toList();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    _stateIsPaused = state == AppLifecycleState.paused;

    if (_stateIsPaused) {
      _pausedStateTimer = _newPausedStateTimer();
      return;
    }

    if (_pausedStateTimer != null) _pausedStateTimer.cancel();
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: FutureBuilder<List<Metadata>>(
        future: _metadatas,
        builder: (
          BuildContext context,
          AsyncSnapshot<List<Metadata>> snapshot,
        ) =>
            (snapshot.connectionState == ConnectionState.done)
                ? TappableTextList(
                    tappableText: _buildTappableText(context, snapshot.data))
                : Center(child: CupertinoActivityIndicator()),
      ),
      navigationBar: CupertinoNavigationBar(
        leading: GestureDetector(
          child: Align(
              child: Text('Lock',
                  style: TextStyle(color: CupertinoColors.destructiveRed)),
              alignment: Alignment(-0.9, 0)),
          onTap: _makeLockOnTapHandler(context),
        ),
        trailing: GestureDetector(
          child: Icon(CupertinoIcons.gear),
          onTap: _makeConfigOnTapHandler(context),
        ),
      ),
    );
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    if (_pausedStateTimer != null) _pausedStateTimer.cancel();
    super.dispose();
  }

  Timer _newPausedStateTimer() {
    const checkPeriod = 30;

    return Timer(Duration(seconds: checkPeriod), () {
      _config.reset();
      Navigator.of(context)
          .pushNamedAndRemoveUntil('/', ModalRoute.withName('/home'));
    });
  }

  Map<String, GestureTapCallback> _buildTappableText(
    BuildContext context,
    List<Metadata> metadatas,
  ) {
    final Map<String, List<Metadata>> metaMap = {};

    metadatas.sort((a, b) => a.id.compareTo(b.id));

    for (var metadata in metadatas) {
      final source = metadata.sourceHost;

      if (metaMap[source] == null) {
        metaMap[source] = [metadata];
      } else {
        metaMap[source].add(metadata);
      }
    }

    final handleOnTap = (List<Metadata> metadatas) => () async =>
        Navigator.of(context).pushNamed('/credentials', arguments: metadatas);

    final Map<String, GestureTapCallback> tappableText = {};

    metaMap.forEach((String key, List<Metadata> value) =>
        tappableText[key] = handleOnTap(value));

    return tappableText;
  }

  GestureTapCallback _makeLockOnTapHandler(BuildContext context) {
    return () {
      _config.reset();
      Navigator.of(context)
          .pushNamedAndRemoveUntil('/', ModalRoute.withName('/home'));
    };
  }

  GestureTapCallback _makeConfigOnTapHandler(BuildContext context) {
    return () async => Navigator.of(context).pushNamed('/config',
        arguments: ConfigScreenArguments(
            connectionConfig: await _config.connectionConfig,
            privateKey: await _config.privateKey));
  }
}
