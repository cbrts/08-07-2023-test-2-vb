#!/usr/bin/env bash
# Start the SSM agent
sudo systemctl enable amazon-ssm-agent
sudo systemctl start amazon-ssm-agent

# Install Ansible
curl https://bootstrap.pypa.io/get-pip.py -o /tmp/get-pip.py
python3 /tmp/get-pip.py --user
python3 -m pip install ansible --user

# Download the Ansible galaxy nginx role and install it
ansible-galaxy install geerlingguy.nginx 
cat <<EOF > playbook.yaml
---
- hosts: localhost
  roles:
    - { role: geerlingguy.nginx }
EOF
ansible-playbook playbook.yaml

# Backup current nginx config and enable default nginx config
sudo mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.old
sudo mv /etc/nginx/nginx.conf.default /etc/nginx/nginx.conf
sudo systemctl restart nginx
