#!/bin/bash
systemctl daemon-reload
systemctl enable gsheet-crm
systemctl start gsheet-crm