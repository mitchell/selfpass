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
  CredentialsRepo client;
  ConfigRepo config;
  Future<List<Metadata>> metadatas;
  bool stateIsPaused = false;
  Timer pausedStateTimer;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    config = Provider.of<ConfigRepo>(context);
    client = Provider.of<CredentialsRepo>(context);

    metadatas = client.getAllMetadata('').toList();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    stateIsPaused = state == AppLifecycleState.paused;

    if (stateIsPaused) {
      pausedStateTimer = newPausedStateTimer();
      return;
    }

    if (pausedStateTimer != null) pausedStateTimer.cancel();
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: FutureBuilder<List<Metadata>>(
        future: metadatas,
        builder: (
          BuildContext context,
          AsyncSnapshot<List<Metadata>> snapshot,
        ) =>
            (snapshot.connectionState == ConnectionState.done)
                ? TappableTextList(
                    tappableText: buildTappableText(context, snapshot.data))
                : Center(child: CupertinoActivityIndicator()),
      ),
      navigationBar: CupertinoNavigationBar(
        leading: CupertinoButton(
          child: Text(
            'Lock',
            style: TextStyle(color: CupertinoColors.destructiveRed),
          ),
          onPressed: makeLockOnTapHandler(context),
          padding: EdgeInsets.zero,
        ),
        trailing: CupertinoButton(
          child: Icon(CupertinoIcons.gear),
          onPressed: makeConfigOnTapHandler(context),
          padding: EdgeInsets.zero,
        ),
      ),
    );
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    if (pausedStateTimer != null) pausedStateTimer.cancel();
    super.dispose();
  }

  Timer newPausedStateTimer() {
    const checkPeriod = 30;

    return Timer(Duration(seconds: checkPeriod), () {
      config.reset();
      Navigator.of(context)
          .pushNamedAndRemoveUntil('/', ModalRoute.withName('/home'));
    });
  }

  Map<String, GestureTapCallback> buildTappableText(
    BuildContext context,
    List<Metadata> metadatas,
  ) {
    final Map<String, List<Metadata>> metaMap = {};

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

  GestureTapCallback makeLockOnTapHandler(BuildContext context) {
    return () {
      config.reset();
      Navigator.of(context)
          .pushNamedAndRemoveUntil('/', ModalRoute.withName('/home'));
    };
  }

  GestureTapCallback makeConfigOnTapHandler(BuildContext context) {
    return () async => Navigator.of(context).pushNamed('/config',
        arguments: ConfigScreenArguments(
            connectionConfig: await config.connectionConfig,
            privateKey: await config.privateKey));
  }
}
