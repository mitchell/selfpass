import 'package:flutter/cupertino.dart';

import '../widgets/tappable_text_list.dart';

class Credentials extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: TappableTextList(tappableText: _buildTappableText(context)),
      navigationBar: const CupertinoNavigationBar(middle: Text('Credentials')),
    );
  }

  static Map<String, GestureTapCallback> _buildTappableText(
      BuildContext context) {
    var handleOnTap = () {};

    Map<String, GestureTapCallback> tappableText = {
      'm@mjfs.us': handleOnTap,
      'm-mjfs': handleOnTap,
      'mitchelljfsimon@gmail.com-mitchelljfsimon': handleOnTap,
    };

    return tappableText;
  }
}
