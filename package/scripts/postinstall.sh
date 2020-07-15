#!/bin/bash
#systemctl daemon-reload
#systemctl enable gsheet-crm
#systemctl start gsheet-crm
/etc/init.d/gsheet-crm start
chkconfig --add gsheet-crm
chkconfig --level 345 gsheet-crm on