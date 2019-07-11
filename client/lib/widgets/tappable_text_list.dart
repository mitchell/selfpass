import 'package:flutter/cupertino.dart';

class TappableTextList extends StatelessWidget {
  final Map<String, GestureTapCallback> tappableText;

  TappableTextList({Key key, this.tappableText}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: _buildListChildren(context),
    );
  }

  List<Widget> _buildListChildren(BuildContext context) {
    List<Widget> widgets = [];

    tappableText.forEach((String text, GestureTapCallback handleOnTap) {
      widgets.add(GestureDetector(
        onTap: handleOnTap,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 15.0),
          decoration: const BoxDecoration(
            border: Border(
              bottom: BorderSide(color: CupertinoColors.lightBackgroundGray),
            ),
          ),
          child: Text(text, textAlign: TextAlign.center),
        ),
      ));
    });

    return widgets;
  }
}
