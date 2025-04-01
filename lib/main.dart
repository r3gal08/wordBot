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
        brightness: Brightness.dark, // Set the theme to dark
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

  List<String> _answers = [];
  int _correctAnswerIndex = -1;
  int? _selectedAnswerIndex;    // Value can be nullable so we specify this here with '?

  Future<void> _sendWordForDefinition() async {
    final String word = _wordController.text.trim();
    if (word.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please enter a word')),
      );
      return;
    }

    // TODO: This should be pulled in via env file or something similar
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

    // TODO: This should be pulled in via env file or something similar
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
        }),
      );

      if (response.statusCode == 200) {
        final responseData = json.decode(response.body);
        setState(() {
          _answers = List<String>.from(responseData['answers']);
          _correctAnswerIndex = responseData['correct_answer'];
          _selectedAnswerIndex = null; // Reset selection
        });
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to send word. Status code: ${response.statusCode}'),
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e')),
      );
    }
  }

  void _checkAnswer(int index) {
    setState(() {
      _selectedAnswerIndex = index;
    });

    bool isCorrect = index == _correctAnswerIndex;

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(isCorrect ? 'Correct!' : 'Wrong! The correct answer is ${_answers[_correctAnswerIndex]}'),
        backgroundColor: isCorrect ? Colors.green : Colors.red,
      ),
    );
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

            // Learn and Definition buttons
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

            // 2x2 Grid for displaying answers
            const SizedBox(height: 20),
            if (_answers.isNotEmpty)
              Expanded(
                child: GridView.builder(
                  itemCount: _answers.length,
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    crossAxisSpacing: 10,
                    mainAxisSpacing: 10,
                  ),
                  itemBuilder: (context, index) {
                    Color buttonColor = Colors.blueGrey; // Default color
                    if (_selectedAnswerIndex != null) {
                      if (index == _correctAnswerIndex) {
                        buttonColor = Colors.green; // Correct answer
                      } else if (index == _selectedAnswerIndex) {
                        buttonColor = Colors.red; // Incorrect answer
                      }
                    }

                    return GestureDetector(
                      onTap: () {
                        if (_selectedAnswerIndex == null) {
                          _checkAnswer(index);
                        }
                      },
                      child: Container(
                        alignment: Alignment.center,
                        decoration: BoxDecoration(
                          color: buttonColor,
                          borderRadius: BorderRadius.circular(10),
                        ),
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          _answers[index],
                          style: const TextStyle(fontSize: 18, color: Colors.white),
                          textAlign: TextAlign.center,
                        ),
                      ),
                    );
                  },
                ),
              ),
          ],
        ),
      ),
    );
  }
}
