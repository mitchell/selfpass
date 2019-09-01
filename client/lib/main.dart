import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import 'repositories/repositories.dart';

import 'screens/authentication.dart';
import 'screens/config.dart';
import 'screens/credential.dart';
import 'screens/credentials.dart';
import 'screens/home.dart';

import 'types/abstracts.dart';
import 'types/screen_arguments.dart';

void main() {
  runApp(Selfpass());
}

class Selfpass extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Provider<ConfigRepo>(
      builder: (BuildContext context) => EncryptedSharedPreferences(),
      child: CupertinoApp(
        title: 'Selfpass',
        theme: CupertinoThemeData(
            textTheme: CupertinoTextThemeData(
              brightness: Brightness.dark,
              textStyle: TextStyle(
                fontSize: 16,
                color: CupertinoColors.lightBackgroundGray,
              ),
            ),
            scaffoldBackgroundColor: CupertinoColors.darkBackgroundGray,
            brightness: Brightness.dark,
            primaryColor: CupertinoColors.activeBlue,
            primaryContrastingColor: CupertinoColors.darkBackgroundGray),
        onGenerateRoute: (RouteSettings settings) {
          String title;
          WidgetBuilder builder;

          switch (settings.name) {
            case '/':
              title = 'Authentication';
              builder = (BuildContext context) => Authentication();
              break;

            case '/home':
              title = 'Hosts';
              builder = (BuildContext context) => Provider<CredentialsRepo>(
                    builder: (BuildContext context) =>
                        GRPCCredentialsClient.getInstance(
                            config: settings.arguments),
                    child: Home(),
                  );
              break;

            case '/credentials':
              title = 'Credentials';
              builder = (BuildContext context) => Provider<CredentialsRepo>(
                    builder: (BuildContext context) =>
                        GRPCCredentialsClient.getInstance(),
                    child: Credentials(settings.arguments),
                  );
              break;

            case '/credential':
              title = 'Credential';
              builder = (BuildContext context) => Provider<CredentialsRepo>(
                    builder: (BuildContext context) =>
                        GRPCCredentialsClient.getInstance(),
                    child: Credential(settings.arguments),
                  );
              break;

            case '/config':
              final ConfigScreenArguments arguments = settings.arguments == null
                  ? ConfigScreenArguments()
                  : settings.arguments;

              title = 'Configuration';
              builder = (BuildContext context) =>
                  Config(arguments.connectionConfig, arguments.privateKey);
              break;
          }

          return CupertinoPageRoute(builder: builder, title: title);
        },
      ),
    );
  }
}
