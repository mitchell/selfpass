import 'package:flutter/cupertino.dart';
// import 'package:provider/provider.dart';

// import '../types/abstracts.dart';
// import '../types/credential.dart';
import '../widgets/tappable_text_list.dart';

class Home extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: TappableTextList(tappableText: _buildTappableText(context)),
      navigationBar:
          const CupertinoNavigationBar(middle: Text('Credentials Hosts')),
    );
  }

  static Map<String, GestureTapCallback> _buildTappableText(
      BuildContext context) {
    var handleOnTap = () => Navigator.of(context).pushNamed('/credentials');
    Map<String, GestureTapCallback> tappableText = {
      "google.com": handleOnTap,
      "amazon.com": handleOnTap,
      "linkedin.com": handleOnTap,
    };

    return tappableText;
  }
}
