import 'package:flutter/cupertino.dart';

import 'screens/authentication.dart';
import 'screens/home.dart';
import 'screens/credentials.dart';

void main() => runApp(Selfpass());

class Selfpass extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return CupertinoApp(
      title: 'Selfpass',
      routes: <String, WidgetBuilder>{
        '/': (BuildContext context) => Authentication(),
        '/home': (BuildContext context) => Home(),
        '/credentials': (BuildContext context) => Credentials(),
      },
    );
  }
}
