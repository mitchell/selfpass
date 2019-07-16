import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';
import '../types/credential.dart';

import '../utils/crypto.dart' as crypto;

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
    final makeOnTapHandler = (String id) => () async {
          showCupertinoDialog(
            context: context,
            builder: (BuildContext context) => CupertinoAlertDialog(
                content: Column(
              children: [
                Text('Decrypting credential...'),
                Container(
                  margin: EdgeInsets.only(top: 10),
                  child: CupertinoActivityIndicator(),
                ),
              ],
            )),
          );

          final config = Provider.of<ConfigRepo>(context);
          final client = Provider.of<CredentialsRepo>(context);

          final Future<String> privateKey = config.privateKey;
          final String password = config.password;

          final credential = await client.get(id);

          credential.password = crypto.decrypt(
            credential.password,
            password,
            await privateKey,
          );

          if (credential.otpSecret.isNotEmpty) {
            credential.otpSecret = crypto.decrypt(
              credential.otpSecret,
              password,
              await privateKey,
            );
          }

          Navigator.of(context)
            ..pop()
            ..pushNamed('/credential', arguments: credential);
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
