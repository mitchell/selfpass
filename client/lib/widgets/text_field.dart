import 'package:flutter/cupertino.dart';

typedef OnSubmittedBuilder = ValueChanged<String> Function(
  BuildContext context,
);

class TextField extends StatelessWidget {
  final OnSubmittedBuilder onSubmittedBuilder;
  final TextEditingController controller;
  final OverlayVisibilityMode clearButtonMode;
  final Widget prefix;
  final Widget suffix;
  final bool obscure;
  final bool autofocus;
  final bool autocorrect;
  final int minLines;
  final int maxLines;

  const TextField({
    this.onSubmittedBuilder,
    this.controller,
    this.obscure = false,
    this.autofocus = false,
    this.minLines,
    this.maxLines,
    this.autocorrect = false,
    this.prefix,
    this.suffix,
    this.clearButtonMode = OverlayVisibilityMode.editing,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 5.0),
      child: CupertinoTextField(
        style: TextStyle(color: CupertinoColors.darkBackgroundGray),
        decoration: BoxDecoration(
          color: CupertinoColors.lightBackgroundGray,
          borderRadius: const BorderRadius.all(Radius.circular(5.0)),
        ),
        clearButtonMode: clearButtonMode,
        textAlign: TextAlign.start,
        onSubmitted: this.onSubmittedBuilder != null
            ? onSubmittedBuilder(context)
            : null,
        controller: controller,
        obscureText: obscure,
        autofocus: autofocus,
        autocorrect: autocorrect,
        minLines: minLines,
        maxLines: maxLines,
        prefix: prefix,
        suffix: suffix,
      ),
    );
  }
}
