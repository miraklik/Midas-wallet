import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';

void main() {
  runApp(const MidasWalletApp());
}

class MidasWalletApp extends StatelessWidget {
  const MidasWalletApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Midas Wallet',
      theme: ThemeData(
        scaffoldBackgroundColor: Colors.black,
        appBarTheme: const AppBarTheme(
          backgroundColor: Colors.black,
          titleTextStyle: TextStyle(
            color: Colors.white,
            fontSize: 18,
            fontWeight: FontWeight.normal,
          ),
          centerTitle: true,
        ),
        colorScheme: ColorScheme.fromSeed(
          brightness: Brightness.dark,
          seedColor: const Color.fromARGB(255, 180, 135, 0),
        ),
        dividerColor: Colors.blueAccent,
        textTheme: TextTheme(
          bodyMedium: const TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.w700,
            fontSize: 16,
          ),
          labelSmall: const TextStyle(
            color: Colors.grey,
            fontWeight: FontWeight.w500,
            fontSize: 12,
          ),
        ),
      ),
      routes: {
        '/': (context) => MainPage(),
        '/coin': (contextl) => InformationTokenScreen(),
      },
      debugShowCheckedModeBanner: false,
    );
  }
}

class MainPage extends StatefulWidget {
  const MainPage({super.key});

  //final String title;

  @override
  State<MainPage> createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Text("Main", style: theme.textTheme.bodyMedium),
        backgroundColor: theme.appBarTheme.backgroundColor,
        leading: IconButton(
          onPressed: () {},
          icon: Icon(Icons.account_box_rounded),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            tooltip: 'Search',
            onPressed: () {},
          ),
          IconButton(
            icon: const Icon(Icons.qr_code_scanner_sharp),
            tooltip: 'QR-code',
            onPressed: () {},
          ),
        ],
      ),
      body: ListView.separated(
        itemCount: 30,
        separatorBuilder: (context, i) => const Divider(),
        itemBuilder: (context, i) {
          const coinName = "Bitcoin";
          return ListTile(
            leading: SvgPicture.asset(
              'images/svg/bitcoin.svg',
              height: 30,
              width: 30,
            ),
            trailing: Text('0.00 USD', style: theme.textTheme.bodyMedium),
            title: Text(coinName, style: theme.textTheme.bodyMedium),
            subtitle: Text('0 BTC', style: theme.textTheme.labelSmall),
            onTap: () {
              Navigator.of(context).pushNamed('/coin', arguments: coinName);
            },
          );
        },
      ),
    );
  }
}

class InformationTokenScreen extends StatefulWidget {
  const InformationTokenScreen({super.key});

  @override
  State<InformationTokenScreen> createState() => _InformationTokenScreenState();
}

class _InformationTokenScreenState extends State<InformationTokenScreen> {
  String? coinName;

  @override
  void didChangeDependencies() {
    final args = ModalRoute.of(context)?.settings.arguments;
    if (args == null) {
      log("You must provide args");
      return;
    }

    if (args is! String) {
      log("You must provide args");
      return;
    }
    coinName = args;
    setState(() {});
    super.didChangeDependencies();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          coinName ?? "Token not found",
          style: Theme.of(context).textTheme.bodyMedium,
        ),
        backgroundColor: Theme.of(context).appBarTheme.backgroundColor,
      ),
    );
  }
}
