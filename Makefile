deploy:
	rm -f /opt/scripts/backupsync
	mkdir -p /opt/scripts/
	cp backupscript /opt/scripts/backupsync
	chmod +x /opt/scripts/backupsync