deploy:
	rm -f /opt/scripts/backupsync
	cp backupscript /opt/scripts/backupsync
	chmod +x /opt/scripts/backupsync