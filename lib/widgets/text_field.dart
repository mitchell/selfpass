import 'package:flutter/cupertino.dart';

typedef OnSubmittedBuilder = ValueChanged<String> Function(
  BuildContext context,
);

class TextField extends StatelessWidget {
  final OnSubmittedBuilder onSubmittedBuilder;
  final TextEditingController controller;
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
  });

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
      ),
    );
  }
}
