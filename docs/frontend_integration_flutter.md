## Flutter Integration: Multi-Message AI with Typing, Aggregation, and Input Control

This guide explains how the frontend should integrate with the backend for multi-message AI responses, with accurate typing indicators, user-input disabling, and support for aggregating multiple quick user messages into a single AI reply.

### Key Backend Behaviors

- The AI composes one overall reply, which the backend splits into several messages and stores sequentially with realistic delays (1–6s) between them.
- The backend maintains a live, in-memory typing tracker per conversation. The tracker reports:
  - `is_typing`: whether more AI chunks are pending
  - `message_index` and `total_messages`: current progress
  - The tracker starts immediately when the AI begins thinking and stops when the last chunk is stored.
- Quick successive user messages within an aggregation window are grouped and produce one AI reply:
  - Aggregation window: ~2s after the first message in the burst
  - Hard cap: if the user keeps typing, AI will start no later than ~6s after the first message

### Endpoints Used

- Send message:
  - `POST /api/v1/conversations/:id/messages`
  - Body minimal example: `{ "type": "text", "text": "hello" }`
  - Response: stored user message

- List messages:
  - `GET /api/v1/conversations/:id/messages`
  - Use for polling new messages (e.g., every 500–800ms while typing is true)

- Typing status (live, preferred):
  - `GET /api/v1/conversations/:id/typing-status`
  - Response data:
    - `is_typing`: boolean
    - `message_index`: number (0-based)
    - `total_messages`: number
    - `last_message_at`: ISO timestamp

### Frontend Logic (minimal)

1) When the user sends a message:
- Immediately disable the text input.
- Call `POST /messages` with the new message.
- Start two polling loops:
  - Typing loop: call `GET /typing-status` every 1s
  - Messages loop: call `GET /messages` every 600ms and append new companion messages

2) While typing-status `is_typing = true`:
- Keep the typing indicator visible.
- Keep input disabled.
- Optionally display progress: `message_index + 1` of `total_messages`.

3) When typing-status becomes false:
- Stop both polling loops.
- Hide typing indicator.
- Re-enable the input.

4) Aggregation of multiple user messages for one AI reply:
- The backend already handles this. If the user sends multiple messages within ~2s bursts (up to ~6s cap), they are treated as one “question” and will produce one AI multi-chunk answer.
- Frontend does not need special logic, other than continuing to disable input and showing typing until the AI completes.

### Minimal Model Hints (Dart)

Use these fields in your models; keep your own structure/style as needed:

```dart
class Message {
  // ... existing fields
  final bool isTyping;       // from backend
  final int messageIndex;    // from backend
  final int totalMessages;   // from backend
}

class TypingStatus {
  final bool isTyping;
  final int messageIndex;
  final int totalMessages;
  final DateTime lastMessageAt;
}
```

### Minimal API Calls (Dart-ish pseudocode)

```dart
Future<void> sendUserMessage(String conversationId, String text) async {
  disableInput();
  await dio.post('$base/conversations/$conversationId/messages', data: {
    'type': 'text',
    'text': text,
  });
  startTypingPoll(conversationId);   // 1s interval
  startMessagePoll(conversationId);  // 600ms interval
}

Future<TypingStatus> getTyping(String id) async {
  final r = await dio.get('$base/conversations/$id/typing-status');
  return TypingStatus.fromJson(r.data['data']);
}

Future<List<Message>> listMessages(String id) async {
  final r = await dio.get('$base/conversations/$id/messages');
  return (r.data['data'] as List).map((j) => Message.fromJson(j)).toList();
}
```

### UI States

- Show typing indicator when `typingStatus.isTyping == true`.
- Disable input when `typingStatus.isTyping == true`.
- Re-enable input when it turns false.
- Append new companion messages as they arrive; final chunk has `isTyping = false`.

### Error Handling

- If typing poll fails, stop polling and re-enable input.
- If message poll fails, stop that poll. Input can stay disabled until typing is false, or re-enable if you stop tracking.

### Notes on Tone

- Backend includes a system directive to reduce idioms/clichés and keep responses natural and human-like. No action needed on the frontend.

### QA Checklist

- Send 1 short user message → input disables → typing true → multiple AI chunks arrive with delays → typing false → input re-enabled.
- Send 2–3 user messages quickly (within ~2s) → single AI multi-chunk response → correct progress and typing.
- Long AI response → delays feel natural (1–6s) and punctuation adds slight pauses.


