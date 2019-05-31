deploy:
	rm -f /opt/scripts/backupsync
	ln -s /var/go/backupSync/current/backupsync /opt/scripts/backupsync