import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import 'repositories/credentials_client.dart';
import 'repositories/config.dart' as repo;

import 'screens/authentication.dart';
import 'screens/credentials.dart';
import 'screens/config.dart';
import 'screens/home.dart';

import 'types/abstracts.dart';

void main() => runApp(Selfpass());

class Selfpass extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Provider<ConfigRepo>(
      builder: (BuildContext context) => repo.Config(),
      child: CupertinoApp(
        title: 'Selfpass',
        onGenerateRoute: (RouteSettings settings) {
          String title;
          WidgetBuilder builder;

          switch (settings.name) {
            case '/':
              title = 'Authentication';
              builder = (BuildContext context) => Authentication();
              break;

            case '/home':
              title = 'Credential Hosts';
              builder = (BuildContext context) => Provider<CredentialsRepo>(
                    builder: (BuildContext context) =>
                        CredentialsClient(settings.arguments),
                    child: Home(),
                  );
              break;

            case '/credentials':
              title = 'Credentials';
              builder =
                  (BuildContext context) => Credentials(settings.arguments);
              break;

            case '/config':
              title = 'Configuration';
              builder = (BuildContext context) => Config(settings.arguments);
              break;
          }

          return CupertinoPageRoute(builder: builder, title: title);
        },
      ),
    );
  }
}
