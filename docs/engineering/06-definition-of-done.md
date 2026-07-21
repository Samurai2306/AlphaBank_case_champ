# Definition of Done

## Для задачи Demo

- [x] Код соответствует `docs/product/01-prd-demo-mvp.md` и contracts  
- [x] Shared schemas обновлены при изменении API/SDUI  
- [x] Unit/contract тесты зелёные  
- [x] Ручной прогон релевантного шага demo script  
- [x] Нет секретов в git  
- [x] Дисклеймер на налоговых поверхностях  
- [x] `deslop` / нет AI-slop в UI copy  
- [x] Документация обновлена, если менялось поведение  

## Для Demo MVP release (pitch)

- [x] P0 фичи D1–D10 (+ Home D10b) из PRD  
- [x] Backend = Go (`apps/api`), health/ready, SSE chat  
- [x] Сценарий «Маша» end-to-end  
- [x] Playwright smoke (landing + home tax hero; API required for home)  
- [x] README: как запустить за <15 минут  
- [x] Offline/scripted fallback на случай отсутствия LLM key  
- [x] `deploy/` stubs готовы к VPS (compose.prod + nginx) — см. `07-server-deployment.md`  

## Для Production increment

- [ ] Соответствие phase в `02-prd-production.md`  
- [ ] On-prem / masking / HITL требования фазы  
- [ ] Eval gate + security scan  
- [ ] Audit events для tool calls  
- [ ] Legal/compliance sign-off на user-facing tax copy  
- [ ] Runbook обновлён  

## Explicitly not Done

- «Работает на моей машине» без README  
- Ответы по налогам без дисклеймера  
- Платеж без confirm  
- Расхождение имён tools между web и docs  
