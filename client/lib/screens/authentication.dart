import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';

import '../widgets/text_field.dart';

class Authentication extends StatefulWidget {
  @override
  _AuthenticationState createState() => _AuthenticationState();
}

class _AuthenticationState extends State<Authentication> {
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController confirmController = TextEditingController();
  bool invalid = false;
  bool passesDontMatch = false;
  ConfigRepo config;
  Future<bool> passwordIsSet;

  @override
  didChangeDependencies() {
    super.didChangeDependencies();
    config = Provider.of<ConfigRepo>(context);
    passwordIsSet = config.passwordIsSet;
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 50.0),
        child: FutureBuilder<bool>(
          future: passwordIsSet,
          builder: (BuildContext context, AsyncSnapshot<bool> snapshot) =>
              snapshot.connectionState == ConnectionState.done
                  ? Column(
                      children: buildColumnChildren(context, snapshot.data))
                  : Center(child: CupertinoActivityIndicator()),
        ),
      ),
    );
  }

  List<Widget> buildColumnChildren(BuildContext context, bool passwordIsSet) {
    List<Widget> children = [
      const Spacer(flex: 4),
      const Flexible(child: Text('Master password:')),
      Flexible(
        child: TextField(
          maxLines: 1,
          obscure: true,
          autofocus: true,
          controller: passwordController,
        ),
      ),
    ];

    if (!passwordIsSet) {
      children.add(const Flexible(child: Text('Re-enter password:')));
      children.add(Flexible(
        child: TextField(
          maxLines: 1,
          obscure: true,
          controller: confirmController,
        ),
      ));
    }

    if (invalid) {
      children.add(const Flexible(
        child: Text(
          'invalid masterpass',
          style: TextStyle(color: CupertinoColors.destructiveRed),
        ),
      ));
    }

    if (passesDontMatch) {
      children.add(const Flexible(
        child: Text(
          'passwords don\'t match',
          style: TextStyle(color: CupertinoColors.destructiveRed),
        ),
      ));
    }

    children.add(Container(
      padding: EdgeInsets.only(top: 20),
      child: CupertinoButton.filled(
        child: Text('Enter'),
        onPressed: buildEnterPressedBuilder(context),
      ),
    ));

    if (passesDontMatch) {
      children.add(const Spacer(flex: 1));
    } else if (invalid || !passwordIsSet) {
      children.add(const Spacer(flex: 2));
    } else {
      children.add(const Spacer(flex: 3));
    }

    return children;
  }

  VoidCallback buildEnterPressedBuilder(BuildContext context) {
    return () async {
      if (await passwordIsSet) {
        if (await config.matchesPasswordHash(passwordController.text)) {
          Navigator.of(context).pushReplacementNamed('/home',
              arguments: await config.connectionConfig);
          return;
        }

        this.setState(() => invalid = true);
        return;
      }

      if (passwordController.text != confirmController.text) {
        this.setState(() {
          passesDontMatch = true;
        });
        return;
      }

      config.setPassword(passwordController.text);
      passwordIsSet = Future<bool>.value(true);
      Navigator.of(context).pushReplacementNamed('/config');
    };
  }
}
