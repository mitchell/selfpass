import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';

import '../widgets/text_field.dart';

class Authentication extends StatefulWidget {
  @override
  _AuthenticationState createState() => _AuthenticationState();
}

class _AuthenticationState extends State<Authentication> {
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _confirmController = TextEditingController();
  bool _invalid = false;
  bool _passesDontMatch = false;
  ConfigRepo _config;
  Future<bool> _passwordIsSet;

  @override
  didChangeDependencies() {
    super.didChangeDependencies();
    _config = Provider.of<ConfigRepo>(context);
    _passwordIsSet = _config.passwordIsSet;
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 50.0),
        child: FutureBuilder<bool>(
          future: _passwordIsSet,
          builder: (BuildContext context, AsyncSnapshot<bool> snapshot) =>
              snapshot.connectionState == ConnectionState.done
                  ? Column(
                      children: _buildColumnChildren(context, snapshot.data))
                  : Center(child: CupertinoActivityIndicator()),
        ),
      ),
    );
  }

  List<Widget> _buildColumnChildren(BuildContext context, bool passwordIsSet) {
    List<Widget> children = [
      const Spacer(flex: 4),
      const Flexible(child: Text('Master password:')),
      Flexible(
        child: TextField(
          maxLines: 1,
          obscure: true,
          autofocus: true,
          controller: _passwordController,
        ),
      ),
    ];

    if (!passwordIsSet) {
      children.add(const Flexible(child: Text('Re-enter password:')));
      children.add(Flexible(
        child: TextField(
          maxLines: 1,
          obscure: true,
          controller: _confirmController,
        ),
      ));
    }

    if (_invalid) {
      children.add(const Flexible(
        child: Text(
          'invalid masterpass',
          style: TextStyle(color: CupertinoColors.destructiveRed),
        ),
      ));
    }

    if (_passesDontMatch) {
      children.add(const Flexible(
        child: Text(
          'passwords don\'t match',
          style: TextStyle(color: CupertinoColors.destructiveRed),
        ),
      ));
    }

    children.add(CupertinoButton(
      child: Text('Enter'),
      onPressed: _buildEnterPressedBuilder(context),
    ));

    if (_passesDontMatch) {
      children.add(const Spacer(flex: 1));
    } else if (_invalid || !passwordIsSet) {
      children.add(const Spacer(flex: 2));
    } else {
      children.add(const Spacer(flex: 3));
    }

    return children;
  }

  VoidCallback _buildEnterPressedBuilder(BuildContext context) {
    return () async {
      if (await _passwordIsSet) {
        if (await _config.matchesPasswordHash(_passwordController.text)) {
          Navigator.of(context).pushReplacementNamed('/home',
              arguments: await _config.connectionConfig);
          return;
        }

        this.setState(() => _invalid = true);
        return;
      }

      if (_passwordController.text != _confirmController.text) {
        this.setState(() {
          _passesDontMatch = true;
        });
        return;
      }

      _config.setPassword(_passwordController.text);
      _passwordIsSet = Future<bool>.value(true);
      Navigator.of(context).pushReplacementNamed('/config');
    };
  }
}
