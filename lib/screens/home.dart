import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/credential.dart';

import '../widgets/tappable_text_list.dart';

class Home extends StatefulWidget {
  @override
  State createState() => _HomeState();
}

class _HomeState extends State<Home> {
  CredentialsRepo _client;
  Future<List<Metadata>> _metadatas;

  @override
  didChangeDependencies() {
    super.didChangeDependencies();
    _client = Provider.of<CredentialsRepo>(context);
  }

  @override
  Widget build(BuildContext context) {
    if (_metadatas == null) {
      _metadatas = _client.getAllMetadata('').toList();
    }

    return CupertinoPageScaffold(
      child: FutureBuilder<List<Metadata>>(
        future: _metadatas,
        builder: (BuildContext context,
                AsyncSnapshot<List<Metadata>> snapshot) =>
            (snapshot.connectionState == ConnectionState.done)
                ? TappableTextList(
                    tappableText: _buildTappableText(context, snapshot.data))
                : Container(),
      ),
      navigationBar: CupertinoNavigationBar(),
    );
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

    final handleOnTap = (List<Metadata> metadatas) => () =>
        Navigator.of(context).pushNamed('/credentials', arguments: metadatas);

    final Map<String, GestureTapCallback> tappableText = {};

    metaMap.forEach((String key, List<Metadata> value) =>
        tappableText[key] = handleOnTap(value));

    return tappableText;
  }
}
