SHELL := /bin/bash

bots-namespace:
	kubectl apply -f cluster/bots_namespace.yaml

# Please define these files in the config directory, you can refer to the example directory to define.
ti-sync-bot-config:
	kubectl apply -f config/ti_sync_bot_config.yaml

database-config:
	kubectl apply -f config/database_config.yaml

deploy: bots-namespace ti-sync-bot-config database-config
	kubectl apply -f cluster

undeploy:
	kubectl delete secret ti-sync-bot-config -n ti-community-bots
	kubectl delete secret database-config -n ti-community-bots
	kubectl delete deployment ti-sync-bot -n ti-community-bots
	kubectl delete service ti-sync-bot -n ti-community-bots
