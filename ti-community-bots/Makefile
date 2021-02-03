SHELL := /bin/bash

bots-namespace:
	kubectl apply -f cluster/bots_namespace.yaml

ti-sync-bot-config:
	kubectl apply -f config/ti_sync_bot_config.yaml -n ti-community-bots

database-config:
	kubectl apply -f config/database_config.yaml -n ti-community-bots

deploy: bots-namespace ti-sync-bot-config database-config
	kubectl apply -f cluster -n ti-community-bots

undeploy:
	kubectl delete secret ti-sync-bot-config -n ti-community-bots
	kubectl delete secret database-config -n ti-community-bots
	kubectl delete deployment ti-sync-bot -n ti-community-bots
	kubectl delete service ti-sync-bot -n ti-community-bots
