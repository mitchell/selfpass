import 'package:flutter/cupertino.dart';

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
    var handleOnTap = () {};

    Map<String, GestureTapCallback> tappableText = {};

    metadatas.forEach(
        (Metadata metadata) => tappableText[metadata.id] = handleOnTap);

    return tappableText;
  }
}
