user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/mk/user-api.mk release-test

social-rpc-dev:
	@make -f deploy/mk/social-rpc.mk release-test

social-api-dev:
	@make -f deploy/mk/social-api.mk release-test

im-ws-dev:
	@make -f deploy/mk/im-ws.mk release-test

im-rpc-dev:
	@make -f deploy/mk/im-rpc.mk release-test

im-api-dev:
	@make -f deploy/mk/im-api.mk release-test

task-mq-dev:
	@make -f deploy/mk/task-mq.mk release-test

release-test: user-rpc-dev user-api-dev social-rpc-dev social-api-dev im-ws-dev im-rpc-dev im-api-dev task-mq-dev

user-api-test: user-api-dev

social-api-test: social-api-dev

task-mq-test: task-mq-dev

im-ws-test: im-ws-dev

install-server:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh

install-server-user-api:
	cd ./deploy/script && chmod +x ./user-api-test.sh && ./user-api-test.sh

install-server-social-api:
	cd ./deploy/script && chmod +x ./social-api-test.sh && ./social-api-test.sh

install-server-user-rpc:
	cd ./deploy/script && chmod +x ./user-rpc-test.sh && ./user-rpc-test.sh

install-server-im-ws:
	cd ./deploy/script && chmod +x ./im-ws-test.sh && ./im-ws-test.sh