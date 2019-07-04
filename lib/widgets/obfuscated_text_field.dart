import 'package:flutter/cupertino.dart';

typedef OnSubmittedBuilder = ValueChanged<String> Function(
  BuildContext context,
);

class ObfuscatedTextField extends StatelessWidget {
  final OnSubmittedBuilder onSubmittedBuilder;

  const ObfuscatedTextField({this.onSubmittedBuilder});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 5.0),
      child: CupertinoTextField(
        decoration: BoxDecoration(
          border: Border.all(color: CupertinoColors.black),
          borderRadius: const BorderRadius.all(Radius.circular(5.0)),
        ),
        clearButtonMode: OverlayVisibilityMode.editing,
        textAlign: TextAlign.center,
        onSubmitted: onSubmittedBuilder(context),
        obscureText: true,
      ),
    );
  }
}
