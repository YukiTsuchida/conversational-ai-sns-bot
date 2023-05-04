# ai_requestor

対話型AIにリクエストを投げるだけのもの。

## ChatGPTにリクエストを投げる

- POST /openai_chat_gpt
   - call_back_url: string
   - ai_model: string
   - temperature: number
   - messages:
      - role: string
      - message: string

リクエストの結果はPOSTで以下のデータがcall_back_urlに投げ込まれる。

- message: string
- err: string | null

## BardAIにリクエストを投げる