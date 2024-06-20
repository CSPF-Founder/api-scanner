# API Security Scanner

## Prerequisites

1.Install Vagrant from the official site, https://developer.hashicorp.com/vagrant/downloads. 

- Please refer to this Installation guide if you face any issues during installation. https://developer.hashicorp.com/vagrant/docs/installation

2.Install Virtualbox from the official site, https://www.virtualbox.org/wiki/Downloads

## Minimum Spec

- 4 GB ram
- 4 CPU cores
- ~10 GB of free disk space

## Installing VM

Download this repository via 

`git clone https://github.com/CSPF-Founder/api-scanner.git`

Or you can download it as a zip file by clicking on "Code" in the top right and clicking "Download zip"

cd into the folder that is created.

In the project folder run the below command.

In Linux:

```
chmod +x setupvm.sh

./setupvm.sh
```

Once the vagrant installation is completed, it will automatically restart in Linux. 

In Windows cmd:

Go to the project folder on cmd and then run the below commands.


```
vagrant up
```

After it has completed

```
vagrant reload
```

In Windows, you need to manually reload as per the above command.


## Accessing the Panel

The API Security Scanner Panel is available on this URL: https://localhost:17443. 

```
Note: If you want to change the port, you can change forwardport in the vagrantfile.
```

For information on how to use the panel refer to [Manual.md](Manual.md)

## Further Reading:


- It is highly recommended to change the default password of the user 'vagrant' and change the SSH keys. 

- If you want to start the VM after your computer restarts you can give `vargant up` on this folder or start from the virtualbox manager. 

- Once up can access the VM by giving the command `vagrant ssh apiscannervm`

## Contributors

Sabari Selvan

Suriya Prakash
