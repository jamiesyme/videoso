# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  #config.vm.box = "ubuntu/xenial64"
  config.vm.box = "xenial64_vb6"

  config.vm.network "forwarded_port", guest: 3000, host: 3000

  config.vm.provider "virtualbox" do |vb|
    vb.memory = 2048
    vb.cpus = 1
  end

  config.vm.synced_folder './web-client', '/web-client'

  # Bootstrap the app
  config.vm.provision "shell", inline: <<-SHELL
    cd /vagrant
    ./bootstrap.sh
  SHELL

  # Disable the Ubuntu console log from being generated.
  # See: https://groups.google.com/forum/#!topic/vagrant-up/eZljy-bddoI
  config.vm.provider "virtualbox" do |vb|
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
  end
end
