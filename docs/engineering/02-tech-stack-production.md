# Tech Stack — Production Blueprint

## Mobile

| Tech | Role |
|------|------|
| React Native (TypeScript) | SuperApp Copilot surface |
| Zustand | Client state |
| React Query | Server state cache |
| Alfa UI Kit / `@alfalab/core-components` (web) + RN kit | Brand components; Storybook: https://alfa-laboratory.github.io/core-components/ |
| react-native-markdown-display | Text fallback |
| Victory Native / Recharts RN | Charts |
| Reanimated + Lottie | Motion |

## Backend microservices

| Layer | Tech |
|-------|------|
| API Gateway | Kong or Nginx |
| Copilot BFF / dialog router | Go |
| AI Swarm | Python FastAPI + LangChain/LlamaIndex/LangGraph |
| Core banking adapters | Java / Spring Boot |
| Event streaming | Apache Kafka |
| Cache | Redis |
| Secrets | Bank Vault |

## AI / ML

| Component | Tech |
|-----------|------|
| Reasoning LLM | Llama 3 70B or Qwen 2.5 72B |
| Router LLM | Mixtral 8x7B |
| Inference | vLLM or TensorRT-LLM |
| Vector DB | Milvus or Qdrant |
| Embeddings | BGE-m3 or E5-large |
| Safety | Llama Guard (or bank equivalent) |
| OCR | Yandex Vision / Tesseract |

## Integrations

- Core transactions API  
- Compliance 115-ФЗ  
- Payments draft/sign  
- FNS / SMEV ENP  
- Push notification platform  
- Gosuslugi / Goskey (signing)  

## Observability

- OpenTelemetry traces across BFF → AI → tools  
- Audit log store (immutable)  
- GPU metrics (tokens/s, queue)  
- Eval regression in CI  

## Compatibility constraint

Идея может использовать внешние AI; **финальный** контур — on-prem и требования ИБ Альфы. Provider interface закладывается в demo.
