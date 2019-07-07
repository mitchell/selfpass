import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/credential.dart';

import '../widgets/tappable_text_list.dart';

class Credentials extends StatelessWidget {
  final List<Metadata> metadatas;

  const Credentials(this.metadatas);

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: TappableTextList(tappableText: _buildTappableText(context)),
      navigationBar: CupertinoNavigationBar(),
    );
  }

  Map<String, GestureTapCallback> _buildTappableText(BuildContext context) {
    var makeOnTapHandler = (String id) => () async {
          final credential =
              await Provider.of<CredentialsRepo>(context).get(id);
          Navigator.of(context).pushNamed('/credential', arguments: credential);
        };

    Map<String, GestureTapCallback> tappableText = {};

    metadatas.forEach((Metadata metadata) {
      var primary = metadata.primary;
      if (metadata.tag != null && metadata.tag != '') {
        primary += "-" + metadata.tag;
      }
      tappableText[primary] = makeOnTapHandler(metadata.id);
    });

    return tappableText;
  }
}
