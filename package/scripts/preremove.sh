#!/bin/bash
systemctl disable gsheet-crm
systemctl stop gsheet-crm
systemctl daemon-reload
