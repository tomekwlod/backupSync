deploy:
	rm -f /opt/scripts/backupsync
	ln -f /var/go/backupSync/current/backupsync /opt/scripts/backupsync