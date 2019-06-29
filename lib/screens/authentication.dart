import 'package:flutter/cupertino.dart';

class Authentication extends StatefulWidget {
  @override
  _AuthenticationState createState() => _AuthenticationState();
}

class _AuthenticationState extends State<Authentication> {
  static const String _masterpass = 'hunter#2';
  bool _invalid = false;

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 50.0),
        child: Column(
          children: _buildColumnChildren(context),
        ),
      ),
    );
  }

  List<Widget> _buildColumnChildren(BuildContext context) {
    final children = [
      const Spacer(flex: 4),
      const Flexible(child: Text('Master password:')),
      Flexible(
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 5.0),
          child: CupertinoTextField(
            decoration: BoxDecoration(
              border: Border.all(color: CupertinoColors.black),
              borderRadius: const BorderRadius.all(Radius.circular(5.0)),
            ),
            clearButtonMode: OverlayVisibilityMode.editing,
            textAlign: TextAlign.center,
            onSubmitted: _makeTextFieldOnSubmittedHandler(context),
            obscureText: true,
          ),
        ),
      ),
    ];

    if (_invalid) {
      children.add(const Flexible(
        child: Text(
          'invalid masterpass',
          style: TextStyle(color: CupertinoColors.destructiveRed),
        ),
      ));
      children.add(const Spacer(flex: 2));
    } else {
      children.add(const Spacer(flex: 3));
    }

    return children;
  }

  ValueChanged<String> _makeTextFieldOnSubmittedHandler(BuildContext context) {
    return (String pass) {
      if (pass != _masterpass) {
        this.setState(() {
          _invalid = true;
        });
        return;
      }

      Navigator.of(context).pushReplacementNamed('/home');
    };
  }
}
