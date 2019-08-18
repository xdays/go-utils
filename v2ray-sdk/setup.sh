#!/bin/bash
# -*- coding: utf-8 -*-
 
cd config
v2ray -c server_api.json &
v2ray -c client_api.json &
