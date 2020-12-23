$shell_docker_installation = <<SCRIPT
apk add docker docker-compose
adduser vagrant docker
rc-update add docker boot
service docker start
SCRIPT

Vagrant.configure("2") do |config|
    config.vm.box = "generic/alpine312"

    config.vm.network "private_network", type: "dhcp"
    config.vm.network "forwarded_port", guest: ENV["NGINX_EXPOSE_PORT"], host: ENV["NGINX_EXPOSE_PORT"], host_ip: "localhost"
    config.vm.network "forwarded_port", guest: ENV["MYSQL_EXPOSE_PORT"], host: ENV["MYSQL_EXPOSE_PORT"], host_ip: "localhost"

    config.vm.synced_folder ".", "/home/vagrant/workspace", type: "nfs"

    config.vm.provision "shell", inline: $shell_docker_installation
end
