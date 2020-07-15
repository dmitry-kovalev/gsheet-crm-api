#!/bin/bash
#systemctl disable gsheet-crm
#systemctl stop gsheet-crm
#systemctl daemon-reload
/etc/init.d/gsheet-crm stop
chkconfig --del gsheet-crm
