# Redis 클러스터 설정

# Redis 노드 IP 및 포트 정의
REDIS_NODES = \
    172.30.0.2:6379 \
    172.30.0.3:6379 \
    172.30.0.4:6379 \
    172.30.0.5:6379 \
    172.30.0.6:6379 \
    172.30.0.7:6379

# 클러스터 복제본 수 (1: 하나의 마스터 노드당 하나의 슬레이브 노드)
CLUSTER_REPLICAS = 1

# Docker Compose 명령어
DOCKER_COMPOSE = docker-compose
DOCKER_EXEC = docker exec -it
DOCKER_FILE = cluster-official-compose.yml

# 기본 명령어
.PHONY: all up cluster create-nodes down clean

# 모든 작업 실행
all: up

install: up cluster

# Docker Compose로 Redis 노드 실행
up:
	@$(DOCKER_COMPOSE) -f $(DOCKER_FILE) up -d
	@echo "Redis 노드 실행 완료."

# Redis 클러스터 생성
cluster:
	@sleep 5 # Redis 노드가 준비될 때까지 대기
	@echo "Redis 클러스터 생성 중..."
	@$(DOCKER_EXEC) redis-node-1 redis-cli --cluster create $(REDIS_NODES) --cluster-replicas $(CLUSTER_REPLICAS)
	@echo "Redis 클러스터가 성공적으로 생성되었습니다."

# 노드 상태 확인
status-nodes:
	@echo "클러스터 노드 상태 확인:"
	@$(DOCKER_EXEC) redis-node-1 redis-cli -c cluster nodes

# Docker Compose 중지
down:
	@$(DOCKER_COMPOSE) down
	@echo "Redis 노드 종료 완료."

# 클러스터 초기화 및 데이터 정리
clean:
	@echo "클러스터 초기화 및 데이터 정리 중..."
	# 각 노드에 대해 CLUSTER RESET HARD 실행
	@$(DOCKER_EXEC) redis-node-1 redis-cli CLUSTER RESET HARD
	@$(DOCKER_EXEC) redis-node-2 redis-cli CLUSTER RESET HARD
	@$(DOCKER_EXEC) redis-node-3 redis-cli CLUSTER RESET HARD
	@$(DOCKER_EXEC) redis-node-4 redis-cli CLUSTER RESET HARD
	@$(DOCKER_EXEC) redis-node-5 redis-cli CLUSTER RESET HARD
	@$(DOCKER_EXEC) redis-node-6 redis-cli CLUSTER RESET HARD
	# Docker Compose로 Redis 노드 종료
	@$(DOCKER_COMPOSE) down
	@echo "Redis 노드 종료 완료."
	# 볼륨 정리
	@docker volume prune -f
	@echo "Redis 클러스터 데이터 및 볼륨 정리 완료."
