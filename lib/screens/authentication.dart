import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

import '../types/abstracts.dart';

import '../widgets/obfuscated_text_field.dart';

class Authentication extends StatefulWidget {
  @override
  _AuthenticationState createState() => _AuthenticationState();
}

class _AuthenticationState extends State<Authentication> {
  bool _invalid = false;
  bool _passesDontMatch = false;
  String _masterpass;
  ConfigRepo _config;
  Future<bool> _passwordIsSet;

  @override
  didChangeDependencies() {
    super.didChangeDependencies();
    _config = Provider.of<ConfigRepo>(context);
  }

  @override
  Widget build(BuildContext context) {
    if (_passwordIsSet == null) {
      _passwordIsSet = _config.passwordSet;
    }

    return CupertinoPageScaffold(
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 50.0),
        child: FutureBuilder<bool>(
          future: _passwordIsSet,
          builder: (BuildContext context, AsyncSnapshot<bool> snapshot) =>
              snapshot.connectionState == ConnectionState.done
                  ? Column(children: _buildColumnChildren(snapshot.data))
                  : Container(),
        ),
      ),
    );
  }

  List<Widget> _buildColumnChildren(bool passwordIsSet) {
    List<Widget> children = [
      const Spacer(flex: 4),
      const Flexible(child: Text('Master password:')),
      Flexible(
        child: ObfuscatedTextField(
            onSubmittedBuilder:
                _buildMasterpassSubmittedBuilder(passwordIsSet)),
      ),
    ];

    if (!passwordIsSet) {
      children.add(const Flexible(child: Text('Re-enter password:')));
      children.add(Flexible(
        child: ObfuscatedTextField(
            onSubmittedBuilder:
                _buildConfirmPassSubmittedBuilder(passwordIsSet)),
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

    if (_passesDontMatch) {
      children.add(const Spacer(flex: 1));
    } else if (_invalid || !passwordIsSet) {
      children.add(const Spacer(flex: 2));
    } else {
      children.add(const Spacer(flex: 3));
    }

    return children;
  }

  OnSubmittedBuilder _buildMasterpassSubmittedBuilder(bool passwordIsSet) {
    return (BuildContext context) {
      return (String pass) async {
        if (passwordIsSet) {
          if (await _config.matchesPasswordHash(pass)) {
            Navigator.of(context).pushReplacementNamed('/home',
                arguments: await _config.connectionConfig);
            return;
          }

          this.setState(() => _invalid = true);
          return;
        }

        _masterpass = pass;
      };
    };
  }

  OnSubmittedBuilder _buildConfirmPassSubmittedBuilder(bool passwordIsSet) {
    return (BuildContext context) {
      return (String pass) async {
        if (pass != _masterpass) {
          this.setState(() {
            _passesDontMatch = true;
          });
          return;
        }

        _config.setPassword(_masterpass);
        _passwordIsSet = Future<bool>.value(true);
        Navigator.of(context).pushReplacementNamed('/home',
            arguments: await _config.connectionConfig);
      };
    };
  }
}
