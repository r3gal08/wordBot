import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Wordlearner3000',
      theme: ThemeData(
        primarySwatch: Colors.deepOrange,
      ),
      home: const MyHomePage(title: 'Wordlearner3000'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});
  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final TextEditingController _wordController = TextEditingController();

  Future<void> _sendWordForDefinition() async {
    final String word = _wordController.text.trim();
    if (word.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please enter a word')),
      );
      return;
    }

    // TODO: Export local host, port, url, etc to a separate file similar to how I did in my bookBot project
    // Replace with your backend URL
    final url = Uri.parse('http://localhost:8080/wordHandler');

    try {
      final response = await http.post(
        url,
        headers: {
          'Content-Type': 'application/json',
        },
        body: jsonEncode({
          'word': word,
          "request": ["definition"]
        }), // Convert object to json string
      );

      // TODO: Here we will handle the json response and display the definition
      if (response.statusCode == 200) {
        // _wordController.clear(); // Clear the text field after success

        // Decode the JSON response
        final responseData = json.decode(response.body);
        final receivedWord = responseData['word'];
        final receivedDefinition = responseData['definition'];

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
              content: Text(
                  'Received word: $receivedWord, Definition: $receivedDefinition')),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              'Failed to send word. Status code: ${response.statusCode}',
            ),
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e')),
      );
    }
  }

Future<void> _sendWordToLearn() async {
  final String word = _wordController.text.trim();
  if (word.isEmpty) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Please enter a word')),
    );
    return;
  }

  // Replace with your backend URL
  final url = Uri.parse('http://localhost:8080/learnHandler');

  try {
    final response = await http.post(
      url,
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'word': word,
        "request": ["learn"]
      }), // Convert object to json string
    );

    if (response.statusCode == 200) {
      final responseData = json.decode(response.body);

      // Parse the answers and correct_answer from the response
      final List<String> answers = List<String>.from(responseData['answers']);
      final int correctAnswer = responseData['correct_answer'];

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            'Answers: $answers, Correct Answer Index: $correctAnswer',
          ),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            'Failed to send word. Status code: ${response.statusCode}',
          ),
        ),
      );
    }
  } catch (e) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Error: $e')),
    );
  }
}

  @override
  void dispose() {
    _wordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            TextField(
              controller: _wordController,
              decoration: const InputDecoration(
                labelText: 'Enter a word',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                ElevatedButton(
                  onPressed: _sendWordForDefinition,
                  child: const Text('Definition'),
                ),
                const SizedBox(width: 10),
                ElevatedButton(
                  onPressed: _sendWordToLearn,
                  child: const Text('Learn'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}