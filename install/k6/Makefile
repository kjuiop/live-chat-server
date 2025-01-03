# 모든 작업 실행
dashboard: dashboard-up

k6: k6-up

k6-2: k6-2-up

dashboard-up:
	@docker-compose up -d prometheus grafana
	@echo "grafana dashboard on"

k6-up:
	@docker-compose up k6


k6-2-up:
	@docker-compose up k6-2
